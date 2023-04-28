/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package boost

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/grpc"
	"github.com/be-io/mesh/client/golang/http"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"os"
	"os/signal"
)

func init() {
	var _ prsim.RuntimeHook = new(Mooter)
}

type Mooter struct {
	http1 mpc.Provider
	http2 mpc.Provider
}

func (that *Mooter) Start(ctx context.Context, runtime prsim.Runtime) error {
	runtimes := macro.Load(prsim.IRuntimeAware).List()
	for _, runner := range runtimes {
		ra, ok := runner.(prsim.RuntimeAware)
		if !ok {
			log.Warn(ctx, "Aware %s cant be cast to RuntimeAware", runner.Att().Name)
			continue
		}
		if err := ra.Init(); nil != err {
			log.Fatal(ctx, "Aware %s exec, %s", runner.Att().Name, err.Error())
		}
	}
	hooks := macro.Load(prsim.IRuntimeHook).List()
	for _, hook := range hooks {
		rh, ok := hook.(prsim.RuntimeHook)
		if !ok {
			log.Warn(ctx, "Hook %s cant be cast to RuntimeHook", hook.Att().Name)
			continue
		}
		if err := rh.Start(ctx, runtime); nil != err {
			log.Fatal(ctx, "Hook %s exec, %s", hook.Att().Name, err.Error())
		}
	}
	http2, ok := macro.Load(mpc.IProvider).Get(grpc.Name).(mpc.Provider)
	if !ok {
		return cause.Errorf("No Provider named %s exist ", grpc.Name)
	}
	http1, ok := macro.Load(mpc.IProvider).Get(http.Name).(mpc.Provider)
	if !ok {
		return cause.Errorf("No Provider named %s exist ", http.Name)
	}
	that.http1 = http1
	that.http2 = http2
	runtime.Submit(func() {
		log.Catch(that.http1.Start(ctx, fmt.Sprintf("0.0.0.0:%d", tool.Runtime.Get().Port+1), nil))
	})
	runtime.Submit(func() {
		log.Catch(that.http2.Start(ctx, fmt.Sprintf("0.0.0.0:%d", tool.Runtime.Get().Port), nil))
	})
	return nil
}

func (that *Mooter) Stop(ctx context.Context, runtime prsim.Runtime) error {
	hooks := macro.Load(prsim.IRuntimeHook).List()
	for _, hook := range hooks {
		rh, ok := hook.(prsim.RuntimeHook)
		if !ok {
			log.Warn(ctx, "Hook %s cant be cast to RuntimeHook", hook.Att().Name)
			continue
		}
		if err := rh.Stop(ctx, runtime); nil != err {
			log.Fatal(ctx, "Hook %s exec, %s", hook.Att().Name, err.Error())
		}
	}
	if nil != that.http2 {
		log.Catch(that.http2.Close())
	}
	if nil != that.http1 {
		log.Catch(that.http1.Close())
	}
	return nil
}

func (that *Mooter) Refresh(ctx context.Context, runtime prsim.Runtime) error {
	return that.Start(ctx, runtime)
}

func (that *Mooter) Wait(ctx context.Context, runtime prsim.Runtime) {
	// Block until signalled.
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	if err := that.Stop(ctx, runtime); nil != err {
		log.Error(ctx, err.Error())
	}
}
