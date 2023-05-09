/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import "github.com/be-io/mesh/client/golang/plugin"

func formatLegacyName(name string) string {
	switch name {
	case "edge-asset", "asset":
		return "asset"
	case "edge-serving", "edge", "serving":
		return "edge"
	case "theta":
		return "theta"
	case "omega":
		return "omega"
	case "gaia-janus", "janus":
		return "janus"
	case "gaia-tensor", "tensor":
		return "tensor"
	case "base-server", "server":
		return "server"
	case "base-client", "client":
		return "client"
	case "gaia-pandora", "pandora":
		return "pandora"
	case "gaia-mesh", "mesh":
		return plugin.Whoami
	case "cube-engine", "cube":
		return "cube"
	default:
		return name
	}
}
