/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"github.com/be-io/mesh/client/golang/macro"
)

var ServiceProxy = new(serviceProxy)

type serviceProxy struct {
}

func (that *serviceProxy) Reference(rtt *macro.Rtt) macro.Caller {
	return &ReferenceHandler{MPI: &macro.MPIAnnotation{Meta: rtt}}
}

func (that *serviceProxy) Proxy(mpi macro.MPI) macro.Caller {
	return &ReferenceHandler{MPI: mpi}
}
