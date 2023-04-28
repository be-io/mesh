/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
)

func init() {
	var _ prsim.Listener = new(PRSIListener)
	macro.Provide(prsim.IListener, new(PRSIListener))
}

const PRSListener = "mesh.plugin.prs.listener"

// PRSIListener main listener.
type PRSIListener struct {
}

func (that *PRSIListener) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.prs.listener"}
}

func (that *PRSIListener) Btt() []*macro.Btt {
	return nil
}

func (that *PRSIListener) Listen(ctx context.Context, event *types.Event) error {
	for _, listener := range macro.Load(prsim.IListener).List() {
		if lister, ok := listener.(prsim.Listener); ok && event.Binding.Match(lister.Btt()...) {
			if err := lister.Listen(ctx, event); nil != err {
				log.Error(ctx, "Listen event with %s/%s, %s. ", event.Binding.Topic, event.Binding.Code, err.Error())
			}
		}
	}
	log.Debug(ctx, "Listen event published with %s/%s. ", event.Binding.Topic, event.Binding.Code)
	return nil
}
