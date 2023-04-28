/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type DistributeOption struct {
	Set  string `index:"0" json:"set" yaml:"set" xml:"set" comment:""`
	Lang string `index:"5" json:"lang" yaml:"lang" xml:"lang" comment:""`
	Addr string `index:"10" json:"addr" yaml:"addr" xml:"addr" comment:""`
}

type TransformOption struct {
	Schema      string `index:"0" json:"schema" yaml:"schema" xml:"schema" comment:""`
	Lang        string `index:"5" json:"lang" yaml:"lang" xml:"lang" comment:""`
	Set         string `index:"10" json:"set" yaml:"set" xml:"set" comment:""`
	Version     string `index:"15" json:"version" yaml:"version" xml:"version" comment:""`
	Describe    string `index:"20" json:"describe" yaml:"describe" xml:"describe" comment:""`
	MeshVersion string `index:"25" json:"mesh_version" yaml:"mesh_version" xml:"mesh_version" comment:""`
}
