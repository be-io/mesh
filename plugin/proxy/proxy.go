/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"context"
	"crypto/x509"
	_ "embed"
	"fmt"
	"github.com/coreos/go-systemd/v22/daemon"
	"github.com/go-acme/lego/v4/challenge"
	gokitmetrics "github.com/go-kit/kit/metrics"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/plugin"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/tool"
	mtypes "github.com/opendatav/mesh/client/golang/types"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"github.com/traefik/traefik/v3/cmd/healthcheck"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"github.com/traefik/traefik/v3/pkg/config/runtime"
	"github.com/traefik/traefik/v3/pkg/config/static"
	"github.com/traefik/traefik/v3/pkg/metrics"
	"github.com/traefik/traefik/v3/pkg/middlewares/accesslog"
	"github.com/traefik/traefik/v3/pkg/provider/acme"
	"github.com/traefik/traefik/v3/pkg/provider/aggregator"
	"github.com/traefik/traefik/v3/pkg/provider/tailscale"
	"github.com/traefik/traefik/v3/pkg/provider/traefik"
	"github.com/traefik/traefik/v3/pkg/safe"
	"github.com/traefik/traefik/v3/pkg/server"
	"github.com/traefik/traefik/v3/pkg/server/middleware"
	"github.com/traefik/traefik/v3/pkg/server/router/tcp"
	"github.com/traefik/traefik/v3/pkg/server/service"
	pcp "github.com/traefik/traefik/v3/pkg/tcp"
	traefiktls "github.com/traefik/traefik/v3/pkg/tls"
	"github.com/traefik/traefik/v3/pkg/tracing"
	"github.com/traefik/traefik/v3/pkg/types"
	"github.com/traefik/traefik/v3/pkg/udp"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"net/http"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
)

func init() {
	plugin.Provide(proxy)
}

const (
	ProviderName    = "mesh"
	TransportX      = "transport_x"
	TransportY      = "transport_y"
	TransportA      = "transport_a"
	TransportB      = "transport_b"
	TransportC      = "traefik"
	TransportD      = "transport_d"
	Letsencrypt     = "letsencrypt"
	PluginRetry     = "retry"
	PluginBarrier   = "barrier"
	PluginHeader    = "header"
	PluginAuthority = "authority"
	PluginErrors    = "errors"
	PluginRewrite   = "rewrite"
	PluginReplace   = "replace"
	PluginHath      = "hath"
)

var proxy = new(meshProxy)

type meshProxy struct {
	TransportX           string                 `json:"plugin.proxy.transport_x" dft:"0.0.0.0:570" usage:"Encrypt data frame message address"`
	TransportY           string                 `json:"plugin.proxy.transport_y" dft:"0.0.0.0:7304" usage:"Data frame message address"`
	TransportA           string                 `json:"plugin.proxy.transport_a" dft:"0.0.0.0:443" usage:"User frontend message address"`
	TransportB           string                 `json:"plugin.proxy.transport_b" dft:"0.0.0.0:80" usage:"User frontend message address"`
	TransportC           string                 `json:"plugin.proxy.transport_c" dft:"0.0.0.0:8866" usage:"Internal service address"`
	TransportD           string                 `json:"plugin.proxy.transport_d" dft:"0.0.0.0:8867" usage:"Frontend service address"`
	Home                 string                 `json:"plugin.proxy.home" dft:"${MESH_HOME}/mesh/proxy/" usage:"path to store certification"`
	MaxConcurrencyStream int32                  `json:"plugin.proxy.stream.max" dft:"400" usage:"Max concurrency stream per connection"`
	MaxTransport         int                    `json:"plugin.proxy.transport.max" dft:"400" usage:"Max transports"`
	StatefulEnable       bool                   `json:"plugin.proxy.stateful.enable" dft:"true" usage:"Stateful route enable"`
	RateLimitEnable      bool                   `json:"plugin.proxy.limit.enable" dft:"true" usage:"Rate limit enable"`
	InsecureEnable       bool                   `json:"plugin.proxy.insecure.enable" dft:"true" usage:"Insecure enable"`
	AccessPath           string                 `json:"plugin.proxy.access.path" dft:"${MESH_HOME}/mesh/proxy/logs/access.log" usage:"path to store access log"`
	AccessBackups        int                    `json:"plugin.proxy.access.backups" dft:"120" usage:"Max archive files backups"`
	AccessSize           int                    `json:"plugin.proxy.access.size" dft:"100" usage:"Max archive files size in megabytes"`
	AccessAge            int                    `json:"plugin.proxy.access.age" dft:"28" usage:"Max archive files age in days"`
	Compress             bool                   `json:"plugin.proxy.access.compress" dft:"false" comment:"Compress the accesslog"`
	InsecureSkip         bool                   `json:"plugin.proxy.insecure.skip" dft:"false" comment:"Insecure skip verify"`
	Server               *server.Server         `json:"-"`
	TCPRouters           map[string]*tcp.Router `json:"-"`
	UDPRouters           map[string]udp.Handler `json:"-"`
	RollingWriter        io.WriteCloser         `json:"-"`
}

func (that *meshProxy) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Name: plugin.Proxy, Flags: meshProxy{}, Create: func() plugin.Plugin {
		return that
	}}
}

func (that *meshProxy) Start(ctx context.Context, runtime plugin.Runtime) {
	log.Catch(runtime.Parse(that))
	config := that.Configuration()
	config.SetEffectiveConfiguration()
	err := config.ValidateConfiguration()
	if nil != err {
		log.Error(ctx, err.Error())
		return
	}
	if err = tool.MakeDir(that.Home); nil != err {
		log.Error(ctx, err.Error())
		return
	}
	// 11G, 120 backups, 100M
	that.RollingWriter = &lumberjack.Logger{
		Filename:   that.AccessPath,
		MaxSize:    that.AccessSize, // megabytes
		MaxBackups: that.AccessBackups,
		MaxAge:     that.AccessAge, //days
		Compress:   that.Compress,  // disabled by default
	}
	if that.Server, err = that.setupServer(ctx, config); nil != err {
		log.Error(ctx, err.Error())
		return
	}
	runtime.Submit(func() {
		http.DefaultTransport.(*http.Transport).Proxy = http.ProxyFromEnvironment
		that.startTraefik(ctx, config)
	})
	runtime.StartHook(func() {
		taskId, err := aware.Scheduler.Period(ctx, time.Second*10, &mtypes.Topic{
			Topic: prsim.RoutePeriodRefresh.Topic,
			Code:  prsim.RoutePeriodRefresh.Code,
		})
		if nil != err {
			log.Error(ctx, err.Error())
		} else {
			log.Info(ctx, "Mesh proxy dynamic routers refresh start with %s. ", taskId)
		}
	})
	that.Register(ctx)
}

func (that *meshProxy) Stop(ctx context.Context, runtime plugin.Runtime) {
	if nil != that.Server {
		that.Server.Close()
	}
	if nil != that.RollingWriter {
		log.Catch(that.RollingWriter.Close())
	}
}

func (that *meshProxy) RecursionBreak(ctx context.Context, uri mtypes.URC) bool {
	for _, addr := range uri.Members() {
		if strings.Index(addr, "0.0.0.0") < 0 || strings.Index(addr, "127.0.0.1") < 0 {
			continue
		}
		adds := strings.Split(addr, ":")
		if len(adds) < 2 {
			continue
		}
		if strings.Contains(that.TransportX, fmt.Sprintf(":%s", adds[1])) || strings.Contains(that.TransportY, fmt.Sprintf(":%s", adds[1])) {
			return true
		}
	}
	return false
}

func (that *meshProxy) startTraefik(ctx context.Context, staticConfiguration *static.Configuration) {
	ntx, _ := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	if staticConfiguration.Ping != nil {
		staticConfiguration.Ping.WithContext(ntx)
	}
	that.Server.Start(ntx)
	sent, err := daemon.SdNotify(false, "READY=1")
	if !sent && nil != err {
		log.Error(ntx, "Failed to notify: %v", err)
	}
	t, err := daemon.SdWatchdogEnabled(false)
	if nil != err {
		log.Error(ntx, "Could not enable Watchdog: %v", err)
	} else if t != 0 {
		// Send a ping each half time given
		t /= 2
		log.Info(ntx, "Watchdog activated with timer duration %s", t)
		safe.Go(func() {
			tick := time.Tick(t)
			for range tick {
				resp, errHealthCheck := healthcheck.Do(*staticConfiguration)
				if resp != nil {
					_ = resp.Body.Close()
				}

				if staticConfiguration.Ping == nil || errHealthCheck == nil {
					if ok, _ := daemon.SdNotify(false, "WATCHDOG=1"); !ok {
						log.Error(ntx, "Fail to tick watchdog")
					}
				} else {
					log.Error(ntx, errHealthCheck.Error())
				}
			}
		})
	}
	log.Info(ntx, "Plugin proxy has been started.")
	that.Server.Wait()
	log.Info(ntx, "Plugin proxy has been shutdown")
}

func (that *meshProxy) setupServer(ctx context.Context, staticConfiguration *static.Configuration) (*server.Server, error) {
	accessLog, err := accesslog.NewHandlerWithFormatWriter(staticConfiguration.AccessLog, that.RollingWriter, logger)
	if nil != err {
		return nil, cause.Error(err)
	}
	routinesPool := safe.NewPool(ctx)
	providerAggregator := aggregator.NewProviderAggregator(*staticConfiguration.Providers)
	if err = providerAggregator.AddProvider(meshGraph); nil != err {
		return nil, cause.Error(err)
	}
	// adds internal provider
	if err = providerAggregator.AddProvider(traefik.New(*staticConfiguration)); nil != err {
		return nil, cause.Error(err)
	}
	// ACME
	tlsManager := traefiktls.NewManager()
	httpChallengeProvider := acme.NewChallengeHTTP()
	tlsChallengeProvider := acme.NewChallengeTLSALPN()
	if err = providerAggregator.AddProvider(tlsChallengeProvider); nil != err {
		return nil, cause.Error(err)
	}
	acmeProviders := that.initACMEProvider(ctx, staticConfiguration, &providerAggregator, tlsManager, httpChallengeProvider, tlsChallengeProvider)
	// Tailscale
	tsProviders := initTailscaleProviders(ctx, staticConfiguration, &providerAggregator)
	// Metrics
	metricRegistries := that.registerMetricClients(ctx, staticConfiguration.Metrics)
	metricsRegistry := metrics.NewMultiRegistry(metricRegistries)
	// Entrypoints
	serverEntryPointsTCP, err := server.NewTCPEntryPoints(staticConfiguration.EntryPoints, staticConfiguration.HostResolver, metricsRegistry)
	if nil != err {
		return nil, cause.Error(err)
	}
	serverEntryPointsUDP, err := server.NewUDPEntryPoints(staticConfiguration.EntryPoints)
	if nil != err {
		return nil, cause.Error(err)
	}
	// Plugins
	pluginBuilder, err := createPluginBuilder(staticConfiguration)
	if nil != err {
		log.Error(ctx, "Plugins are disabled because an error has occurred. %s", err.Error())
	}
	// Providers plugins
	for name, conf := range staticConfiguration.Providers.Plugin {
		if nil == pluginBuilder {
			break
		}
		provider, err := pluginBuilder.BuildProvider(name, conf)
		if nil != err {
			return nil, cause.Errorf("plugin: failed to build provider: %w", err)
		}
		if err = providerAggregator.AddProvider(provider); nil != err {
			return nil, cause.Errorf("plugin: failed to add provider: %w", err)
		}
	}

	// Service manager factory
	var spiffeX509Source *workloadapi.X509Source
	if staticConfiguration.Spiffe != nil && staticConfiguration.Spiffe.WorkloadAPIAddr != "" {
		log.Info(ctx, "Waiting on SPIFFE SVID delivery, workloadAPIAddr is %s", staticConfiguration.Spiffe.WorkloadAPIAddr)
		spiffeX509Source, err = workloadapi.NewX509Source(ctx, workloadapi.WithClientOptions(workloadapi.WithAddr(staticConfiguration.Spiffe.WorkloadAPIAddr)))
		if nil != err {
			return nil, cause.Errorf("unable to create SPIFFE x509 source: %w", err)
		}
		log.Info(ctx, "Successfully obtained SPIFFE SVID.")
	}
	roundTripperManager := service.NewRoundTripperManager(spiffeX509Source)
	dialerManager := pcp.NewDialerManager(spiffeX509Source)
	acmeHTTPHandler := that.getHTTPChallengeHandler(acmeProviders, httpChallengeProvider)
	managerFactory := service.NewManagerFactory(*staticConfiguration, routinesPool, metricsRegistry, roundTripperManager, acmeHTTPHandler)
	tracer, tracerCloser := setupTracing(ctx, staticConfiguration.Tracing)
	chainBuilder := middleware.NewChainBuilder(metricsRegistry, accessLog, tracer)
	routerFactory := server.NewRouterFactory(*staticConfiguration, managerFactory, tlsManager, chainBuilder, pluginBuilder, metricsRegistry, dialerManager)
	// Watcher
	watcher := server.NewConfigurationWatcher(routinesPool, providerAggregator, that.getDefaultsEntrypoints(ctx, staticConfiguration), "internal")
	// TLS
	watcher.AddListener(func(conf dynamic.Configuration) {
		tlsManager.UpdateConfigs(ctx, conf.TLS.Stores, conf.TLS.Options, conf.TLS.Certificates)
		gauge := metricsRegistry.TLSCertsNotAfterTimestampGauge()
		for _, certificate := range tlsManager.GetServerCertificates() {
			that.appendCertMetric(gauge, certificate)
		}
	})
	// Metrics
	watcher.AddListener(func(_ dynamic.Configuration) {
		metricsRegistry.ConfigReloadsCounter().Add(1)
		metricsRegistry.LastConfigReloadSuccessGauge().Set(float64(time.Now().Unix()))
	})
	// Server Transports
	watcher.AddListener(func(conf dynamic.Configuration) {
		roundTripperManager.Update(conf.HTTP.ServersTransports)
		dialerManager.Update(conf.TCP.ServersTransports)
	})
	// Switch router
	watcher.AddListener(that.switchRouter(routerFactory, serverEntryPointsTCP, serverEntryPointsUDP))
	// Metrics
	if metricsRegistry.IsEpEnabled() || metricsRegistry.IsRouterEnabled() || metricsRegistry.IsSvcEnabled() {
		var eps []string
		for key := range serverEntryPointsTCP {
			eps = append(eps, key)
		}
		watcher.AddListener(func(conf dynamic.Configuration) {
			metrics.OnConfigurationUpdate(conf, eps)
		})
	}
	// TLS challenge
	watcher.AddListener(tlsChallengeProvider.ListenConfiguration)
	// Certificate Resolvers
	resolverNames := map[string]struct{}{}
	// ACME
	for _, p := range acmeProviders {
		resolverNames[p.ResolverName] = struct{}{}
		watcher.AddListener(p.ListenConfiguration)
	}
	// Tailscale
	for _, p := range tsProviders {
		resolverNames[p.ResolverName] = struct{}{}
		watcher.AddListener(p.HandleConfigUpdate)
	}
	// Certificate resolver logs
	watcher.AddListener(func(config dynamic.Configuration) {
		for rtName, rt := range config.HTTP.Routers {
			if rt.TLS == nil || rt.TLS.CertResolver == "" {
				continue
			}
			if _, ok := resolverNames[rt.TLS.CertResolver]; !ok {
				log.Warn(ctx, "Router %s uses a non-existent certificate resolver: %s", rtName, rt.TLS.CertResolver)
			}
		}
	})
	return server.NewServer(routinesPool, serverEntryPointsTCP, serverEntryPointsUDP, watcher, chainBuilder, accessLog, tracerCloser), nil
}

func (that *meshProxy) getHTTPChallengeHandler(acmeProviders []*acme.Provider, httpChallengeProvider http.Handler) http.Handler {
	var acmeHTTPHandler http.Handler
	for _, p := range acmeProviders {
		if p != nil && p.HTTPChallenge != nil {
			acmeHTTPHandler = httpChallengeProvider
			break
		}
	}
	return acmeHTTPHandler
}

func (that *meshProxy) getDefaultsEntrypoints(ctx context.Context, staticConfiguration *static.Configuration) []string {
	var defaultEntryPoints []string

	// Determines if at least one EntryPoint is configured to be used by default.
	var hasDefinedDefaults bool
	for _, ep := range staticConfiguration.EntryPoints {
		if ep.AsDefault {
			hasDefinedDefaults = true
			break
		}
	}

	for name, cfg := range staticConfiguration.EntryPoints {
		// By default, all entrypoints are considered.
		// If at least one is flagged, then only flagged entrypoints are included.
		if hasDefinedDefaults && !cfg.AsDefault {
			continue
		}

		protocol, err := cfg.GetProtocol()
		if nil != err {
			// Should never happen because Traefik should not start if protocol is invalid.
			log.Error(ctx, "Invalid protocol: %v", err)
		}
		if protocol != "udp" && name != static.DefaultInternalEntryPointName {
			defaultEntryPoints = append(defaultEntryPoints, name)
		}
	}
	sort.Strings(defaultEntryPoints)
	return defaultEntryPoints
}

func (that *meshProxy) switchRouter(routerFactory *server.RouterFactory, serverEntryPointsTCP server.TCPEntryPoints, serverEntryPointsUDP server.UDPEntryPoints) func(conf dynamic.Configuration) {
	return func(conf dynamic.Configuration) {
		rtConf := runtime.NewConfig(conf)
		routers, udpRouters := routerFactory.CreateRouters(rtConf)
		serverEntryPointsTCP.Switch(routers)
		serverEntryPointsUDP.Switch(udpRouters)
		that.TCPRouters = routers
		that.UDPRouters = udpRouters
	}
}

// initACMEProvider creates and registers acme.Provider instances corresponding to the configured ACME certificate resolvers.
func (that *meshProxy) initACMEProvider(ctx context.Context, c *static.Configuration, providerAggregator *aggregator.ProviderAggregator, tlsManager *traefiktls.Manager, httpChallengeProvider, tlsChallengeProvider challenge.Provider) []*acme.Provider {
	localStores := map[string]*acme.LocalStore{}
	var resolvers []*acme.Provider
	for name, resolver := range c.CertificatesResolvers {
		if resolver.ACME == nil {
			continue
		}
		if localStores[resolver.ACME.Storage] == nil {
			localStores[resolver.ACME.Storage] = acme.NewLocalStore(resolver.ACME.Storage)
		}
		p := &acme.Provider{
			Configuration:         resolver.ACME,
			Store:                 localStores[resolver.ACME.Storage],
			ResolverName:          name,
			HTTPChallengeProvider: httpChallengeProvider,
			TLSChallengeProvider:  tlsChallengeProvider,
		}
		if err := providerAggregator.AddProvider(p); nil != err {
			log.Error(ctx, "The ACME resolver %q is skipped from the resolvers list because: %v", name, err)
			continue
		}
		p.SetTLSManager(tlsManager)
		p.SetConfigListenerChan(make(chan dynamic.Configuration))
		resolvers = append(resolvers, p)
	}
	return resolvers
}

// initTailscaleProviders creates and registers tailscale.Provider instances corresponding to the configured Tailscale certificate resolvers.
func initTailscaleProviders(ctx context.Context, cfg *static.Configuration, providerAggregator *aggregator.ProviderAggregator) []*tailscale.Provider {
	var providers []*tailscale.Provider
	for name, resolver := range cfg.CertificatesResolvers {
		if resolver.Tailscale == nil {
			continue
		}

		tsProvider := &tailscale.Provider{ResolverName: name}

		if err := providerAggregator.AddProvider(tsProvider); nil != err {
			log.Error(ctx, "Unable to create Tailscale provider %s: %v", name, err)
			continue
		}

		providers = append(providers, tsProvider)
	}

	return providers
}

func (that *meshProxy) registerMetricClients(ctx context.Context, metricsConfig *types.Metrics) []metrics.Registry {
	if metricsConfig == nil {
		return nil
	}
	var registries []metrics.Registry
	if metricsConfig.Prometheus != nil {
		prometheusRegister := metrics.RegisterPrometheus(ctx, metricsConfig.Prometheus)
		if prometheusRegister != nil {
			registries = append(registries, prometheusRegister)
			log.Debug(ctx, "Configured Prometheus metrics")
		}
	}
	if metricsConfig.Datadog != nil {
		registries = append(registries, metrics.RegisterDatadog(ctx, metricsConfig.Datadog))
		log.Debug(ctx, "Configured Datadog metrics: pushing to %s once every %s",
			metricsConfig.Datadog.Address, metricsConfig.Datadog.PushInterval)
	}
	if metricsConfig.StatsD != nil {
		registries = append(registries, metrics.RegisterStatsd(ctx, metricsConfig.StatsD))
		log.Debug(ctx, "Configured StatsD metrics: pushing to %s once every %s",
			metricsConfig.StatsD.Address, metricsConfig.StatsD.PushInterval)
	}
	if metricsConfig.InfluxDB2 != nil {
		influxDB2Register := metrics.RegisterInfluxDB2(ctx, metricsConfig.InfluxDB2)
		if influxDB2Register != nil {
			registries = append(registries, influxDB2Register)
			log.Debug(ctx, "Configured InfluxDB v2 metrics: pushing to %s (%s org/%s bucket) once every %s",
				metricsConfig.InfluxDB2.Address, metricsConfig.InfluxDB2.Org, metricsConfig.InfluxDB2.Bucket, metricsConfig.InfluxDB2.PushInterval)
		}
	}
	if metricsConfig.OpenTelemetry != nil {
		openTelemetryRegistry := metrics.RegisterOpenTelemetry(ctx, metricsConfig.OpenTelemetry)
		if openTelemetryRegistry != nil {
			registries = append(registries, openTelemetryRegistry)
			log.Debug(ctx, "Configured OpenTelemetry metrics: pushing to %s once every %s",
				metricsConfig.OpenTelemetry.Address, metricsConfig.OpenTelemetry.PushInterval.String())
		}
	}
	return registries
}

func (that *meshProxy) appendCertMetric(gauge gokitmetrics.Gauge, certificate *x509.Certificate) {
	sort.Strings(certificate.DNSNames)
	labels := []string{
		"cn", certificate.Subject.CommonName,
		"serial", certificate.SerialNumber.String(),
		"sans", strings.Join(certificate.DNSNames, ","),
	}
	notAfter := float64(certificate.NotAfter.Unix())
	gauge.With(labels...).Set(notAfter)
}

func setupTracing(ctx context.Context, conf *static.Tracing) (trace.Tracer, io.Closer) {
	if nil == conf {
		return nil, nil
	}
	tracer, closer, err := tracing.NewTracing(conf)
	if err != nil {
		log.Warn(ctx, "Unable to create tracer")
		return nil, nil
	}

	return tracer, closer
}
