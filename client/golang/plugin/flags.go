/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package plugin

import (
	"fmt"
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/schema"
	"github.com/be-io/mesh/client/golang/tool"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var flags = &structureFlags{Manifest: map[string]map[string]*structureFlags{}, Regex: regexp.MustCompile("\\$\\{(.*?)\\}")}

type structureFlags struct {
	Once     sync.Once                             `json:"-" yaml:"-"`
	Regex    *regexp.Regexp                        `json:"-" yaml:"-"`
	Manifest map[string]map[string]*structureFlags `json:"-" yaml:"-"`
	Key      string                                `json:"key" yaml:"key"`
	Dft      string                                `json:"dft" yaml:"dft"`
	Usage    string                                `json:"usage" yaml:"usage"`
	Expr     []string                              `json:"expr" yaml:"expr"`
	Kind     reflect.Type                          `json:"kind" yaml:"kind"`
}

func (that *structureFlags) inspectType(structure interface{}) map[string]*structureFlags {
	sfs := map[string]*structureFlags{}
	proto := func() reflect.Type {
		if v := reflect.TypeOf(structure); v.Kind() == reflect.Ptr {
			return v.Elem()
		}
		return reflect.TypeOf(structure)
	}()
	for index := 0; index < proto.NumField(); index++ {
		name := proto.Field(index).Tag.Get("json")
		if "" == strings.TrimSpace(name) {
			name = proto.Field(index).Name
		}
		if "-" == name {
			continue
		}
		flag := &structureFlags{
			Key:   proto.Field(index).Tag.Get("key"),
			Dft:   proto.Field(index).Tag.Get("dft"),
			Usage: proto.Field(index).Tag.Get("usage"),
			Expr:  []string{proto.Field(index).Tag.Get("expr")},
			Kind:  proto.Field(index).Type,
		}
		matches := that.Regex.FindAllStringSubmatch(flag.Dft, 10)
		for _, pair := range matches {
			flag.Expr = append(flag.Expr, pair[1])
		}
		sfs[that.toHump(name)] = flag
	}
	return sfs
}

func (that *structureFlags) inspect() map[string]map[string]*structureFlags {
	that.Once.Do(func() {
		for _, plugin := range LoadAll() {
			attr := plugin.Ptt()
			that.Manifest[that.toHump(string(attr.Name))] = that.inspectType(attr.Flags)
		}
	})
	return that.Manifest
}

func (that *structureFlags) variables(args ...string) map[string]string {
	home, _ := os.UserHomeDir()
	pwd, _ := filepath.Abs(".")
	vars := map[string]string{
		"CPU":       strconv.Itoa(runtime.NumCPU()),
		"HOME":      home,
		"TEMP":      os.TempDir(),
		"PWD":       pwd,
		"MESH_HOME": tool.Anyone(os.Getenv("MESH_HOME"), home),
		"IPH":       tool.IPHex.Get(),
	}
	for _, env := range os.Environ() {
		if strings.Index(env, "=") > -1 {
			index := strings.Index(env, "=")
			vars[env[0:index]] = env[index+1:]
		}
	}
	for _, arg := range os.Args {
		if strings.Index(arg, "--") == 0 && strings.Index(arg, "=") > -1 {
			index := strings.Index(arg, "=")
			vars[arg[2:index]] = arg[index+1:]
		}
	}
	for _, arg := range args {
		if strings.Index(arg, "--") == 0 && strings.Index(arg, "=") > -1 {
			index := strings.Index(arg, "=")
			vars[arg[2:index]] = arg[index+1:]
		}
	}
	return vars
}

func (that *structureFlags) readConfAsStruct(conf string, ptr interface{}) error {
	vars := that.variables()
	props, err := that.readAsMap(Format, conf)
	if nil != err {
		return err
	}
	if nil == props {
		props = map[string]interface{}{}
	}
	configuration := that.inspectType(ptr)
	for fn, flag := range configuration {
		hump := that.toHump(fn)
		underline := that.toUnderline(fn)
		upperHump := strings.ToUpper(hump)
		upperUnderline := strings.ToUpper(underline)
		lowerHump := strings.ToLower(hump)
		value := that.matchBestFlag(flag, vars, props, hump, underline, upperHump, upperUnderline, lowerHump)
		fns := that.compatibleUnderlineOrHump(fn)
		for _, name := range fns {
			props[name] = value
		}
	}
	bytes, err := codec.Jsonizer.Marshal(props)
	if nil != err {
		return err
	}
	return codec.Jsonizer.Unmarshal(bytes, ptr)
}

func (that *structureFlags) readAsStruct(pluginName string, ptr interface{}, args ...string) error {
	if nil == that.inspect()[pluginName] {
		return nil
	}
	vars := that.variables(args...)
	tree, err := that.readAsMap(Format, Config)
	if nil != err {
		return err
	}
	for pn, plugin := range that.inspect() {
		if pn != pluginName {
			continue
		}
		pns := that.compatibleUnderlineOrHump(pn)
		props := map[string]interface{}{}
		for _, name := range pns {
			if kv, ok := tree[name].(map[string]interface{}); ok && len(kv) > 0 {
				props = kv
			}
		}
		for _, name := range pns {
			tree[name] = props
		}
		for fn, flag := range plugin {
			hump := that.toHump(fn)
			underline := that.toUnderline(fn)
			upperHump := strings.ToUpper(hump)
			upperUnderline := strings.ToUpper(underline)
			lowerHump := strings.ToLower(hump)
			value := that.matchBestFlag(flag, vars, props, hump, underline, upperHump, upperUnderline, lowerHump)
			fns := that.compatibleUnderlineOrHump(fn)
			for _, name := range fns {
				props[name] = value
			}
		}
	}
	bytes, err := codec.Jsonizer.Marshal(tree[pluginName])
	if nil != err {
		return err
	}
	return codec.Jsonizer.Unmarshal(bytes, ptr)
}

func (that *structureFlags) inspectAsFormatted() ([]byte, error) {
	tree := that.inspect()
	switch Format {
	case "yaml":
		return yaml.Marshal(tree)
	case "json":
		return codec.Jsonizer.Marshal(tree)
	default:
		return nil, fmt.Errorf("Cant resolve the format type of %s ", Format)
	}
}

func (that *structureFlags) readAsMap(format string, config string) (map[string]interface{}, error) {
	input, err := that.readAsBytes(config)
	if nil != err {
		return nil, err
	}
	var ptr map[string]interface{}
	switch format {
	case "yaml":
		return ptr, yaml.Unmarshal(input, &ptr)
	case "json":
		return ptr, codec.Jsonizer.Unmarshal(input, &ptr)
	default:
		return nil, fmt.Errorf("Cant resolve the format type of %s ", format)
	}
}

func (that *structureFlags) readAsBytes(config string) ([]byte, error) {
	if strings.Index(config, "http:") == 0 || strings.Index(config, "https:") == 0 {
		response, err := http.Get(config)
		if nil != err {
			return nil, err
		}
		defer func() { _ = response.Body.Close() }()
		var input []byte
		_, err = response.Body.Read(input)
		return input, err
	}
	_, err := os.Stat(config)
	if os.IsExist(err) {
		return os.ReadFile(config)
	}
	if filepath.IsAbs(config) {
		return that.readConfCreateIfAbsent(config)
	}
	return []byte(config), nil
}

func (that *structureFlags) compatibleUnderlineOrHump(name string) []string {
	if strings.Index(name, "_") < 0 && strings.ToLower(name) == name {
		return []string{name}
	}
	if strings.Index(name, "_") < 0 {
		return []string{name, that.toUnderline(name), strings.ToLower(name), that.toDot(that.toUnderline(name))}
	}
	if strings.ToLower(name) == name {
		return []string{name, that.toHump(name), that.toDot(name)}
	}
	return []string{name, strings.ToLower(name), that.toDot(strings.ToLower(name))}
}

func (that *structureFlags) toHump(name string) string {
	return schema.Hump(name)
}

func (that *structureFlags) toUnderline(name string) string {
	return schema.Underline(name)
}

func (that *structureFlags) toDot(name string) string {
	return strings.ReplaceAll(name, "_", ".")
}

func (that *structureFlags) matchBestFlag(flag *structureFlags, vars map[string]string, props map[string]interface{}, keys ...string) interface{} {
	var value interface{}
	var cks []string
	for _, key := range keys {
		cks = append(cks, key, strings.ReplaceAll(key, ".", "_"), strings.ReplaceAll(key, "_", "."))
	}
	for _, key := range cks {
		if strings.Contains("home,HOME", key) {
			continue
		}
		if nil != props[key] && "" != props[key] {
			value = that.convertWithFlag(flag, fmt.Sprintf("%v", props[key]))
			break
		}
		if "" != vars[key] {
			value = that.convertWithFlag(flag, vars[key])
			break
		}
	}
	if nil == value || "" == value {
		value = that.convertWithFlag(flag, flag.Dft)
	}
	for _, expr := range flag.Expr {
		if "" != expr && "" != vars[expr] {
			value = strings.ReplaceAll(fmt.Sprintf("%v", value), fmt.Sprintf("${%s}", expr), vars[expr])
		}
	}
	return value
}

func (that *structureFlags) convertWithFlag(flag *structureFlags, value string) interface{} {
	switch flag.Kind.Name() {
	case "string":
		return value
	case "bool":
		return value == "true"
	case "Duration":
		v, _ := time.ParseDuration(value)
		return v
	case "int", "int32", "int64", "uint", "uint32", "uint64":
		v, _ := strconv.Atoi(value)
		return v
	default:
		for key, informer := range informers {
			if strings.ContainsAny(key, flag.Kind.Name()) {
				if v, err := informer(value); nil == err {
					return v
				}
			}
		}
		return value
	}
}
func (that *structureFlags) readConfCreateIfAbsent(conf string) ([]byte, error) {
	dir := filepath.Dir(conf)
	_, err := os.Stat(dir)
	if nil == err || os.IsNotExist(err) {
		return that.gracefulReadConf(conf)
	}
	if err := os.MkdirAll(dir, 0755); nil != err {
		return nil, err
	}
	return that.gracefulReadConf(conf)
}

func (that *structureFlags) gracefulReadConf(conf string) ([]byte, error) {
	fd, err := that.createConfIfAbsent(conf)
	if nil != err {
		return nil, err
	}
	defer func() {
		if err := fd.Close(); nil != err {
			fmt.Println(err.Error())
		}
	}()
	return ioutil.ReadAll(fd)
}

func (that *structureFlags) createConfIfAbsent(conf string) (*os.File, error) {
	_, err := os.Stat(conf)
	if nil != err && os.IsNotExist(err) {
		return os.Create(conf)
	}
	if nil != err {
		return nil, err
	}
	return os.Open(conf)
}
