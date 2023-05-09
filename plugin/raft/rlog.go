/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package raft

import (
	"context"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/lni/dragonboat/v4/logger"
)

func init() {
	var _ logger.ILogger = new(rlog)
}

type rlog struct {
	name  string
	level log.Level
	ctx   context.Context
}

func (that *rlog) SetLevel(level logger.LogLevel) {
	switch level {
	case logger.CRITICAL:
		that.level = log.WARN
	case logger.ERROR:
		that.level = log.ERROR
	case logger.WARNING:
		that.level = log.WARN
	case logger.INFO:
		that.level = log.INFO
	case logger.DEBUG:
		that.level = log.DEBUG
	}
}

func (that *rlog) Debugf(format string, args ...interface{}) {
	log.Print(that.ctx, that.level, format, args...)
}

func (that *rlog) Infof(format string, args ...interface{}) {
	log.Print(that.ctx, that.level, format, args...)
}

func (that *rlog) Warningf(format string, args ...interface{}) {
	log.Print(that.ctx, that.level, format, args...)
}

func (that *rlog) Errorf(format string, args ...interface{}) {
	log.Print(that.ctx, that.level, format, args...)
}

func (that *rlog) Panicf(format string, args ...interface{}) {
	log.Print(that.ctx, that.level, format, args...)
}
