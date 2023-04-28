/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package main

import (
	"bytes"
	"context"
	"embed"
	"flag"
	"fmt"
	"github.com/be-io/mesh/client/golang/log"
	ark "github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/tool"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/imports"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"
)

//go:embed mp*.*
var home embed.FS

func main() {
	var include, exclude, module, super, macro, mpx, dist string
	flag.StringVar(&mpx, "p", "github.com/be-io/mesh/client/golang/proxy", "Go mesh program interface proxy path.")
	flag.StringVar(&macro, "r", "github.com/be-io/mesh/client/golang/macro", "Go macro path.")
	flag.StringVar(&super, "s", "", "Go mesh program interface proxy path.")
	flag.StringVar(&module, "m", "", "Go module of the proxied path.")
	flag.StringVar(&include, "i", "", "Include bundle, if appointed will proxy the bundle only.")
	flag.StringVar(&exclude, "e", "", "Exclude bundle, if appointed will not proxy the bundle.")
	flag.StringVar(&dist, "d", "proxy", "Proxy source file dist path.")
	flag.Parse()
	ctx := ark.Context()
	pwd, err := os.Getwd()
	if nil != err {
		log.Error(ctx, err.Error())
		return
	}
	proxy := &meshable{
		MPX:          mpx,
		Macro:        macro,
		Src:          pwd,
		Dst:          filepath.Join(pwd, dist),
		Module:       module,
		Super:        super,
		Include:      strings.Split(include, ","),
		Exclude:      strings.Split(exclude, ","),
		declarations: map[string]map[string]*Declaration{}}
	proxy.Proxy(ctx)
}

// Kind is ast.TypeSpec with the associated comment map.
type Kind struct {
	Package     *ast.Package
	Suckage     *ast.Package
	Super       *ast.TypeSpec
	Type        *ast.TypeSpec
	File        *ast.File
	SuFile      *ast.File
	Files       *token.FileSet
	Comment     *ast.CommentGroup
	Methods     []*Method
	Imports     []string
	Attachments map[string]string
	Annotations map[string]map[string]string
}

// Method represents a function signature.
type Method struct {
	Name        string
	Params      []*Param
	Res         []*Param
	Comments    string
	Attachments map[string]string
	Annotations map[string]map[string]string
}

// Param represents a parameter in a function or method signature.
type Param struct {
	Name string
	Type string
}

// Declaration include the comments.
type Declaration struct {
	Name     string
	Package  *ast.Package
	File     *ast.File
	Decl     ast.Node
	Comments *ast.CommentGroup
}

type meshable struct {
	MPX          string
	Macro        string
	Src          string
	Dst          string
	Module       string
	Super        string
	Include      []string
	Exclude      []string
	declarations map[string]map[string]*Declaration
}

// Proxy is static proxy.
func (that *meshable) Proxy(ctx context.Context) {
	packages := map[string]*ast.Package{}
	var dirs []string
	dirs = append(dirs, that.Src)
	for len(dirs) > 0 {
		dir := dirs[0]
		dirs = dirs[1:]
		if !that.Filter(dir) {
			continue
		}
		entries, err := os.ReadDir(dir)
		if nil != err {
			log.Error(ctx, err.Error())
			return
		}
		for _, entry := range entries {
			if entry.IsDir() {
				dirs = append(dirs, fmt.Sprintf("%s%s%s", dir, string(filepath.Separator), entry.Name()))
			}
		}
		ps, err := parser.ParseDir(token.NewFileSet(), dir, func(info os.FileInfo) bool {
			return !strings.Contains(info.Name(), "_test.go") && that.Filter(fmt.Sprintf("%s%s%s", dir, string(filepath.Separator), info.Name()))
		}, parser.ParseComments)
		if nil != err {
			log.Error(ctx, err.Error())
		}
		for _, p := range ps {
			packages[dir] = p
		}
	}

	// share one sets across the whole package
	sets := token.NewFileSet()
	types := map[string]map[string]*Kind{}
	for dir, pkg := range packages {
		if nil == types[pkg.Name] {
			types[pkg.Name] = map[string]*Kind{}
		}
		kinds, err := that.ParseTypes(ctx, dir, pkg, sets)
		if nil != err {
			log.Error(ctx, err.Error())
			return
		}
		for name, kind := range kinds {
			types[pkg.Name][name] = kind
		}
	}
	if err := that.StaticDynamic(ctx, types, that.Dst); nil != err {
		log.Error0(err.Error())
	}
}

// Filter true include, false exclude.
func (that *meshable) Filter(dir string) bool {
	if that.Src == dir {
		return true
	}
	path := strings.Join(strings.Split(strings.ReplaceAll(dir, that.Src, ""), string(filepath.Separator)), "/")
	for _, ex := range that.Exclude {
		if "" != ex && strings.Index(path, fmt.Sprintf("/%s", ex)) == 0 {
			return false
		}
	}
	for _, in := range that.Include {
		if "" != in && strings.Index(path, fmt.Sprintf("/%s", in)) == 0 || strings.Index(fmt.Sprintf("/%s", in), path) == 0 {
			return true
		}
	}
	if len(that.Include) > 0 && "" != strings.Join(that.Include, "") {
		return false
	}
	return true
}

// ParseTypes locates the *ast.TypeSpec for type id in the import path.
func (that *meshable) ParseTypes(ctx context.Context, dir string, pkg *ast.Package, sets *token.FileSet) (map[string]*Kind, error) {
	types := map[string]*Kind{}
	for file := range pkg.Files {
		fast, err := parser.ParseFile(sets, file, nil, parser.ParseComments)
		if nil != err {
			log.Error(ctx, err.Error())
			continue
		}
		for _, dcl := range fast.Decls {
			decl, ok := dcl.(*ast.GenDecl)
			if !ok || len(decl.Specs) < 1 || !(decl.Tok == token.TYPE || decl.Tok == token.VAR) {
				continue
			}
			if vt, ok := decl.Specs[0].(*ast.TypeSpec); ok && decl.Tok == token.TYPE {
				if _, ok := vt.Type.(*ast.InterfaceType); ok {
					for _, sp := range decl.Specs {
						spec := sp.(*ast.TypeSpec)
						methods, err := that.ParseMethods(ctx, sets, pkg, spec)
						if nil != err {
							log.Error(ctx, err.Error())
							continue
						}
						annotations := that.CommentsVars(ctx, vt.Name.Name, that.FlattenCommentGroup(decl.Doc))
						types[fmt.Sprintf("%s.%s", pkg.Name, spec.Name.Name)] = &Kind{
							Package:     pkg,
							Suckage:     &ast.Package{},
							Super:       spec,
							Type:        spec,
							File:        fast,
							SuFile:      fast,
							Files:       sets,
							Comment:     decl.Doc,
							Methods:     methods,
							Imports:     that.CollectImports(dir, pkg, fast),
							Attachments: map[string]string{"Name": that.AssignName(spec.Name.Name, annotations)},
							Annotations: annotations,
						}
					}
				}
				continue
			}
			if vs, ok := decl.Specs[0].(*ast.ValueSpec); ok {
				if nil == vs.Type {
					continue
				}
				reference := that.ParseDeclaration(ctx, sets, pkg, fast, vs.Type)
				super, ok := reference.Decl.(*ast.TypeSpec)
				if !ok {
					continue
				}
				service := that.ParseDeclaration(ctx, sets, pkg, fast, that.ParseVarValue(vs))
				ts, ok := service.Decl.(*ast.TypeSpec)
				if !ok {
					continue
				}
				methods, err := that.ParseMethods(ctx, sets, pkg, super)
				if nil != err {
					log.Error(ctx, err.Error())
					continue
				}
				for _, method := range methods {
					if nil != method.Annotations["MPI"] && "" != method.Annotations["MPI"]["Name"] {
						method.Annotations["MPI"]["Name"] = fmt.Sprintf("strings.ReplaceAll(%s, \"${mesh.name}\", tool.Name.Get())", method.Annotations["MPI"]["Name"])
					}
				}
				annotations := that.CommentsVars(ctx, service.Name, that.FlattenCommentGroup(service.Comments))
				if len(annotations) < 1 {
					continue
				}
				kind := &Kind{
					Package:     service.Package,
					Suckage:     reference.Package,
					Super:       super,
					Type:        ts,
					File:        fast,
					SuFile:      reference.File,
					Files:       sets,
					Comment:     service.Comments,
					Methods:     methods,
					Imports:     that.CollectImports(dir, service.Package, fast, reference.File),
					Attachments: map[string]string{"Name": that.AssignName(ts.Name.Name, annotations), "MPX": that.AssignMPXName(dir, service.Package, fast, reference.File)},
					Annotations: annotations,
				}
				// tn := tool.Ternary(that.HasAnnotation(kind, "Binding"), ts.Name.Name, super.Name.Name)
				types[fmt.Sprintf("%s.%s", pkg.Name, ts.Name.Name)] = kind
			}
		}
	}
	return types, nil
}

// CollectImportMap will collect and unique imports.
func (that *meshable) CollectImportMap(dir string, pkg *ast.Package, files ...*ast.File) map[string]string {
	homer := fmt.Sprintf("%s%s", that.Module, strings.ReplaceAll(dir, that.Src, ""))
	importers := map[string]string{filepath.Base(that.Macro): fmt.Sprintf("\"%s\"", that.Macro)}
	for _, file := range files {
		for _, importer := range file.Imports {
			if "" != importers[importer.Path.Value] && (nil == importer.Name || "" == importer.Name.Name) {
				continue
			}
			im := func() string {
				if nil != importer.Name && "" != importer.Name.Name {
					return importer.Name.Name
				} else {
					return filepath.Base(strings.ReplaceAll(importer.Path.Value, "\"", ""))
				}
			}()
			importers[im] = importer.Path.Value
		}
	}
	if that.Src != that.Dst {
		if "" != importers[filepath.Base(homer)] {
			importers[fmt.Sprintf("%smps", filepath.Base(homer))] = fmt.Sprintf("\"%s\"", homer)
		} else {
			importers[filepath.Base(homer)] = fmt.Sprintf("\"%s\"", homer)
		}
	}
	if that.MPX != fmt.Sprintf("%s%s", that.Module, strings.ReplaceAll(that.Dst, that.Src, "")) {
		importers[fmt.Sprintf("%smpx", filepath.Base(that.MPX))] = fmt.Sprintf("\"%s\"", that.MPX)
	}
	return importers
}

// CollectImports will collect and unique imports.
func (that *meshable) CollectImports(dir string, pkg *ast.Package, files ...*ast.File) []string {
	importers := that.CollectImportMap(dir, pkg, files...)
	var paths []string
	for name, path := range importers {
		if name == filepath.Base(strings.ReplaceAll(path, "\"", "")) {
			paths = append(paths, path)
		} else {
			paths = append(paths, fmt.Sprintf("%s %s", name, path))
		}
	}
	return paths
}

// AssignMPXName will assign the mesh proxy interface name.
func (that *meshable) AssignMPXName(dir string, pkg *ast.Package, files ...*ast.File) string {
	if that.MPX == fmt.Sprintf("%s%s", that.Module, strings.ReplaceAll(that.Dst, that.Src, "")) {
		return ""
	}
	importers := that.CollectImportMap(dir, pkg, files...)
	for name, importer := range importers {
		if that.MPX == strings.ReplaceAll(importer, "\"", "") {
			return fmt.Sprintf("%s.", name)
		}
	}
	return fmt.Sprintf("%s.", filepath.Base(that.MPX))
}

// ParseVarValue parse the variable value type.
func (that *meshable) ParseVarValue(vs *ast.ValueSpec) ast.Expr {
	if len(vs.Values) > 0 {
		if vc, ok := vs.Values[0].(*ast.CallExpr); ok {
			return vc.Args[0]
		}
	}
	return nil
}

// ParseDeclaration parse the variable type.
func (that *meshable) ParseDeclaration(ctx context.Context, sets *token.FileSet, pkg *ast.Package, file *ast.File, vs ast.Expr) *Declaration {
	switch vt := vs.(type) {
	case *ast.SelectorExpr:
		for _, importer := range file.Imports {
			if dc := that.ImportType(ctx, importer.Path.Value, vt.Sel.Name, sets); nil != dc {
				return dc
			}
		}
	case *ast.Ident:
		for _, source := range pkg.Files {
			for _, dcl := range source.Decls {
				if decl, ok := dcl.(*ast.GenDecl); ok && decl.Tok == token.TYPE && len(decl.Specs) > 0 {
					if dt, ok := decl.Specs[0].(*ast.TypeSpec); ok && dt.Name.Name == vt.Name {
						return &Declaration{
							Name:     vt.Name,
							Package:  pkg,
							File:     file,
							Decl:     dt,
							Comments: decl.Doc,
						}
					}
				}
			}
		}
	}
	return new(Declaration)
}

// ImportType filter the type in source file.
func (that *meshable) ImportType(ctx context.Context, path string, name string, sets *token.FileSet) *Declaration {
	if nil == that.declarations[path] {
		that.declarations[path] = map[string]*Declaration{}
	}
	if nil != that.declarations[path][name] {
		return that.declarations[path][name]
	}
	pkg, err := build.Import(strings.ReplaceAll(path, "\"", ""), "", build.ImportComment)
	if nil != err {
		log.Error(ctx, err.Error())
		return new(Declaration)
	}
	sources, err := parser.ParseDir(sets, pkg.Dir, func(info os.FileInfo) bool { return !strings.Contains(info.Name(), "_test.go") }, parser.ParseComments)
	if nil != err {
		log.Error(ctx, err.Error())
		return new(Declaration)
	}
	for _, source := range sources {
		for _, fast := range source.Files {
			for _, dcl := range fast.Decls {
				if decl, ok := dcl.(*ast.GenDecl); ok && decl.Tok == token.TYPE && len(decl.Specs) > 0 {
					if vt, ok := decl.Specs[0].(*ast.TypeSpec); ok {
						that.declarations[path][vt.Name.Name] = &Declaration{
							Name:     vt.Name.Name,
							Package:  source,
							File:     fast,
							Decl:     vt,
							Comments: decl.Doc,
						}
					}
				}
			}
		}
	}
	return that.declarations[path][name]
}

// FilterTypeSpec filter the type in source file.
func (that *meshable) FilterTypeSpec(name string, file *ast.File) *ast.TypeSpec {
	for _, dcl := range file.Decls {
		decl, ok := dcl.(*ast.GenDecl)
		if !ok || len(decl.Specs) < 1 || !(decl.Tok == token.TYPE || decl.Tok == token.VAR) {
			continue
		}
		for _, spec := range decl.Specs {
			if ts, ok := spec.(*ast.TypeSpec); ok && ts.Name.Name == name {
				return ts
			}
		}
	}
	return nil
}

// ParseMethods returns the set of methods required to implement iface.
// It is called ParseMethods rather than methods because the
// function descriptions are functions; there is no receiver.
// Special case for the built-in error interface.
func (that *meshable) ParseMethods(ctx context.Context, files *token.FileSet, pack *ast.Package, kind *ast.TypeSpec) ([]*Method, error) {
	var methods []*Method
	name := kind.Name.Name
	if name == "error" {
		// The error interface is built-in.
		methods = append(methods, &Method{
			Name:        "Error",
			Res:         []*Param{{Type: "string"}},
			Attachments: map[string]string{"Name": "Error"},
		})
		return methods, nil
	}
	astType, ok := kind.Type.(*ast.InterfaceType)
	if !ok || nil == astType.Methods {
		return methods, nil
	}
	for _, method := range astType.Methods.List {
		if len(method.Names) == 0 {
			// Embedded interface: recurse
			// embedded, err := that.ParseMethods(kind.FullType(method.Type), "srcDir")
			// if nil != err {
			//     return nil, cause.Error(err)
			// }
			continue
		}
		methods = append(methods, that.ParseMethod(ctx, files, pack, method))
	}
	return methods, nil
}

// gofmt pretty-prints e.
func (that *meshable) gofmt(files *token.FileSet, expr ast.Expr) string {
	var buf bytes.Buffer
	log.Catch(printer.Fprint(&buf, files, expr))
	return buf.String()
}

// FullType returns the fully qualified type of e.
// Examples, assuming package net/http:
//
//	FullType(int) => "int"
//	FullType(Handler) => "http.Handler"
//	FullType(io.Reader) => "io.Reader"
//	FullType(*Request) => "*http.Request"
func (that *meshable) FullType(files *token.FileSet, pack *ast.Package, expr ast.Expr) string {
	ast.Inspect(expr, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Ident:
			// Using ParseTypes instead of IsExported here would be
			// more accurate, but it'd be crazy expensive, and if
			// the type isn't exported, there's no point trying
			// to implement it anyway.
			if n.IsExported() {
				n.Name = pack.Name + "." + n.Name
			}
		case *ast.SelectorExpr:
			return false
		}
		return true
	})
	return that.gofmt(files, expr)
}

func (that *meshable) ParseParam(files *token.FileSet, pack *ast.Package, field *ast.Field) []*Param {
	var params []*Param
	typ := that.FullType(files, pack, field.Type)
	for _, name := range field.Names {
		params = append(params, &Param{Name: name.Name, Type: typ})
	}
	// Handle anonymous ParseParam
	if len(params) == 0 {
		params = []*Param{{Type: typ}}
	}
	return params
}

func (that *meshable) ParseMethod(ctx context.Context, files *token.FileSet, pack *ast.Package, field *ast.Field) *Method {
	fn := &Method{Name: field.Names[0].Name, Attachments: map[string]string{}}
	typ := field.Type.(*ast.FuncType)
	if typ.Params != nil {
		for _, field := range typ.Params.List {
			for _, param := range that.ParseParam(files, pack, field) {
				// only for method parameters:
				// assign a blank identifier "_" to an anonymous parameter
				if param.Name == "" {
					param.Name = "_"
				}
				fn.Params = append(fn.Params, param)
			}
		}
	}
	if typ.Results != nil {
		for _, f := range typ.Results.List {
			fn.Res = append(fn.Res, that.ParseParam(files, pack, f)...)
		}
	}
	fn.Comments = that.FlattenCommentGroup(field.Doc)
	fn.Annotations = that.CommentsVars(ctx, fn.Name, fn.Comments)
	fn.Attachments["Name"] = that.AssignName(fn.Name, fn.Annotations)
	return fn
}

// CommentsVars reports whether commentGroups precedes a field.
func (that *meshable) CommentsVars(ctx context.Context, name string, comment string) map[string]map[string]string {
	annotations := map[string]map[string]string{}
	explains := strings.Split(comment, "\n")
	for _, explain := range explains {
		regex, err := regexp.Compile("// @([a-zA-Z0-9]+(\\(.*\\))?)")
		if nil != err {
			log.Error(ctx, err.Error())
			return annotations
		}
		if strings.Index(explain, ")") > -1 {
			explain = explain[0 : strings.Index(explain, ")")+1]
		}
		expressions := regex.FindAllString(explain, 31)
		for _, expr := range expressions {
			index := strings.Index(expr, "(")
			if index < 0 {
				annotations[strings.ReplaceAll(expr, "// @", "")] = map[string]string{"Name": fmt.Sprintf("\"%s\"", name)}
				continue
			}
			mac := strings.TrimSpace(strings.ReplaceAll(expr[0:index], "// @", ""))
			metadata := expr[index+1 : len(expr)-1]
			if strings.Index(expr, "=") < 0 {
				annotations[mac] = map[string]string{"Name": strings.TrimSpace(metadata)}
				continue
			}
			annotations[mac] = map[string]string{}
			pairs := strings.Split(metadata, ",")
			for _, pair := range pairs {
				kv := strings.Split(strings.TrimSpace(pair), "=")
				if len(kv) > 1 {
					annotations[mac][strings.TrimSpace(tool.FistUpper(kv[0]))] = strings.TrimSpace(kv[1])
				}
			}
		}
	}
	return annotations
}

// CommentsType reports the comments of type key
func (that *meshable) CommentsType(comments ast.CommentMap, name string) []*ast.CommentGroup {
	var cms []*ast.CommentGroup
	for key, comment := range comments {
		if decl, ok := key.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
			if len(decl.Specs) < 1 {
				continue
			}
			if ts, ok := decl.Specs[0].(*ast.TypeSpec); ok && ts.Name.Name == name {
				cms = append(cms, comment...)
			}
		}
	}
	return cms
}

// FlattenCommentGroup flattens the comment map to a string.
// This function must be used at the point when m is expected to have a single
// element.
func (that *meshable) FlattenCommentGroup(comments ...*ast.CommentGroup) string {
	var explain strings.Builder
	for _, commentGroup := range comments {
		if nil == commentGroup {
			continue
		}
		for index, comment := range commentGroup.List {
			explain.WriteString(comment.Text)
			// add an end-of-line character if this is '//'-style comment
			if index < len(commentGroup.List)-1 && comment.Text[1] == '/' {
				explain.WriteString("\n")
			}
		}
	}

	// for '/*'-style comments, make sure to append EOL character to the comment
	// block
	return explain.String()
}

func (that *meshable) HasAnnotation(kind *Kind, annotations ...string) bool {
	for _, annotation := range annotations {
		if len(kind.Annotations[annotation]) > 0 {
			return true
		}
	}
	for _, method := range kind.Methods {
		for _, annotation := range annotations {
			if len(method.Annotations[annotation]) > 0 {
				return true
			}
		}
	}
	return false
}

func (that *meshable) AssignName(dft string, annotations map[string]map[string]string) string {
	if nil != annotations["SPI"] {
		return tool.Anyone(annotations["SPI"]["Name"], "macro.MeshMPI")
	}
	if nil != annotations["MPI"] {
		return tool.Anyone(annotations["MPI"]["Name"], "macro.MeshMPI")
	}
	if nil != annotations["MPS"] {
		return annotations["MPS"]["Name"]
	}
	if nil != annotations["Binding"] {
		return fmt.Sprintf("\"%-%s\"", annotations["Binding"]["Topic"], annotations["Binding"]["Code"])
	}
	return fmt.Sprintf("\"%s\"", dft)
}

// StaticDynamic prints nicely formatted method stubs
// for fns using receiver expression receiver.
// If receiver is not a valid receiver expression,
// genStubs will panic.
// genStubs won't generate stubs for
// already implemented methods of receiver.
func (that *meshable) StaticDynamic(ctx context.Context, types map[string]map[string]*Kind, dest string) error {
	for _, kinds := range types {
		if _, err := os.Stat(dest); os.IsNotExist(err) {
			log.Panic(os.MkdirAll(dest, os.ModePerm))
		}
		latch := sync.WaitGroup{}
		for _, kind := range kinds {
			if len(kind.Methods) < 1 || !that.HasAnnotation(kind, "MPI", "MPS", "Binding") {
				continue
			}
			latch.Add(1)
			that.GoStaticDynamic(ctx, kind, dest, func() { latch.Done() })
		}
		latch.Wait()
	}
	return nil
}

// GoStaticDynamic execute in go routine
func (that *meshable) GoStaticDynamic(ctx context.Context, kind *Kind, dest string, finalizer func()) {
	go func() {
		defer finalizer()
		marker := func(kind *Kind) *template.Template {
			if that.HasAnnotation(kind, "MPS", "Binding") {
				return mps
			}
			return mpi
		}
		log.Info(ctx, kind.Type.Name.Name)
		kind.Attachments["comma"] = "`"
		kind.Attachments["placeholder"] = "%s"
		if that.Src == that.Dst {
			kind.Attachments["pn"] = filepath.Base(that.Module)
			kind.Attachments["sn"] = kind.Type.Name.String()
			kind.Attachments["in"] = kind.Super.Name.Name
			kind.Attachments["an"] = kind.Type.Name.Name
		} else {
			kind.Attachments["pn"] = "proxy"
			kind.Attachments["sn"] = fmt.Sprintf("%s.%s", kind.Package.Name, kind.Type.Name)
			kind.Attachments["in"] = fmt.Sprintf("%s.%s", kind.Suckage.Name, kind.Super.Name.Name)
			kind.Attachments["an"] = fmt.Sprintf("%s%s.%s", kind.Package.Name, tool.Ternary(kind.Package.Name == kind.Suckage.Name, "mps", ""), kind.Type.Name)
		}
		var buf bytes.Buffer
		if err := marker(kind).Execute(&buf, kind); nil != err {
			log.Info(ctx, buf.String())
			log.Error(ctx, err.Error())
			return
		}
		puff, err := imports.Process("", buf.Bytes(), &imports.Options{Comments: true, TabIndent: true, TabWidth: 8, FormatOnly: false})
		if nil != err {
			log.Info(ctx, buf.String())
			log.Error(ctx, err.Error())
			return
		}
		pretty, err := format.Source(puff)
		if nil != err {
			log.Info(ctx, string(puff))
			log.Error(ctx, err.Error())
			return
		}
		prefix := tool.Ternary(that.HasAnnotation(kind, "MPS", "Binding"), "mps", "mpi")
		name := tool.Ternary(that.HasAnnotation(kind, "Binding"), strings.ToLower(kind.Type.Name.Name), strings.ToLower(kind.Super.Name.Name))
		fd, err := os.Create(filepath.Join(dest, fmt.Sprintf("%s_%s.go", prefix, name)))
		if nil != err {
			log.Info(ctx, string(pretty))
			log.Error(ctx, err.Error())
			return
		}
		defer func() { log.Catch(fd.Close()) }()
		if _, err = fd.Write(pretty); nil != err {
			log.Info(ctx, string(pretty))
			log.Error(ctx, err.Error())
			return
		}
	}()
}

var fns = map[string]interface{}{
	"Neq": func(x string, y string) bool {
		return x != y
	},
	"Eq": func(x int, y int) bool {
		return x > y
	},
	"NoReturn": func(params []*Param) bool {
		if len(params) < 1 {
			return true
		}
		if len(params) == 1 && params[0].Type == "error" {
			return true
		}
		return false
	},
	"ReturnType": func(params []*Param) string {
		for _, param := range params {
			if params[0].Type == "error" {
				continue
			}
			return param.Type
		}
		return "interface{}"
	},
	"HasParameters": func(method *Method) bool {
		return len(method.Params) > 1 && method.Params[0].Type == "context.Context"
	},
	"FistUpper": func(v string) string {
		return tool.FistUpper(v)
	},
	"FistLower": func(v string) string {
		return tool.FistLower(v)
	},
	"Minus": func(x int, y int) int {
		return x - y
	},
	"Multiply": func(x int, y int, z int) int {
		return x*y + z
	},
	"LowerUnder": func(v string) string {
		x, _ := regexp.Compile("[A-Z]")
		return strings.ToLower(x.ReplaceAllString(tool.FistLower(v), "_$0"))
	},
	"Has": func(annotations map[string]map[string]string, key string) bool {
		return nil != annotations && len(annotations[key]) > 0
	},
	"String": func(annotations map[string]map[string]string, key string) string {
		switch key {
		case "Binding":
			topic := strings.ReplaceAll(annotations[key]["Topic"], "\"", "")
			code := strings.ReplaceAll(annotations[key]["Code"], "\"", "")
			return fmt.Sprintf("\"%s.%s\"", topic, tool.Ternary(code == "", "*", code))
		default:
			return annotations[key]["Name"]
		}
	},
}

var mpi = template.Must(template.New("mpi.gohtml").Funcs(fns).ParseFS(home, "mpi.gohtml"))
var mps = template.Must(template.New("mps.gohtml").Funcs(fns).ParseFS(home, "mps.gohtml"))
