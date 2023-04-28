/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package tool

import (
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/panjf2000/ants/v2"
	"runtime"
)

var SharedRoutines = new(macro.Once[*ants.Pool]).With(func() *ants.Pool {
	pool, err := ants.NewPool(runtime.NumCPU()*100, ants.WithPanicHandler(func(err interface{}) {
		log.Error0("%v", err)
	}))
	if nil != err {
		panic(err.Error())
	}
	return pool
})
