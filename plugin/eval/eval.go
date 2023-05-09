/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package eval

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/dsa"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	_ "gonum.org/v1/gonum"
	"hash/fnv"
	"reflect"
)

//go:generate go install github.com/traefik/yaegi/cmd/yaegi@v0.14.3
//go:generate yaegi extract gonum.org/v1/gonum/mat

func init() {
	var _ prsim.Evaluator = engine
	macro.Provide(prsim.IEvaluator, engine)
}

const (
	Name   = "yaegi"
	Prefix = "mesh.eval.expr"
)

var (
	Symbols     = map[string]map[string]reflect.Value{}
	engine      = &evaluator{scripts: dsa.NewStringMap[*program]()}
	interceptor = new(macro.Once[*interp.Interpreter]).With(func() *interp.Interpreter {
		ctx := mpc.Context()
		it := interp.New(interp.Options{})
		if err := it.Use(stdlib.Symbols); nil != err {
			log.Error(ctx, err.Error())
		}
		if err := it.Use(Symbols); nil != err {
			log.Error(ctx, err.Error())
		}
		return it
	})
	vs = `
	package %s
	import "context"
	func Invoke(ctx context.Context, args map[string]string, dft string) (string, error) {
		return "%s", nil
	}
	`
	es = `
	package %s
	import (
		"context"
		"fmt"
	)
	func Invoke(ctx context.Context, args map[string]string, dft string) (string, error) {
		return %s, nil
	}
	`
	ss = `
	package %s
	%s
	`
)

type program struct {
	signature string
	version   int
}

type evaluator struct {
	scripts dsa.Map[string, *program]
}

func (that *evaluator) Att() *macro.Att {
	return &macro.Att{Name: Name}
}

func (that *evaluator) Compile(ctx context.Context, script *types.Script) (string, error) {
	if "" == script.Code {
		script.Code = script.Name
	}
	if "" == script.Code || "" == script.Kind || "" == script.Name {
		return "", cause.Errorable(cause.Validate)
	}
	key := fmt.Sprintf("%s.%s", Prefix, script.Code)
	entity, err := aware.KV.Get(ctx, key)
	if nil != err {
		return "", cause.Error(err)
	}
	if nil == entity {
		entity = new(types.Entity)
	}
	if _, err = entity.Wrap(script); nil != err {
		return "", cause.Error(err)
	}
	if err = aware.KV.Put(ctx, key, entity); nil != err {
		return "", cause.Error(err)
	}
	_, err = that.scripts.Update(script.Code, func(k string, v *program) (*program, error) {
		if nil == v {
			return that.doCompile(ctx, script.Code, 0)
		}
		return that.doCompile(ctx, script.Code, v.version+1)
	})
	return script.Code, cause.Error(err)
}

func (that *evaluator) doCompile(ctx context.Context, code string, version int) (*program, error) {
	key := fmt.Sprintf("%s.%s", Prefix, code)
	entity, err := aware.KV.Get(ctx, key)
	if nil != err {
		return nil, cause.Error(err)
	}
	if nil == entity || !entity.Present() {
		log.Warn(ctx, "Script %s not present. ", code)
		return &program{signature: "", version: -1}, nil
	}
	var script types.Script
	if err = entity.TryReadObject(&script); nil != err {
		log.Warn(ctx, "Script %s cant be resolve, %s", code, err.Error())
		return &program{signature: "", version: -1}, nil
	}
	hash := fnv.New32a()
	if _, err = hash.Write([]byte(fmt.Sprintf("%s%d", code, version))); nil != err {
		return &program{signature: "", version: -1}, nil
	}
	pname := fmt.Sprintf("p%d", hash.Sum32())
	expr := func() string {
		if script.Kind == types.VALUE {
			return fmt.Sprintf(vs, pname, script.Expr)
		}
		if script.Kind == types.EXPRESSION {
			return fmt.Sprintf(es, pname, script.Expr)
		}
		return fmt.Sprintf(ss, pname, script.Expr)
	}()
	if _, err = interceptor.Get().Eval(expr); nil != err {
		return nil, cause.Errorc(cause.UnexpectedSyntax, err)
	}
	return &program{signature: fmt.Sprintf("%s.Invoke", pname), version: 0}, nil
}

func (that *evaluator) Exec(ctx context.Context, code string, args map[string]string, dft string) (string, error) {
	p, err := that.scripts.PutIfe(code, func(k string) (*program, error) { return that.doCompile(ctx, code, 0) })
	if nil != err {
		return "", cause.Error(err)
	}
	if p.version < 0 {
		log.Warn(ctx, "Script %s not present. ", code)
		return dft, nil
	}
	exec, err := interceptor.Get().Eval(p.signature)
	if nil != err {
		return "", cause.Error(err)
	}
	ex, ok := exec.Interface().(func(ctx context.Context, args map[string]string, dft string) (string, error))
	if !ok {
		return "", cause.Errorable(cause.UnexpectedSyntax)
	}
	v, err := ex(ctx, args, dft)
	if nil != err {
		return "", cause.Error(err)
	}
	return v, nil
}

func (that *evaluator) Dump(ctx context.Context, feature map[string]string) ([]*types.Script, error) {
	keys, err := aware.KV.Keys(ctx, Prefix)
	if nil != err {
		return nil, cause.Error(err)
	}
	var scripts []*types.Script
	for _, key := range keys {
		entity, err := aware.KV.Get(ctx, key)
		if nil != err {
			return nil, cause.Error(err)
		}
		if nil == entity.Buffer {
			continue
		}
		var script types.Script
		if err = entity.TryReadObject(&script); nil != err {
			return nil, cause.Error(err)
		}
		scripts = append(scripts, &script)
	}
	return scripts, nil
}

func (that *evaluator) Index(ctx context.Context, index *types.Paging) (*types.Page, error) {
	scripts, err := that.Dump(ctx, map[string]string{})
	if nil != err {
		return nil, cause.Error(err)
	}
	return new(types.Page).Reset(index, int64(len(scripts)), scripts), nil
}
