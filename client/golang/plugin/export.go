/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package plugin

import (
	"fmt"
	"reflect"
)

type factories map[string]Factory

var instance = factories{}
var informers = map[string]Informer{}

// Provide os executor provider with name in init function
func Provide(spi Factory) {
	name := spi.Ptt().Name
	if nil != instance[string(name)] {
		panic(fmt.Errorf("Duplicate spi provider %s of %s has be provide in %s, please check and disable anyone ",
			string(name), reflect.ValueOf(spi).Type().String(), reflect.ValueOf(instance[string(name)]).Type().String()))
	}
	instance[string(name)] = spi
}

// Conv register Informer
func Conv(name string, informer Informer) {
	informers[name] = informer
}

// Load the plugin factory
func Load(name Name) Factory {
	return instance[string(name)]
}

// LoadAll load all plugin factories
func LoadAll() map[string]Factory {
	return instance
}

// LoadC is load plugin as a Container to start
func LoadC(name Name) Container {
	if nil == Load(name) {
		panic(fmt.Errorf("Cant found plugin named %s ", string(name)))
	}
	return &sharedContainer{name: name}
}

// ForkRuntime will fork a plugin runtime of plugin if absent.
func ForkRuntime(name Name, flags ...string) Runtime {
	return &sharedContainer{name: name, flags: flags}
}

// Inspect will scan all plugin to parse the bootstrap arguments.
func Inspect() ([]byte, error) {
	return flags.inspectAsFormatted()
}

// Parse the conf
func Parse(conf string, ptr interface{}) error {
	return flags.readConfAsStruct(conf, ptr)
}
