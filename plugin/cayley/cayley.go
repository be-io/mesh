/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cayley

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/dsa"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/plugin"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/plugin/cayley/graph"
	"path/filepath"
)

func init() {
	plugin.Provide(stores)
}

const Name = "cayley"

var stores = &cayley{graphs: dsa.NewStringMap[*graph.Graph]()}

type cayley struct {
	Home   string `json:"cayley.home" dft:"${MESH_HOME}/mesh/cayley" usage:"Directory to store cayley bolt data."`
	graphs dsa.Map[string, *graph.Graph]
}

func (that *cayley) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Name: plugin.GDB, WaitAny: true, Flags: cayley{}, Create: func() plugin.Plugin {
		return that
	}}
}

func (that *cayley) Start(ctx context.Context, runtime plugin.Runtime) {
	log.Catch(runtime.Parse(that))
	if err := tool.MakeDir(that.Home); nil != err {
		log.Error(ctx, err.Error())
	}
}

func (that *cayley) Stop(ctx context.Context, runtime plugin.Runtime) {

}

func (that *cayley) OpenGraph(ctx context.Context, name string) (*graph.Graph, error) {
	return that.graphs.PutIfe(name, func(k string) (*graph.Graph, error) {
		if err := tool.MakeDir(that.Home); nil != err {
			return nil, cause.Error(err)
		}
		return graph.NewGraph(filepath.Join(that.Home, name)), nil
	})
}
