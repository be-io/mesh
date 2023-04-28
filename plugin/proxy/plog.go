/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"bytes"
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"github.com/traefik/traefik/v3/pkg/middlewares/accesslog"
	"strconv"
	"time"
)

func init() {
	var _ prsim.Hodor = logger
	logrus.SetReportCaller(true)
	macro.Provide(prsim.IHodor, logger)
}

const (
	defaultValue = "-"
)

var logger = new(plog)

var (
	MeshUrn           = prsim.MeshUrn.Keys("request_")
	MeshTraceId       = prsim.MeshTraceId.Keys("request_")
	MeshSpanId        = prsim.MeshSpanId.Keys("request_")
	MeshFromInstId    = prsim.MeshFromInstId.Keys("request_")
	MeshFromNodeId    = prsim.MeshFromNodeId.Keys("request_")
	MeshIncomingHost  = prsim.MeshIncomingHost.Keys("request_")
	MeshOutgoingHost  = prsim.MeshOutgoingHost.Keys("request_")
	MeshIncomingProxy = prsim.MeshIncomingProxy.Keys("request_")
	MeshOutgoingProxy = prsim.MeshOutgoingProxy.Keys("request_")
	MeshSubset        = prsim.MeshSubset.Keys("request_")
)

type plog struct {
}

func (that *plog) Att() *macro.Att {
	return &macro.Att{Name: "plugin.proxy.plog.hodor"}
}

func (that *plog) X() string {
	return ""
}

func (that *plog) Stats(ctx context.Context, features []string) (map[string]string, error) {
	return nil, nil
}

func (that *plog) Debug(ctx context.Context, features map[string]string) error {
	if nil != features && "" != features["log.level"] {
		level, err := strconv.Atoi(features["log.level"])
		if nil != err {
			return cause.Error(err)
		}
		if log.ALL.Is(level) || log.STACK.Is(level) || log.DEBUG.Is(level) {
			zlog.Logger.Level(zerolog.DebugLevel)
			return nil
		}
		if log.INFO.Is(level) {
			zlog.Logger.Level(zerolog.InfoLevel)
			return nil
		}
		if log.WARN.Is(level) {
			zlog.Logger.Level(zerolog.WarnLevel)
			return nil
		}
		if log.ERROR.Is(level) {
			zlog.Logger.Level(zerolog.ErrorLevel)
			return nil
		}
		if log.FATAL.Is(level) {
			zlog.Logger.Level(zerolog.FatalLevel)
		}
	}
	return nil
}

func (that *plog) Format(entry *logrus.Entry) ([]byte, error) {
	buff := &bytes.Buffer{}
	timestamp := defaultValue
	if v, ok := entry.Data[accesslog.StartUTC]; ok {
		timestamp = v.(time.Time).Local().Format(log.DateFormat)
	} else if v, ok := entry.Data[accesslog.StartLocal]; ok {
		timestamp = v.(time.Time).Local().Format(log.DateFormat)
	}

	var elapsedMillis int64
	if v, ok := entry.Data[accesslog.Duration]; ok {
		elapsedMillis = v.(time.Duration).Nanoseconds() / 1000000
	}
	_, err := fmt.Fprintf(buff, "%s,%s,%s,%v,%s,%s,%s,%s,%v,%v,%dms,%s,%s,%s,%s,%s,%s,%s,%s%s,%s,%s,%v,%s,%s,%s,%s\n",
		timestamp,
		that.toLog(entry.Data, accesslog.ClientHost),
		that.toLog(entry.Data, accesslog.ClientUsername),
		that.toLog(entry.Data, accesslog.OriginStatus),
		that.toLog(entry.Data, accesslog.RequestProtocol),
		that.toLog(entry.Data, accesslog.RequestMethod),
		that.toLog(entry.Data, accesslog.ServiceURL),
		that.toLog(entry.Data, accesslog.RouterName),
		that.toLog(entry.Data, accesslog.OriginContentSize),
		that.toLog(entry.Data, accesslog.RetryAttempts),
		elapsedMillis,
		that.toLog(entry.Data, accesslog.TLSVersion),
		that.toLog(entry.Data, accesslog.TLSCipher),
		that.toLog(entry.Data, MeshTraceId...),
		that.toLog(entry.Data, MeshSpanId...),
		that.toLog(entry.Data, MeshFromInstId...),
		that.toLog(entry.Data, MeshFromNodeId...),
		that.toLog(entry.Data, MeshSubset...),
		that.toLog(entry.Data, append(MeshUrn, accesslog.RequestHost)...),
		that.toLog(entry.Data, accesslog.RequestPath),
		that.toLog(entry.Data, accesslog.RequestRefererHeader),
		that.toLog(entry.Data, accesslog.RequestUserAgentHeader),
		that.toLog(entry.Data, accesslog.RequestCount),
		that.toLog(entry.Data, MeshIncomingHost...),
		that.toLog(entry.Data, MeshOutgoingHost...),
		that.toLog(entry.Data, MeshIncomingProxy...),
		that.toLog(entry.Data, MeshOutgoingProxy...),
	)
	return buff.Bytes(), err
}

func (that *plog) toLog(fields logrus.Fields, keys ...string) interface{} {
	for _, key := range keys {
		if v, ok := fields[key]; ok {
			if nil == v || "" == v {
				continue
			}
			switch val := v.(type) {
			case string:
				return tool.Anyone(val, defaultValue)
			case fmt.Stringer:
				return tool.Anyone(val.String(), defaultValue)
			default:
				return v
			}
		}
	}
	return defaultValue
}
