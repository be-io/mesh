/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cmd

import (
	"bytes"
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"os"
	"os/exec"
	"strings"
)

func init() {
	var _ prsim.EndpointSticker[*string, *string] = new(meshEndpoint)
	macro.Provide(prsim.IEndpointSticker, new(meshEndpoint))
}

type meshEndpoint struct {
}

func (that *meshEndpoint) Att() *macro.Att {
	return &macro.Att{Name: "mesh.dot.exe"}
}

func (that *meshEndpoint) Rtt() *macro.Rtt {
	return &macro.Rtt{Name: "mesh.dot.exe"}
}

func (that *meshEndpoint) I() *string {
	var x string
	return &x
}

func (that *meshEndpoint) O() *string {
	var x string
	return &x
}

func (that *meshEndpoint) Stick(ctx context.Context, varg *string) (*string, error) {
	writer := &bytes.Buffer{}
	cmd := exec.Command(os.Args[0], strings.Split(*varg, " ")...)
	cmd.Stdout = writer
	cmd.Stderr = writer
	if err := cmd.Run(); nil != err {
		return nil, cause.Error(err)
	}
	v := writer.String()
	return &v, nil
}
