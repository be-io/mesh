/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package dns

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/plugin"
	"github.com/be-io/mesh/client/golang/prsim"
	dp "github.com/be-io/mesh/plugin/dns/plugin"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"os"
	"runtime"
)

func init() {
	var dns = new(meshDNS)
	//caddy.RegisterCaddyfileLoader("flag", caddy.LoaderFunc(dns.ConfLoader))
	caddy.SetDefaultCaddyfileLoader("default", caddy.LoaderFunc(dns.DefaultLoader))
	plugin.Provide(dns)
}

const serverType = "dns"

// meshDNS is wrapper of dns server daemon.
type meshDNS struct {
	Conf string `json:"conf" yaml:"conf"`
	// List installed plugins
	Plugins bool `json:"plugins" yaml:"plugins"`
	// Path to write pid file
	PIDFile string `json:"pid_file" yaml:"pid_file"`
	// Show version
	Version bool `json:"version" yaml:"version"`
	// Quiet mode (no initialization output)
	Quiet bool `json:"quiet" yaml:"quiet"`

	instance *caddy.Instance
}

func (that *meshDNS) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Name: plugin.DNS, Flags: meshDNS{}, Create: func() plugin.Plugin {
		return that
	}}
}

func (that *meshDNS) Start(ctx context.Context, run plugin.Runtime) {
	plugins := append(dnsserver.Directives, dp.Server().Name())
	// don't show init stuff from caddy
	dnsserver.Quiet = that.Quiet
	dnsserver.Directives = plugins
	caddy.Quiet = that.Quiet
	caddy.PidFile = that.PIDFile
	caddy.DefaultConfigFile = "Corefile"
	caddy.AppName = string(plugin.DNS)
	caddy.AppVersion = prsim.Version
	caddy.TrapSignals()

	log.Info(ctx, "DNS daemon %s-%s(%s %s %s) start", caddy.AppName, caddy.AppVersion, runtime.GOOS, runtime.GOARCH, runtime.Version())
	log.Info(ctx, "Read conf from %s, write pid to %s", that.Conf, that.PIDFile)

	// Get Core file input
	profile, err := caddy.LoadCaddyfile(serverType)
	if nil != err {
		log.Error(ctx, "DNS daemon start failure, %s", err.Error())
		return
	}

	// Start your engines
	instance, err := caddy.Start(profile)
	if nil != err {
		log.Error(ctx, "DNS daemon start failure, %s", err.Error())
		return
	}

	// Twiddle your thumbs
	that.instance = instance
	// that.instance.Wait()
}

func (that *meshDNS) Stop(ctx context.Context, runtime plugin.Runtime) {
	if nil == that.instance {
		return
	}
	if err := that.instance.Stop(); nil != err {
		log.Error(ctx, "DNS daemon stop failure, %s", err.Error())
	}
	that.instance.Wait()
}

// ConfLoader loads the Caddyfile using the -conf flag.
func (that *meshDNS) ConfLoader(serverType string) (caddy.Input, error) {
	if that.Conf == "" {
		return caddy.CaddyfileInput{
			Filepath:       that.Conf,
			ServerTypeName: serverType,
		}, nil
	}

	if that.Conf == "stdin" {
		return caddy.CaddyfileFromPipe(os.Stdin, serverType)
	}

	contents, err := os.ReadFile(that.Conf)
	if nil != err {
		return nil, err
	}
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       that.Conf,
		ServerTypeName: serverType,
	}, nil
}

// DefaultLoader loads the Corefile from the current working directory.
func (that *meshDNS) DefaultLoader(serverType string) (caddy.Input, error) {
	contents, err := os.ReadFile(caddy.DefaultConfigFile)
	if nil != err && os.IsNotExist(err) {
		return caddy.CaddyfileInput{
			Contents:       []byte(fmt.Sprintf(".:53 {\nmesh\nlog\n}\n")),
			Filepath:       caddy.DefaultConfigFile,
			ServerTypeName: serverType,
		}, nil
	}
	if nil != err {
		return nil, err
	}
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       caddy.DefaultConfigFile,
		ServerTypeName: serverType,
	}, nil
}
