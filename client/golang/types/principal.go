/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import "fmt"

func init() {
	var _ fmt.Stringer = new(Principal)
	var _ fmt.Stringer = new(Location)
}

type Principal struct {
	NodeId string `index:"0" json:"node_id" xml:"node_id" yaml:"node_id" comment:"node_id"`
	InstId string `index:"5" json:"inst_id" xml:"inst_id" yaml:"inst_id" comment:"inst_id"`
}

type Location struct {
	Principal
	IP   string `index:"10" json:"ip" xml:"ip" yaml:"ip" comment:"ip"`
	Port string `index:"15" json:"port" xml:"port" yaml:"port" comment:"port"`
	Host string `index:"20" json:"host" xml:"host" yaml:"host" comment:"host"`
	Name string `index:"25" json:"name" xml:"name" yaml:"name" comment:"name"`
}

func (that *Principal) String() string {
	return fmt.Sprintf("{\"node_id\":\"%s\",\"inst_id\":\"%s\"}", that.NodeId, that.InstId)
}

func (that *Location) String() string {
	return fmt.Sprintf("{\"node_id\":\"%s\",\"inst_id\":\"%s\",\"ip\":\"%s\",\"host\":\"%s\",\"port\":\"%s\",\"name\":\"%s\"}",
		that.NodeId, that.InstId, that.IP, that.Host, that.Port, that.Name)
}
