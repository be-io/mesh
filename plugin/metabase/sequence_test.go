/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package metabase

import (
	"fmt"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/plugin"
	"testing"
)

func TestInsertIfAbsent(t *testing.T) {
	ctx := macro.Context()
	dsn := "root:@tcp(127.0.0.1:3306)/mesh"
	container := plugin.LoadC("metabase")
	container.Start(ctx, fmt.Sprintf("--dsn=%s", dsn))
	defer container.Stop(ctx)

	for index := 0; index < 100; index++ {
		id, err := ss.Next(macro.Context(), "INST_ID", 8)
		if nil != err {
			t.Error(err)
			return
		}
		t.Log(id)
	}
}

func TestSequenceConcurrency(t *testing.T) {
	metabase.DSN = "root:@tcp(127.0.0.1:3306)/mesh"
	var sequences []chan []string
	var kinds []string
	for index := 0; index < 20; index++ {
		sequences = append(sequences, make(chan []string, 1))
		if index < 10 {
			kinds = append(kinds, "INST_ID")
		} else {
			kinds = append(kinds, "ASSET_ID")
		}
	}
	for index := 0; index < 20; index++ {
		s := sequences[index]
		kind := kinds[index]
		go func() {
			defer func() { close(s) }()
			for idx := 0; idx < 20; idx++ {
				id, err := ss.Next(macro.Context(), kind, 8)
				if nil != err {
					t.Error(err)
					return
				}
				t.Log(kind, id)
				s <- []string{kind, id}
			}
		}()
	}
	keys := map[string]string{}
	for _, s := range sequences {
		for x, ok := []string{}, true; ok; x, ok = <-s {
			if !ok {
				break
			}
			if len(x) > 0 {
				keys[fmt.Sprintf("%s-%s", x[0], x[1])] = x[1]
				continue
			}
		}
	}
	t.Log(len(keys))
}
