/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package boost

import (
	"context"
	"github.com/gofrs/flock"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/plugin"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/tool"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

type Daemon struct {
	PIDFd   string
	Name    string
	Hooks   []prsim.RuntimeHook
	Waiter  prsim.Waiter
	Runtime prsim.Runtime
}

func (that *Daemon) tryLockEvenIfAbsent(ctx context.Context) (plugin.Locker, int, error) {
	if err := tool.MakeDir(filepath.Dir(that.PIDFd)); nil != err {
		return nil, -1, cause.Error(err)
	}
	if err := tool.MakeFile(that.PIDFd); nil != err {
		return nil, -1, cause.Error(err)
	}
	locker := flock.New(that.PIDFd)
	text, err := os.ReadFile(that.PIDFd)
	if nil != err || "" == string(text) {
		return locker, -1, err
	}
	pid, err := strconv.Atoi(string(text))
	return locker, pid, err
}

func (that *Daemon) isProcessLive(pid int) bool {
	if process, _ := os.FindProcess(pid); nil != process {
		return nil == process.Signal(syscall.Signal(0))
	}
	return false
}

func (that *Daemon) writePID(ctx context.Context, pid int, npid int) {
	if that.isProcessLive(pid) {
		log.Warn(ctx, "%s already boot with pid %d. ", that.Name, pid)
	}
	if err := os.WriteFile(that.PIDFd, []byte(strconv.Itoa(npid)), 0644); nil != err {
		log.Error(ctx, err.Error())
	}
}

func (that *Daemon) Run(ctx context.Context) error {
	locker, pid, err := that.tryLockEvenIfAbsent(ctx)
	if nil != err {
		return cause.Error(err)
	}
	ok, err := locker.TryLock()
	if nil != err {
		return cause.Error(err)
	}
	if !ok && that.isProcessLive(pid) {
		log.Warn(ctx, "%s cant lock %s, already boot with pid %d. ", that.Name, that.PIDFd, pid)
		return nil
	}
	log.AddProcessShutdownHook(func() error {
		log.Catch(locker.Unlock())
		log.Catch(os.Remove(that.PIDFd))
		return nil
	})
	that.writePID(ctx, pid, os.Getpid())
	for _, hook := range that.Hooks {
		if err = hook.Start(ctx, that.Runtime); nil != err {
			return cause.Error(err)
		}
	}
	log.Info(ctx, "%s has been started with %d. ", that.Name, os.Getpid())
	if nil != that.Waiter {
		that.Waiter.Wait(ctx, that.Runtime)
	}
	return nil
}

func (that *Daemon) Shutdown(ctx context.Context) error {
	locker, pid, err := that.tryLockEvenIfAbsent(ctx)
	if nil != err {
		return cause.Error(err)
	}
	ok, err := locker.TryLock()
	if nil != err {
		return cause.Error(err)
	}
	if ok {
		log.Warn(ctx, "None available %s process to shutdown! ", that.Name)
		return nil
	}
	if pid < 1 {
		log.Info(ctx, "%s has been shutdown now! ", that.Name)
		return nil
	}
	for _, hook := range that.Hooks {
		if err = hook.Stop(ctx, that.Runtime); nil != err {
			log.Error(ctx, err.Error())
		}
	}
	if err = syscall.Kill(pid, syscall.SIGKILL); nil != err {
		return cause.Error(err)
	}
	log.Catch(locker.Unlock())
	log.Catch(os.Remove(that.PIDFd))
	log.Info(ctx, "%s has been shutdown now! ", that.Name)
	return nil
}
