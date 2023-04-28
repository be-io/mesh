/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package system

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"os"
	"strconv"
	"strings"
)

func init() {
	var _ prsim.Builtin = new(MeshBuiltin)
	macro.Provide(prsim.IBuiltin, new(MeshBuiltin))
}

var _ prsim.Builtin = new(MeshBuiltin)

// MeshBuiltin
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type MeshBuiltin struct {
}

func (that *MeshBuiltin) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}

func (that *MeshBuiltin) Doc(ctx context.Context, name string, formatter string) (string, error) {
	return "", nil
}

func (that *MeshBuiltin) Version(ctx context.Context) (*types.Versions, error) {
	mv := os.Getenv(strings.ToUpper(fmt.Sprintf("%s_version", tool.Name.Get())))
	version := tool.Anyone(mv, prsim.Version)
	return &types.Versions{
		Version: version,
		Infos: map[string]string{
			fmt.Sprintf("%s.commit_id", tool.Name.Get()): prsim.CommitID,
			fmt.Sprintf("%s.os", tool.Name.Get()):        prsim.GOOS,
			fmt.Sprintf("%s.arch", tool.Name.Get()):      prsim.GOARCH,
			fmt.Sprintf("%s.version", tool.Name.Get()):   version,
		},
	}, nil
}

func (that *MeshBuiltin) Debug(ctx context.Context, features map[string]string) error {
	if nil != features && "" != features["log.level"] {
		level, err := strconv.Atoi(features["log.level"])
		if nil != err {
			return cause.Error(err)
		}
		log.SetLevel(log.Level(level))
	}
	if nil != features && "" != features["mesh.mode"] {
		mode, err := strconv.Atoi(features["mesh.mode"])
		if nil != err {
			return cause.Error(err)
		}
		macro.SetMode(macro.Modes(mode))
	}
	for _, door := range macro.Load(prsim.IHodor).List() {
		if hodor, ok := door.(prsim.Hodor); ok {
			if err := hodor.Debug(ctx, features); nil != err {
				return cause.Error(err)
			}
		}
	}
	return nil
}

func (that *MeshBuiltin) Stats(ctx context.Context, features []string) (map[string]string, error) {
	indies := map[string]string{}
	for _, door := range macro.Load(prsim.IHodor).List() {
		if hodor, ok := door.(prsim.Hodor); ok {
			stats, err := hodor.Stats(ctx, features)
			if nil != err {
				return nil, cause.Error(err)
			}
			for key, value := range stats {
				indies[key] = value
			}
		}
	}
	return indies, nil
}

func (that *MeshBuiltin) Fallback(ctx context.Context) error {
	mtx := mpc.ContextWith(ctx)
	if strings.Index(mtx.GetUrn(), "404.") == 0 {
		return cause.Errorable(cause.NoProvider)
	}
	if strings.Index(mtx.GetUrn(), "503.") == 0 {
		return cause.Errorable(cause.ServiceUnavailable)
	}
	if strings.Index(mtx.GetUrn(), "618.") == 0 {
		return cause.Errorable(cause.NetNotWeave)
	}
	return cause.Errorable(cause.Unknown)
}
