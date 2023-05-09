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

type Cursor struct {
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
	S string            `index:"0" json:"s" xml:"s" yaml:"s" comment:"Subject"`
	P string            `index:"1" json:"p" xml:"p" yaml:"p" comment:"Predicate"`
	O string            `index:"2" json:"o" xml:"o" yaml:"o" comment:"Object"`
	A map[string]string `index:"3" json:"l" xml:"l" yaml:"l" comment:"Attribute"`
}

// VertexQuad S#L=O
func VertexQuad(vertex string, attrs map[string]string) *Quad {
	return &Quad{S: vertex, P: string(Has), O: vertex, A: attrs}
}

// LinkQuad S#L->O
func LinkQuad(src string, dst string, attrs map[string]string) *Quad {
	return &Quad{S: src, P: string(Direct), O: dst, A: attrs}
}
