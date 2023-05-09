/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/system"
	mtypes "github.com/be-io/mesh/client/golang/types"
	"os"
	"strings"
	"time"
)

func InspectVersion(ctx prsim.Context, nodeId string) (*mtypes.Versions, error) {
	key := fmt.Sprintf("mesh.plugin.dynamic.version.%s", strings.ToLower(nodeId))
	version := &mtypes.Versions{}
	if err := system.GetWithCache(ctx, aware.Cache, key, version); nil != err {
		log.Error(ctx, err.Error())
	}
	if "" != version.Version {
		return version, nil
	}
	if strings.Contains(os.Getenv("MESH_VERSION"), nodeId) {
		var versions map[string]string
		_, err := aware.Codec.DecodeString(os.Getenv("MESH_VERSION"), &versions)
		if nil != err {
			log.Error(ctx, err.Error())
		}
		if nil != versions && "" != versions[nodeId] {
			return &mtypes.Versions{Version: versions[nodeId]}, nil
		}
	}
	mtx := ctx.Resume(ctx)
	mtx.SetAttribute(mpc.TimeoutKey, time.Second*3)
	mtx.GetPrincipals().Push(&mtypes.Principal{NodeId: nodeId})
	defer func() { mtx.GetPrincipals().Pop() }()
	version, err := aware.RemoteNet.Version(mtx, nodeId)
	if nil != err {
		return nil, cause.Errorf("Check node version, %s", err.Error())
	}
	if "" == version.Version {
		return nil, cause.Errorf("Check node version, version=%s", version.Version)
	}
	log.Info(ctx, "Check node version, %s=%s", nodeId, version.Version)
	if err = system.PutWithCache(ctx, aware.Cache, key, version, time.Hour*3); nil != err {
		log.Error(ctx, err.Error())
	}
	return version, nil
}
