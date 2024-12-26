/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package plugin

import (
	"context"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"os"
	"os/signal"
	"runtime/debug"
	"sort"
	"sync"
	"time"
)

func init() {
	var _ Container = new(sharedContainer)
}

type sharedContainer struct {
	mustRun       bool             // optional
	waitAny       bool             // optional
	priority      int              // optional
	name          Name             // required
	once          sync.Once        // optional
	flags         []string         // optional
	shutdownHooks []func()         // optional
	startingHooks []func()         // optional
	waitAnyLocks  []Locker         // optional
	parent        *sharedContainer // optional
	myself        Plugin           // optional
	children      sharedContainers // optional
}

func (that *sharedContainer) Init(flags ...string) {
	attr := Load(that.name).Ptt()
	that.once.Do(func() {
		that.flags = append(that.flags, flags...)
		that.myself = attr.Create()
		that.priority = attr.Priority
		that.waitAny = attr.WaitAny
	})
}

func (that *sharedContainer) RunAware(ctx context.Context) {
	for _, some := range macro.Load(prsim.IRuntimeAware).List() {
		if aware, ok := some.(prsim.RuntimeAware); ok {
			if err := aware.Init(); nil != err {
				log.Error(ctx, "Prepare runtime with unexpected %s", err.Error())
				return
			}
		}
	}
}

func (that *sharedContainer) Start(ctx context.Context, flags ...string) {
	defer func() {
		if err := recover(); nil != err {
			log.Error(ctx, string(debug.Stack()))
			log.Error(ctx, "Mesh start %s with unexpected err, %v", that.name, err)
			that.Stop(ctx)
		}
	}()
	// Root plugin's parent is nil
	if nil == that.parent {
		that.RunAware(ctx)
	}
	if nil != that.parent {
		log.Info(ctx, "Mesh starting plugin %s", that.name)
	}
	that.Init(flags...)
	that.myself.Start(ctx, that)
	sort.Sort(that.children)

	for _, child := range that.children {
		if child.waitAny && !that.mustRun {
			log.Info(ctx, "Mesh post starting mutex plugin %s", child.name)
			continue
		}
		child.mustRun = that.mustRun
		child.Start(ctx, flags...)
	}

	for _, hook := range that.startingHooks {
		hook()
	}

	// Root plugin's parent is nil
	if nil == that.parent {
		that.WaitLocks(ctx, time.Minute*5)
		for _, child := range that.children {
			if child.waitAny {
				child.mustRun = true
				child.Start(ctx, flags...)
			}
		}
		log.Debug(ctx, "Mesh started.")
	}
}

func (that *sharedContainer) Stop(ctx context.Context) {
	for _, child := range that.children {
		child.Stop(ctx)
	}

	that.myself.Stop(ctx, that)

	for _, hook := range that.shutdownHooks {
		hook()
	}
}

func (that *sharedContainer) Wait(ctx context.Context) {
	// Block until signalled.
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	that.Stop(ctx)
}

func (that *sharedContainer) WaitLocks(ctx context.Context, duration time.Duration) {
	if len(that.waitAnyLocks) < 1 {
		return
	}
	ticker := time.NewTicker(time.Second * 6)
	timeout := time.NewTicker(duration)
	defer func() {
		ticker.Stop()
		timeout.Stop()
	}()
	waited := func() bool {
		for _, locker := range that.waitAnyLocks {
			ok, err := locker.TryLock()
			if nil != err {
				log.Warn(ctx, "Try to startup mutex plugin, %s", err.Error())
				return false
			}
			if !ok {
				return false
			}
		}
		return true
	}
	if waited() {
		log.Info(ctx, "Mutex plugins wait startup succeed. ")
		return
	}
	for {
		select {
		case _, open := <-ticker.C:
			if !open {
				log.Warn(ctx, "Mutex plugin wait ticker has been closed. ")
				return
			}
			if !waited() {
				continue
			}
			log.Info(ctx, "Mutex plugins wait startup succeed. ")
			return
		case <-timeout.C:
			log.Warn(ctx, "Mutex plugins wait startup has been timeout. ")
			return
		}
	}
}

func (that *sharedContainer) WaitAny(ctx context.Context, locker Locker) {
	if nil != that.parent {
		that.parent.WaitAny(ctx, locker)
	} else {
		that.waitAnyLocks = append(that.waitAnyLocks, locker)
	}
}

func (that *sharedContainer) Plugins(plugins ...Name) {
	for _, plugin := range plugins {
		that.Load(plugin)
	}
}

func (that *sharedContainer) Parse(ptr interface{}) error {
	return flags.readAsStruct(string(that.name), ptr, that.flags...)
}

func (that *sharedContainer) Parse0(ptr interface{}, conf string) error {
	return flags.readConfAsStruct(conf, ptr)
}

func (that *sharedContainer) Load(plugin Name) {
	for _, children := range that.children {
		if children.name == plugin {
			return
		}
	}
	child := &sharedContainer{name: plugin, parent: that}
	child.Init(that.flags...)
	that.children = append(that.children, child)
}

func (that *sharedContainer) Submit(routine func()) {
	log.Devour(func() {
		go routine()
	})
}

func (that *sharedContainer) StartHook(hook func()) {
	if nil != that.parent && !that.mustRun {
		that.parent.StartHook(hook)
	} else {
		that.startingHooks = append(that.startingHooks, hook)
	}
}

func (that *sharedContainer) ShutdownHook(hook func()) {
	if nil != that.parent {
		that.parent.ShutdownHook(hook)
	} else {
		that.shutdownHooks = append(that.shutdownHooks, hook)
	}
}

type sharedContainers []*sharedContainer

func (that sharedContainers) Len() int {
	return len(that)
}

func (that sharedContainers) Less(i, j int) bool {
	return that[j].priority < that[i].priority
}

func (that sharedContainers) Swap(i, j int) {
	tmp := that[j]
	that[j] = that[i]
	that[i] = tmp
}
