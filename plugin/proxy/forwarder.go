package proxy

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/traefik/traefik/v2/pkg/server/middleware"
	"net/http"
	"strings"
)

func init() {
	var _ http.Handler = new(forwarder)
	var _ prsim.Listener = forwarders
	middleware.Provide(forwarders)
	macro.Provide(prsim.IListener, forwarders)
}

var forwarders = &forwarderMiddleware{license: &types.License{}, forwards: map[string]string{}}

type forwarder struct {
	name string
	next http.Handler
}

func (that *forwarder) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if nil == proxy.TCPRouters || nil == proxy.TCPRouters[TransportY] {
		that.next.ServeHTTP(writer, request)
		return
	}
	if target, ok := forwarders.Forwarding(request.Host); ok {
		if "" == target {
			http.NotFoundHandler().ServeHTTP(writer, request)
			return
		}
		request.Host = target
		request.Header.Set("Host", target)
		proxy.TCPRouters[TransportY].GetHTTPHandler().ServeHTTP(writer, request)
	} else {
		that.next.ServeHTTP(writer, request)
	}
}

type forwarderMiddleware struct {
	license  *types.License
	services []*types.Service
	forwards map[string]string
}

func (that *forwarderMiddleware) Name() string {
	return fmt.Sprintf("%s@%s", PluginForwarder, ProviderName)
}

func (that *forwarderMiddleware) Priority() int {
	return 0
}

func (that *forwarderMiddleware) Scope() int {
	return 1
}

func (that *forwarderMiddleware) New(ctx context.Context, next http.Handler, name string) (http.Handler, error) {
	return &forwarder{next: next, name: name}, nil
}

func (that *forwarderMiddleware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.proxy.forwarder"}
}

func (that *forwarderMiddleware) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.LicenseImports, prsim.RegistryEventRefresh}
}

func (that *forwarderMiddleware) Listen(ctx context.Context, event *types.Event) error {
	if event.Binding.Match(prsim.LicenseImports) {
		var license *types.License
		if err := event.TryGetObject(&license); nil != err {
			return cause.Error(err)
		}
		that.license = tool.Anyone(license, new(types.License))
	}
	if event.Binding.Match(prsim.RegistryEventRefresh) {
		var registrations types.MetadataRegistrations
		if err := event.TryGetObject(&registrations); nil != err {
			return cause.Error(err)
		}
		that.services = registrations.Of(types.METADATA).InferService()
	}
	forwards := map[string]string{}
	names := map[string]string{}
	for _, name := range that.license.SuperURN {
		sn := strings.Join(strings.Split(name, ".")[1:], ".")
		names[name] = sn
		forwards[sn] = ""
	}
	for _, service := range that.services {
		urn := types.FromURN(ctx, service.URN)
		if sn := names[urn.Name]; "" != sn {
			forwards[sn] = service.URN
		}
	}
	that.forwards = forwards
	return nil
}

// Forwarding must emit 404
func (that *forwarderMiddleware) Forwarding(urn string) (string, bool) {
	if len(that.license.SuperURN) < 1 {
		return "", false
	}
	ctx := mpc.Context()
	uname := types.FromURN(ctx, urn)
	for name, to := range that.forwards {
		if "" != name && strings.HasSuffix(uname.Name, name) {
			return to, true
		}
	}
	return "", false
}
