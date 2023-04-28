/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cmd

import (
	_ "github.com/be-io/mesh/client/golang/proxy"
	_ "github.com/be-io/mesh/container/node"
	_ "github.com/be-io/mesh/container/operator"
	_ "github.com/be-io/mesh/container/panel"
	_ "github.com/be-io/mesh/container/proxy"
	_ "github.com/be-io/mesh/container/server"
	_ "github.com/be-io/mesh/proxy"
	_ "github.com/be-io/mesh/ptp"
)
