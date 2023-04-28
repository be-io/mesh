/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package nsqio

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/log"
	"os"
	"syscall"
	"time"
)

type DirLock struct {
	dir string
	fd  *os.File
}

func New(dir string) *DirLock {
	if "" == dir {
		dir, _ = os.Getwd()
	}
	return &DirLock{
		dir: dir,
	}
}

func (that *DirLock) Lock() error {
	f, err := os.Open(that.dir)
	if err != nil {
		return err
	}
	that.fd = f
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return fmt.Errorf("cannot flock directory %s - %s (possibly in use by another instance of nsqd)", that.dir, err)
	}
	return nil
}

func (that *DirLock) Unlock() error {
	defer func() { log.Catch(that.fd.Close()) }()
	return syscall.Flock(int(that.fd.Fd()), syscall.LOCK_UN)
}

func (that *DirLock) WaitLock(ctx context.Context, duration time.Duration, fn func()) {
	ticker := time.NewTicker(time.Second * 5)
	timeout := time.NewTicker(duration)
	defer func() {
		ticker.Stop()
		timeout.Stop()
	}()
	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				log.Warn(ctx, "NSQ wait ticker has been closed. ")
				return
			}
			if err := that.Lock(); nil != err {
				log.Warn(ctx, "Try to startup nsq, %s", err.Error())
				continue
			}
			log.Catch(that.Unlock())
			fn()
			log.Info(ctx, "NSQ wait startup succeed. ")
			return
		case <-timeout.C:
			log.Warn(ctx, "NSQ wait startup has been timeout. ")
			return
		}
	}
}
