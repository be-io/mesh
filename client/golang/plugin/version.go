/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package plugin

var (
	Whoami = "mesh"
	Format = "yaml"
	Config = "{}"
	Mode   = 0
	Debug  = false
)

const (
	SERVER   runMode = 1
	NODE     runMode = 2
	OPERATOR runMode = 4
	PANEL    runMode = 8
	PROXY    runMode = 16
	INC      runMode = 33
)

const Banner = `
	 __  __           _     
	|  \/  | ___  ___| |__  
	| |\/| |/ _ \/ __| '_ \      A lightweight, distributed, relational
	| |  | |  __/\__ \ | | |     network architecture for MPC
	|_|  |_|\___||___/_| |_|

	(v%s, build %s)

	
`

const Usage = `
Usage:  mesh [COMMAND] [OPTIONS] [ARG...]

A lightweight, distributed, relational network architecture for MPC.

Options:
      --server    Run mesh as mesh server.
      --operator  Run mesh as k8s operator.
      --node      Run mesh as mesh node.
      --panel     Run mesh as mesh panel.
      --proxy     Run mesh as mesh proxy.
      --inc       Run mesh as mesh inc mode.
  -c, --config    Mesh configuration input, it can be url/path/body, default is {}.
  -f, --format    Mesh configuration format, it can be json/yaml, default is yaml.
  -d, --debug     Set the run mode to debug with more information.
  -h, --help      Print more information on a command.
  -v, --version   Print version information and quit.

Commands:
  inspect        Return low-level information on mesh plugin objects.
  start          Start mesh with the run mode(server/operator/node/panel).
  stop           Stop the mesh process if available.
  status         Display a live stream of mesh cluster resource usage statistics.
  dump           Dump the mesh cluster metadata tables.
  exec           Execute the mesh actor, Such mesh exec actor -i '{}'.

Run 'mesh COMMAND --help' for more information on a command.
`

type runMode int

func (that runMode) Enable(enable bool) {
	if enable {
		Mode = int(that) | Mode
	}
}

func (that runMode) Match() bool {
	return (int(that) & Mode) == int(that)
}
