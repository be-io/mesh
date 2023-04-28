/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package schema

type Kind struct {
	Pkg       string                       `json:"0"`
	Name      string                       `json:"1"`
	Imports   []string                     `json:"2"`
	Macros    map[string]map[string]string `json:"3"`
	Comments  []string                     `json:"4"`
	Modifier  int                          `json:"5"`
	Variables []*Variable                  `json:"6"`
	Methods   []*Method                    `json:"7"`
	Supers    []*Kind                      `json:"8"`
	Traits    []*Kind                      `json:"9"`
	Signature string                       `json:"10"`
}

type Variable struct {
	Name     string                       `json:"0"`
	Macros   map[string]map[string]string `json:"1"`
	Comments []string                     `json:"2"`
	Modifier int                          `json:"3"`
	Kind     *Kind                        `json:"4"`
	Value    string                       `json:"5"`
}

type Method struct {
	Name       string                       `json:"0"`
	Macros     map[string]map[string]string `json:"1"`
	Comments   []string                     `json:"2"`
	Modifier   int                          `json:"3"`
	Parameters []*Parameter                 `json:"4"`
	Returns    []*Return                    `json:"5"`
	Causes     []*Throwable                 `json:"6"`
}

type Parameter struct {
	Name     string                       `json:"0"`
	Macros   map[string]map[string]string `json:"1"`
	Comments []string                     `json:"2"`
	Kind     *Kind                        `json:"3"`
	Value    string                       `json:"4"`
}

type Return struct {
	Name     string   `json:"0"`
	Comments []string `json:"1"`
	Kind     *Kind    `json:"2"`
}

type Throwable struct {
	Name     string   `json:"0"`
	Comments []string `json:"1"`
	Kind     *Kind    `json:"2"`
}

type Set struct {
	Name        string
	Version     string
	Describe    string
	MeshVersion string
}

type Tree struct {
	Root      string
	Childrens []*Tree
}
