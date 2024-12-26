/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package schema

import (
	"context"
	"github.com/opendatav/mesh/client/golang/cause"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
)

const my = "__my__"

var Runtime Schema = &meshableSchema{manifests: map[string]map[string][]*Kind{}}

var ISchema = (*Schema)(nil)

type Schema interface {

	// Scan will scan and new type reference
	Scan(ctx context.Context, root string) (map[string][]*Kind, error)

	// Refine will scan and define type reference.
	Refine(schema string) (interface{}, error)

	// Define a schema into runtime
	Define(ctx context.Context, reference interface{}) (string, error)

	// Import the mesh interface schemas.
	Import(ctx context.Context, types map[string]map[string][]*Kind) error

	// Export the schema with codec
	Export(ctx context.Context) (map[string]map[string][]*Kind, error)
}

type meshableManifest struct {
	Instance string
	Types    map[string][]*Kind
}

func (that *meshableManifest) Refine(schema string) (interface{}, bool) {
	return that.Types, true
}

type meshableSchema struct {
	manifests map[string]map[string][]*Kind
}

func (that *meshableSchema) Scan(ctx context.Context, root string) (map[string][]*Kind, error) {
	if nil != that.manifests[my] {
		return that.manifests[my], nil
	}
	sets := token.NewFileSet()
	pkg, err := build.Import(root, "", build.ImportComment)
	if nil != err {
		return nil, cause.Error(err)
	}
	fast, err := parser.ParseFile(sets, pkg.Dir, nil, parser.ParseComments)
	if nil != err {
		return nil, cause.Error(err)
	}
	for _, dcl := range fast.Decls {
		switch decl := dcl.(type) {
		case *ast.GenDecl:
			if len(decl.Specs) < 0 {

			}
		}
	}
	return that.manifests[my], nil
}

func (that *meshableSchema) Refine(schema string) (interface{}, error) {
	return nil, nil
}

func (that *meshableSchema) Define(ctx context.Context, reference interface{}) (string, error) {
	return "", nil
}

func (that *meshableSchema) Import(ctx context.Context, types map[string]map[string][]*Kind) error {
	return nil
}

func (that *meshableSchema) Export(ctx context.Context) (map[string]map[string][]*Kind, error) {
	return that.manifests, nil
}
