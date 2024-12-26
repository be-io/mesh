/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package boost

import (
	"bytes"
	"context"
	"fmt"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/crypt"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
	"time"
)

func MeshID(ctx context.Context, nodeId string, instId string, name string, crt string) (string, error) {
	keys, err := aware.KMS.ApplyRoot(ctx, &types.KeyCsr{
		CNO:      nodeId,
		PNO:      nodeId,
		Domain:   fmt.Sprintf("%s.%s", nodeId, types.MeshDomain),
		Subject:  name,
		Length:   2048,
		ExpireAt: types.Time(time.Now().AddDate(100, 0, 0)),
		Mail:     fmt.Sprintf("%s@%s", nodeId, types.MeshDomain),
		IsCA:     true,
		CaCert:   "",
		CaKey:    "",
	})
	if nil != err {
		return "", cause.Error(err)
	}
	keySet := types.KeysSet(keys)
	environ := &types.Environ{
		Version:  types.EnvironVersion,
		NodeId:   nodeId,
		InstId:   instId,
		InstName: name,
		RootCrt:  keySet.Get(types.RootCaCrtKey),
		RootKey:  keySet.Get(types.RootCaPrivateKey),
		NodeCrt:  crt,
	}
	buff, err := aware.YAML.Encode(environ)
	if nil != err {
		return "", cause.Error(err)
	}
	rk, err := aware.Singularity.Reveal(ctx)
	if nil != err {
		return "", cause.Error(err)
	}
	pub, err := crypt.PEM.InformKey(rk.PublicKey)
	if nil != err {
		return "", cause.Error(err)
	}
	eff, err := aware.RSA2.Encrypt(ctx, buff.Bytes(), map[string][]byte{string(prsim.PublicKey): pub})
	if nil != err {
		return "", cause.Error(err)
	}
	return crypt.Base64.FormKey(eff, "MESH ID")
}

func FromMeshID(ctx context.Context, meshId string) (*types.Environ, error) {
	mff, err := crypt.Base64.InformKey(meshId)
	if nil != err {
		return nil, cause.Error(err)
	}
	rk, err := aware.Singularity.Reveal(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	privateKey, err := crypt.PEM.InformKey(rk.PrivateKey)
	if nil != err {
		return nil, cause.Error(err)
	}
	eff, err := aware.RSA2.Decrypt(ctx, mff, map[string][]byte{string(prsim.PrivateKey): privateKey})
	if nil != err {
		return nil, cause.Error(err)
	}
	var environ types.Environ
	_, err = aware.YAML.Decode(bytes.NewBuffer(eff), &environ)
	if nil != err {
		return nil, cause.Error(err)
	}
	return &environ, nil
}
