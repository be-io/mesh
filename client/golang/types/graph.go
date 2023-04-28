/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type Predicate string

const (
	Has    Predicate = "#"
	Direct Predicate = "->"
)

type Vertex struct {
	ID    string            `index:"0" json:"id" xml:"id" yaml:"id" comment:"Vertex identity"`
	Label string            `index:"1" json:"label" xml:"label" yaml:"label" comment:"Vertex label"`
	Attrs map[string]string `index:"2" json:"attrs" xml:"attrs" yaml:"attrs" comment:"Vertex attributes"`
	Raw   []byte            `index:"3" json:"raw" xml:"raw" yaml:"raw" comment:"Vertex raw data"`
}

type Side struct {
	Src   string            `index:"0" json:"src" xml:"src" yaml:"src" comment:"Vertex source identity"`
	Dst   string            `index:"1" json:"dst" xml:"dst" yaml:"dst" comment:"Vertex destination identity"`
	Label string            `index:"1" json:"label" xml:"label" yaml:"label" comment:"Vertex label"`
	Attrs map[string]string `index:"2" json:"attrs" xml:"attrs" yaml:"attrs" comment:"Vertex attributes"`
	Raw   []byte            `index:"3" json:"raw" xml:"raw" yaml:"raw" comment:"Vertex raw data"`
}

type VertexLabel struct {
}

type Tuple struct {
	S     string `index:"0" json:"s" xml:"s" yaml:"s" comment:"Subject"`
	P     string `index:"1" json:"p" xml:"p" yaml:"p" comment:"Predicate"`
	Depth int64  `index:"2" json:"dep" xml:"dep" yaml:"dep" comment:"Depth"`
}

// Triple (S, P, O)
type Triple struct {
	S string `index:"0" json:"s" xml:"s" yaml:"s" comment:"Subject"`
	P string `index:"1" json:"p" xml:"p" yaml:"p" comment:"Predicate"`
	O string `index:"2" json:"o" xml:"o" yaml:"o" comment:"Object"`
}

// Quad (S, P, O, L{})
type Quad struct {
	S string `index:"0" json:"s" xml:"s" yaml:"s" comment:"Subject"`
	P string `index:"1" json:"p" xml:"p" yaml:"p" comment:"Predicate"`
	O string `index:"2" json:"o" xml:"o" yaml:"o" comment:"Object"`
	L string `index:"3" json:"l" xml:"l" yaml:"l" comment:"Label"`

	Attrs map[string]string `index:"5" json:"attrs" xml:"attrs" yaml:"attrs" comment:"Vertex attributes"`
	Raw   []byte            `index:"6" json:"raw" xml:"raw" yaml:"raw" comment:"Vertex raw data"`
}

// VertexQuad S#L=O
func VertexQuad(vertex string, attrs map[string]string, raw []byte) *Quad {
	return &Quad{S: vertex, P: string(Has), O: vertex, Attrs: attrs, Raw: raw}
}

// LinkQuad S#L->O
func LinkQuad(src string, dst string, attrs map[string]string) *Quad {
	return &Quad{S: src, P: string(Direct), O: dst, Attrs: attrs}
}
