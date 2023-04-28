/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package nsqio

import (
	"crypto/tls"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/nsqio/nsq/nsqd"
	"strconv"
	"strings"
	"time"
)

type nsqOption struct {
	Config                   string `json:"plugin.nsq.config" dft:"" usage:"path to config file"`
	LogPrefix                string `json:"plugin.nsq.log.prefix" dft:"[mesh] " usage:"log message prefix"`
	NodeId                   string `json:"plugin.nsq.node.id" dft:"" usage:"unique part for message IDs, (int) in range [0,1024) (default is hash of hostname)"`
	HTTPSAddress             string `json:"plugin.nsq.https.address" dft:"" usage:"<addr>:<port> to listen on for HTTPS clients"`
	HTTPAddress              string `json:"plugin.nsq.http.address" dft:"" usage:"<addr>:<port> to listen on for HTTP clients"`
	TCPAddress               string `json:"plugin.nsq.tcp.address" dft:"" usage:"<addr>:<port> to listen on for TCP clients"`
	AuthHTTPAddresses        string `json:"plugin.nsq.http.auth.address" dft:"" usage:"<addr>:<port> or a full url to query auth server (may be given multiple times)"`
	BroadcastAddress         string `json:"plugin.nsq.broadcast.address" dft:"" usage:"address that will be registered with lookupd (defaults to the OS hostname)"`
	BroadcastTCPPort         string `json:"plugin.nsq.broadcast.tcp.port" dft:"" usage:"TCP port that will be registered with lookupd (defaults to the TCP port that this nsqd is listening on)"`
	BroadcastHTTPPort        string `json:"plugin.nsq.broadcast.http.port" dft:"" usage:"HTTP port that will be registered with lookupd (defaults to the HTTP port that this nsqd is listening on)"`
	LookupdTCPAddrs          string `json:"plugin.nsq.lookupd.tcp.address" dft:"" usage:"lookupd TCP address (may be given multiple times)"`
	HTTPClientConnectTimeout string `json:"plugin.nsq.http.client.connect.timeout" dft:"" usage:"timeout for HTTP connect"`
	HTTPClientRequestTimeout string `json:"plugin.nsq.http.client.request.timeout" dft:"" usage:"timeout for HTTP request"`
	// diskqueue options
	DataPath        string        `json:"plugin.nsq.data.path" dft:"${MESH_HOME}/mesh/nsq/${IPH}" usage:"path to store disk-backed messages"`
	MemQueueSize    int64         `json:"plugin.nsq.mem.queue.size" dft:"64" usage:"number of messages to keep in memory (per topic/channel)"`
	MaxBytesPerFile int64         `json:"plugin.nsq.max.bytes.per.file" dft:"" usage:"number of bytes per diskqueue file before rolling"`
	SyncEvery       int64         `json:"plugin.nsq.sync.every" dft:"8" usage:"number of messages per diskqueue fsync"`
	SyncTimeout     time.Duration `json:"plugin.nsq.sync.timeout" dft:"10s" usage:"duration of time per diskqueue fsync"`

	QueueScanWorkerPoolMax  int `json:"queue-scan-worker-pool-max" dft:"" usage:"max concurrency for checking in-flight and deferred message timeouts"`
	QueueScanSelectionCount int `json:"queue-scan-selection-count" dft:"" usage:"number of channels to check per cycle (every 100ms) for in-flight and deferred timeouts"`
	// msg and command options
	MsgTimeout    time.Duration `json:"plugin.nsq.msg.timeout" dft:"1h" usage:"default duration to wait before auto-requeing a message"`
	MaxMsgTimeout time.Duration `json:"plugin.nsq.msg.timeout.max" dft:"3d" usage:"maximum duration before a message will timeout"`
	MaxMsgSize    int64         `json:"plugin.nsq.msg.size.max" dft:"" usage:"maximum size of a single message in bytes"`
	MaxReqTimeout time.Duration `json:"plugin.nsq.req.timeout.max" dft:"1h" usage:"maximum requeuing timeout for a message"`
	MaxBodySize   int64         `json:"plugin.nsq.body.size.max" dft:"" usage:"maximum size of a single command body"`
	// client overridable configuration options
	MaxHeartbeatInterval   string `json:"max-heartbeat-interval" dft:"" usage:"maximum client configurable duration of time between client heartbeats"`
	MaxRdyCount            int64  `json:"max-rdy-count" dft:"" usage:"maximum RDY count for a client"`
	MaxOutputBufferSize    int64  `json:"max-output-buffer-size" dft:"" usage:"maximum client configurable size (in bytes) for a client output buffer"`
	MaxOutputBufferTimeout string `json:"max-output-buffer-timeout" dft:"" usage:"maximum client configurable duration of time between flushing to a client"`
	MinOutputBufferTimeout string `json:"min-output-buffer-timeout" dft:"" usage:"minimum client configurable duration of time between flushing to a client"`
	OutputBufferTimeout    string `json:"output-buffer-timeout" dft:"" usage:"default duration of time between flushing data to clients"`
	MaxChannelConsumers    int    `json:"max-channel-consumers" dft:"" usage:"maximum channel consumers connection count per nsqd instance (default 0, i.e., unlimited)"`
	// statsd integration options
	StatsdAddress          string `json:"statsd-address" dft:"" usage:"UDP <addr>:<port> of a statsd daemon for pushing stats"`
	StatsdInterval         string `json:"statsd-interval" dft:"" usage:"duration between pushing to statsd"`
	StatsdMemStats         bool   `json:"statsd-mem-stats" dft:"" usage:"toggle sending memory and GC stats to statsd"`
	StatsdPrefix           string `json:"statsd-prefix" dft:"" usage:"prefix used for keys sent to statsd (%s for host replacement)"`
	StatsdUDPPacketSize    int    `json:"statsd-udp-packet-size" dft:"" usage:"the size in bytes of statsd UDP packets"`
	StatsdExcludeEphemeral bool   `json:"statsd-exclude-ephemeral" dft:"" usage:"Skip ephemeral topics and channels when sending stats to statsd"`
	// End to end percentile flags
	E2eProcessingLatencyPercentiles string `json:"e2e-processing-latency-percentile" dft:"" usage:"message processing time percentiles (as float (0, 1.0]) to track (can be specified multiple times or comma separated '1.0,0.99,0.95', default none)"`
	E2EProcessingLatencyWindowTime  string `json:"e2e-processing-latency-window-time" dft:"" usage:"calculate end to end latency quantiles for this duration of time (ie: 60s would only show quantile calculations from the past 60 seconds)"`
	// TLS config
	TLSCert             string `json:"tls-cert" dft:"" usage:"path to certificate file"`
	TLSKey              string `json:"tls-key" dft:"" usage:"path to key file"`
	TLSClientAuthPolicy string `json:"tls-client-auth-policy" dft:"" usage:"client certificate auth policy ('require' or 'require-verify')"`
	TLSRootCAFile       string `json:"tls-root-ca-file" dft:"" usage:"path to certificate authority file"`
	TLSRequired         string `json:"tls-required" dft:"" usage:"require TLS for client connections (true, false, tcp-https)"`
	TLSMinVersion       string `json:"tls-min-version" dft:"" usage:"minimum SSL/TLS version acceptable ('ssl3.0', 'tls1.0', 'tls1.1', or 'tls1.2')"`
	// compression
	DeflateEnabled  bool `json:"deflate" dft:"true" usage:"enable deflate feature negotiation (client compression)"`
	MaxDeflateLevel int  `json:"max-deflate-level" dft:"" usage:"max deflate compression level a client can negotiate (> values == > nsqd CPU usage)"`
	SnappyEnabled   bool `json:"snappy" dft:"true" usage:"enable snappy feature negotiation (client compression)"`
}
type tlsRequiredOption int

func (t *tlsRequiredOption) Set(s string) error {
	s = strings.ToLower(s)
	if s == "tcp-https" {
		*t = nsqd.TLSRequiredExceptHTTP
		return nil
	}
	required, err := strconv.ParseBool(s)
	if required {
		*t = nsqd.TLSRequired
	} else {
		*t = nsqd.TLSNotRequired
	}
	return err
}

func (t *tlsRequiredOption) Get() interface{} { return int(*t) }

func (t *tlsRequiredOption) String() string {
	return strconv.FormatInt(int64(*t), 10)
}

func (t *tlsRequiredOption) IsBoolFlag() bool { return true }

type tlsMinVersionOption uint16

func (t *tlsMinVersionOption) Set(s string) error {
	s = strings.ToLower(s)
	switch s {
	case "":
		return nil
	case "ssl3.0":
		*t = tls.VersionSSL30
	case "tls1.0":
		*t = tls.VersionTLS10
	case "tls1.1":
		*t = tls.VersionTLS11
	case "tls1.2":
		*t = tls.VersionTLS12
	default:
		return cause.Errorf("unknown tlsVersionOption %q", s)
	}
	return nil
}

func (t *tlsMinVersionOption) Get() interface{} { return uint16(*t) }

func (t *tlsMinVersionOption) String() string {
	return strconv.FormatInt(int64(*t), 10)
}
