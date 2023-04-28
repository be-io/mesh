/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"fmt"
	mtypes "github.com/be-io/mesh/client/golang/types"
	ptypes "github.com/traefik/paerser/types"
	"github.com/traefik/traefik/v3/pkg/config/static"
	"github.com/traefik/traefik/v3/pkg/ping"
	"github.com/traefik/traefik/v3/pkg/provider/acme"
	"github.com/traefik/traefik/v3/pkg/types"
	"os"
	"time"
)

func (that *meshProxy) Domains() []types.Domain {
	var domains []types.Domain
	for _, suf := range []string{"net", "cn", "com"} {
		domains = append(domains, types.Domain{
			Main: fmt.Sprintf("%s.%s", mtypes.CN, suf),
			SANs: []string{fmt.Sprintf("*.%s.%s", mtypes.CN, suf)},
		})
	}
	return domains
}

func (that *meshProxy) Configuration() *static.Configuration {
	// "https://auth.acme-dns.io"
	_ = os.Setenv("ACME_DNS_API_BASE", "https://acme.zerossl.com/v2/DV90")
	_ = os.Setenv("ACME_DNS_STORAGE_PATH", fmt.Sprintf("%s%s.json", that.Home, "acme"))
	return &static.Configuration{
		Global: &static.Global{
			CheckNewVersion: true,
		},
		EntryPoints: static.EntryPoints{
			TransportX: &static.EntryPoint{
				Address:          that.TransportX,
				ForwardedHeaders: &static.ForwardedHeaders{},
				Transport: &static.EntryPointsTransport{
					LifeCycle: &static.LifeCycle{
						RequestAcceptGraceTimeout: ptypes.Duration(time.Second * 32),
						GraceTimeOut:              ptypes.Duration(time.Second * 32),
					},
					RespondingTimeouts: &static.RespondingTimeouts{
						IdleTimeout: ptypes.Duration(time.Minute * 1),
					},
				},
				HTTP2: &static.HTTP2Config{
					MaxConcurrentStreams: proxy.MaxConcurrencyStream,
				},
			},
			TransportY: &static.EntryPoint{
				Address:          that.TransportY,
				ForwardedHeaders: &static.ForwardedHeaders{},
				Transport: &static.EntryPointsTransport{
					LifeCycle: &static.LifeCycle{
						RequestAcceptGraceTimeout: ptypes.Duration(time.Second * 32),
						GraceTimeOut:              ptypes.Duration(time.Second * 32),
					},
					RespondingTimeouts: &static.RespondingTimeouts{
						IdleTimeout: ptypes.Duration(time.Minute * 1),
					},
				},
				HTTP2: &static.HTTP2Config{
					MaxConcurrentStreams: proxy.MaxConcurrencyStream,
				},
			},
			TransportA: &static.EntryPoint{
				Address: that.TransportA,
				ForwardedHeaders: &static.ForwardedHeaders{
					Insecure: true,
				},
				Transport: &static.EntryPointsTransport{
					LifeCycle: &static.LifeCycle{
						RequestAcceptGraceTimeout: ptypes.Duration(time.Second * 32),
						GraceTimeOut:              ptypes.Duration(time.Second * 32),
					},
					RespondingTimeouts: &static.RespondingTimeouts{
						IdleTimeout: ptypes.Duration(time.Minute * 1),
					},
				},
				HTTP: static.HTTPConfig{
					TLS: &static.TLSConfig{
						Options:      mtypes.LocalNodeId,
						CertResolver: mtypes.LocalNodeId,
					},
				},
				HTTP2: &static.HTTP2Config{
					MaxConcurrentStreams: proxy.MaxConcurrencyStream,
				},
			},
			TransportB: &static.EntryPoint{
				Address: that.TransportB,
				ForwardedHeaders: &static.ForwardedHeaders{
					Insecure: true,
				},
				Transport: &static.EntryPointsTransport{
					LifeCycle: &static.LifeCycle{
						RequestAcceptGraceTimeout: ptypes.Duration(time.Second * 32),
						GraceTimeOut:              ptypes.Duration(time.Second * 32),
					},
					RespondingTimeouts: &static.RespondingTimeouts{
						IdleTimeout: ptypes.Duration(time.Minute * 1),
					},
				},
				HTTP: static.HTTPConfig{
					Redirections: &static.Redirections{
						EntryPoint: &static.RedirectEntryPoint{
							To:     TransportA,
							Scheme: "https",
						},
					},
				},
				HTTP2: &static.HTTP2Config{
					MaxConcurrentStreams: proxy.MaxConcurrencyStream,
				},
			},
			TransportC: &static.EntryPoint{
				Address: that.TransportC,
				ForwardedHeaders: &static.ForwardedHeaders{
					Insecure: true,
				},
				Transport: &static.EntryPointsTransport{
					LifeCycle: &static.LifeCycle{
						RequestAcceptGraceTimeout: ptypes.Duration(time.Second * 32),
						GraceTimeOut:              ptypes.Duration(time.Second * 32),
					},
					RespondingTimeouts: &static.RespondingTimeouts{
						IdleTimeout: ptypes.Duration(time.Minute * 1),
					},
				},
				HTTP: static.HTTPConfig{
					Redirections: &static.Redirections{
						EntryPoint: &static.RedirectEntryPoint{
							To:     TransportA,
							Scheme: "https",
						},
					},
				},
				HTTP2: &static.HTTP2Config{
					MaxConcurrentStreams: proxy.MaxConcurrencyStream,
				},
			},
			TransportD: &static.EntryPoint{
				Address: that.TransportD,
				ForwardedHeaders: &static.ForwardedHeaders{
					Insecure: true,
				},
				Transport: &static.EntryPointsTransport{
					LifeCycle: &static.LifeCycle{
						RequestAcceptGraceTimeout: ptypes.Duration(time.Second * 8),
						GraceTimeOut:              ptypes.Duration(time.Second * 32),
					},
					RespondingTimeouts: &static.RespondingTimeouts{
						IdleTimeout: ptypes.Duration(time.Minute * 1),
					},
				},
				HTTP: static.HTTPConfig{
					Redirections: &static.Redirections{
						EntryPoint: &static.RedirectEntryPoint{
							To:     TransportD,
							Scheme: "https",
						},
					},
				},
				HTTP2: &static.HTTP2Config{
					MaxConcurrentStreams: proxy.MaxConcurrencyStream,
				},
			},
		},
		CertificatesResolvers: map[string]static.CertificateResolver{
			Letsencrypt: {
				ACME: &acme.Configuration{
					Email:          "coyzeng@gmail.com",
					CAServer:       "https://acme.zerossl.com/v2/DV90",
					PreferredChain: "",
					Storage:        fmt.Sprintf("%s%s.json", that.Home, "acme"),
					KeyType:        "EC384",
					EAB: &acme.EAB{
						Kid:         "ZaI5UXvg0Vw63POFsoBbSA",
						HmacEncoded: "WWup77c0lzXvA5-zeNgfaJnwvIEbZK_b5gEgWAWFInzeLRaQUI72wrUgv86kPdJH4IW9bCNFw0qgMGNLEDJojA",
					},
					CertificatesDuration: 365 * 24,
					DNSChallenge: &acme.DNSChallenge{
						Provider:  "acme-dns",
						Resolvers: []string{"1.1.1.1:53", "8.8.8.8:53"},
					},
				},
			},
			mtypes.LocalNodeId: {
				ACME: &acme.Configuration{
					Email:          "coyzeng@gmail.com",
					CAServer:       "https://acme.zerossl.com/v2/DV90",
					PreferredChain: "",
					Storage:        fmt.Sprintf("%s%s.json", that.Home, mtypes.LocalNodeId),
					KeyType:        "EC384",
					EAB: &acme.EAB{
						Kid:         "ZaI5UXvg0Vw63POFsoBbSA",
						HmacEncoded: "WWup77c0lzXvA5-zeNgfaJnwvIEbZK_b5gEgWAWFInzeLRaQUI72wrUgv86kPdJH4IW9bCNFw0qgMGNLEDJojA",
					},
					CertificatesDuration: 365 * 24,
					TLSChallenge:         &acme.TLSChallenge{}}},
		},
		Providers: &static.Providers{
			ProvidersThrottleDuration: ptypes.Duration(2 * time.Second),
		},
		ServersTransport: &static.ServersTransport{
			MaxIdleConnsPerHost: that.MaxTransport,
			InsecureSkipVerify:  true,
		},
		AccessLog: &types.AccessLog{
			FilePath:      that.AccessPath,
			Format:        "common",
			BufferingSize: 200,
			Fields: &types.AccessLogFields{
				DefaultMode: "keep",
				Names:       map[string]string{},
				Headers: &types.FieldHeaders{
					DefaultMode: "keep",
					Names: map[string]string{
						"RequestHost": "keep",
					},
				},
			},
		},
		//Tracing: &static.Tracing{
		//	ServiceName:   plugin.Whoami,
		//	SpanNameLimit: 255,
		//	Jaeger:        nil,
		//	Zipkin:        nil,
		//	Datadog:       nil,
		//	Instana:       nil,
		//	Haystack:      nil,
		//	Elastic:       nil,
		//},
		//Metrics: &types.Metrics{
		//	Prometheus: nil,
		//	Datadog:    nil,
		//	StatsD:     nil,
		//	InfluxDB:   nil,
		//	InfluxDB2:  nil,
		//},
		HostResolver: &types.HostResolverConfig{},
		API: &static.API{
			Insecure:  true,
			Dashboard: true,
			Debug:     false,
		},
		Ping: &ping.Handler{
			EntryPoint: TransportY,
		},
		Log: &types.TraefikLog{
			Level:  "INFO",
			Format: "common",
		},
		Experimental: nil,
	}
}
