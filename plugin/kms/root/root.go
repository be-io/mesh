/*
 * Copyright (c) 2000, 2099, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package root

import (
	"bytes"
	"context"
	"embed"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/prsim"
	"io/ioutil"
)

//go:embed root.yaml
var rootKey embed.FS

func Key(ctx context.Context, yaml codec.Codec) (*prsim.RootKey, error) {
	// Read the root key
	fd, err := rootKey.Open("root.yaml")
	if nil != err {
		return nil, cause.Error(err)
	}
	defer func() { log.Catch(fd.Close()) }()
	keys, err := ioutil.ReadAll(fd)
	if nil != err {
		return nil, cause.Error(err)
	}
	var rk prsim.RootKey
	if _, err = yaml.Decode(bytes.NewBuffer(keys), &rk); nil != err {
		return nil, cause.Error(err)
	}
	return &rk, nil
}
