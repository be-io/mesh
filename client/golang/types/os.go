/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type Cluster struct {
	Version string         `index:"0" json:"version" xml:"version" yaml:"version" comment:"Cluster version"`
	Name    string         `index:"1" json:"name" xml:"name" yaml:"name" comment:"Cluster name"`
	Memo    string         `index:"2" json:"memo" xml:"memo" yaml:"memo" comment:"Cluster memo"`
	Icon    string         `index:"3" json:"icon" xml:"icon" yaml:"icon" comment:"Cluster icon"`
	Status  int64          `index:"4" json:"status" xml:"status" yaml:"status" comment:"Cluster status"`
	Hosts   []*ClusterHost `index:"5" json:"hosts" xml:"hosts" yaml:"hosts" comment:"Cluster hosts"`
	Spec    *ClusterSpec   `index:"6" json:"spec" xml:"spec" yaml:"spec" comment:"Cluster spec"`
}

type ClusterHost struct {
	Host   string `index:"0" json:"host" xml:"host" yaml:"host" comment:"Cluster host"`
	PubKey string `index:"1" json:"pub_key" xml:"pub_key" yaml:"pub_key" comment:"Cluster public key"`
}

type ClusterSpec struct {
	CPU       int64 `index:"0" json:"cpu" xml:"cpu" yaml:"cpu" comment:"Cluster cpu spec"`
	GPU       int64 `index:"1" json:"gpu" xml:"gpu" yaml:"gpu" comment:"Cluster gpu spec"`
	Mem       int64 `index:"2" json:"mem" xml:"mem" yaml:"mem" comment:"Cluster mem spec"`
	Disk      int64 `index:"3" json:"disk" xml:"disk" yaml:"disk" comment:"Cluster disk spec"`
	Bandwidth int64 `index:"4" json:"bandwidth" xml:"bandwidth" yaml:"bandwidth" comment:"Cluster bandwidth spec"`
}

type Workspace struct {
	Version string `index:"0" json:"version" xml:"version" yaml:"version" comment:"Cluster version"`
	Name    string `index:"1" json:"name" xml:"name" yaml:"name" comment:"Cluster name"`
	Memo    string `index:"2" json:"memo" xml:"memo" yaml:"memo" comment:"Cluster memo"`
	Icon    string `index:"3" json:"icon" xml:"icon" yaml:"icon" comment:"Cluster icon"`
	Status  int64  `index:"4" json:"status" xml:"status" yaml:"status" comment:"Cluster status"`
	Cluster string `index:"5" json:"cluster" xml:"cluster" yaml:"cluster" comment:"Cluster name"`
}

type OSCharts struct {
	Version string `index:"0" json:"version" xml:"version" yaml:"version" comment:"Chart version"`
	Name    string `index:"1" json:"name" xml:"name" yaml:"name" comment:"Chart name"`
	Memo    string `index:"2" json:"memo" xml:"memo" yaml:"memo" comment:"Chart memo"`
	Config  string `index:"3" json:"config" xml:"config" yaml:"config" comment:"Chart config"`
}

type Operation struct {
}
