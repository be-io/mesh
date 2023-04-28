/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type MeshQL struct {
	Name string `index:"0" json:"name" xml:"name" yaml:"name" comment:"Name"`
	Expr string `index:"1" json:"expr" xml:"expr" yaml:"expr" comment:"Expr"`
}

type Vec1 struct {
	Name string `index:"0" json:"name" xml:"name" yaml:"name" comment:"Name"`
	X    string `index:"1" json:"x" xml:"x" yaml:"x" comment:"X"`
}

type Vec2 struct {
	Name string `index:"0" json:"name" xml:"name" yaml:"name" comment:"Name"`
	X    string `index:"1" json:"x" xml:"x" yaml:"x" comment:"X"`
	Y    string `index:"2" json:"y" xml:"y" yaml:"y" comment:"Y"`
}

type Vertex struct {
	Name  string            `index:"0" json:"name" xml:"name" yaml:"name" comment:"Name"`
	Group string            `index:"1" json:"group" xml:"group" yaml:"group" comment:"Group"`
	Attrs map[string]string `index:"2" json:"attrs" xml:"attrs" yaml:"attrs" comment:"Attributes"`
	X     string            `index:"3" json:"x" xml:"x" yaml:"x" comment:"X"`
}

type Side struct {
}

type Cursor struct {
	Name  string `index:"0" json:"name" xml:"name" yaml:"name" comment:"Name"`
	Depth int64  `index:"1" json:"depth" xml:"depth" yaml:"depth" comment:"Depth"`
	X     string `index:"2" json:"x" xml:"x" yaml:"x" comment:"X"`
}

type Quad struct {
	Name      string `index:"0" json:"name" xml:"name" yaml:"name" comment:"Name"`
	Subject   string `index:"1" json:"subject" xml:"subject" yaml:"subject" comment:"Subject"`
	Predicate string `index:"2" json:"predicate" xml:"predicate" yaml:"predicate" comment:"Predicate"`
	Object    string `index:"3" json:"object" xml:"object" yaml:"object" comment:"Object"`
	Label     string `index:"4" json:"label" xml:"label" yaml:"label" comment:"Label"`
}
