# Mesh Golang Client

[![Build Status](https://travis-ci.org/ducesoft/babel.svg?branch=master)](https://travis-ci.org/ducesoft/babel)
[![Financial Contributors on Open Collective](https://opencollective.com/babel/all/badge.svg?label=financial+contributors)](https://opencollective.com/babel) [![codecov](https://codecov.io/gh/babel/babel/branch/master/graph/badge.svg)](https://codecov.io/gh/babel/babel)
![license](https://img.shields.io/github/license/ducesoft/babel.svg)

中文版 [README](README_CN.md)

## Introduction

Mesh Golang client develop kits base on Golang1.18.

```bash
go get github.com/opendatav/mesh/client/golang
```

## Features

Any

## Get Started

Static proxy mesh program interface.

```bash
go generate
```

Mesh client provider life cycle hooks.

```go
package main

import (
	"context"
	"fmt"
	"gitlab.trustbe.net/middleware/gaia-ark/internal/builtin"
	_ "gitlab.trustbe.net/middleware/gaia-ark/internal/kms"
	_ "gitlab.trustbe.net/middleware/gaia-ark/panel"
	_ "gitlab.trustbe.net/middleware/gaia-ark/proxy"
	"github.com/opendatav/mesh/client/golang/boost"
	_ "github.com/opendatav/mesh/client/golang/grpc"
	"github.com/opendatav/mesh/client/golang/log"
	_ "github.com/opendatav/mesh/client/golang/proxy"
	"github.com/opendatav/mesh/client/golang/prsim"
	_ "github.com/opendatav/mesh/client/golang/system"
	"os"
)

func main() {
	ctx := context.Background()
	runtime := &prsim.MeshRuntime{}
	mooter := new(boost.Mooter)
	if err := mooter.Start(ctx, runtime); nil != err {
		log.Fatal(ctx, err.Error())
		return
	}
	_, err := os.Stdout.WriteString(fmt.Sprintf(builtin.Banner, builtin.Version, builtin.CommitID))
	if nil != err {
		log.Catch(err)
	}
	log.Info(ctx, "Ark has been started. ")
	mooter.Wait(ctx, runtime)
}

```

