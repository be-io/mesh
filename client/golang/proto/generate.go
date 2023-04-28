/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/tool"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	protoc()
}

func protoc() {
	var src, dst, mod, java, python string
	flag.StringVar(&src, "s", "", "Source dir.")
	flag.StringVar(&dst, "d", "", "Generate Source dir.")
	flag.StringVar(&java, "j", "", "Generate Java language.")
	flag.StringVar(&python, "p", "", "Generate Python language.")
	flag.StringVar(&mod, "m", "github.com/be-io/mesh", "Generate Source module.")
	flag.Parse()
	ctx := macro.Context()
	pwd, err := os.Getwd()
	if nil != err {
		log.Error(ctx, err.Error())
		return
	}
	workspace := fmt.Sprintf("%s%s%s", pwd, string(filepath.Separator), src)
	protos := map[string]string{}
	var dirs []string
	dirs = append(dirs, workspace)
	for len(dirs) > 0 {
		dir := dirs[0]
		dirs = dirs[1:]
		entries, err := os.ReadDir(dir)
		if nil != err {
			log.Error(ctx, err.Error())
			return
		}
		for _, entry := range entries {
			path := strings.ReplaceAll(fmt.Sprintf("%s%s%s", dir, string(filepath.Separator), entry.Name()), "//", "/")
			if entry.IsDir() {
				dirs = append(dirs, path)
			} else if strings.Contains(entry.Name(), ".proto") {
				protos[entry.Name()] = strings.ReplaceAll(path, workspace, "")
			}
		}
	}
	command := bytes.Buffer{}
	command.WriteString("protoc ")
	command.WriteString(tool.Ternary("" != java, fmt.Sprintf("--java_out=%s ", java), ""))
	command.WriteString(tool.Ternary("" != python, fmt.Sprintf("--python_out=%s ", python), ""))
	command.WriteString("--go_out=. ")
	command.WriteString("--go-grpc_out=require_unimplemented_servers=false:. ")
	command.WriteString(fmt.Sprintf("-I=%s ", pwd))
	command.WriteString(strings.ReplaceAll(fmt.Sprintf("--go-grpc_opt=module=%s ", mod), "//", ""))
	command.WriteString(strings.ReplaceAll(fmt.Sprintf("--go_opt=module=%s ", mod), "//", ""))
	for _, proto := range protos {
		command.WriteString(fmt.Sprintf("%s%s", src, proto))
		command.WriteString(" ")
	}
	build := exec.Command("bash", "-c", command.String())
	build.Env = os.Environ()
	build.Dir = pwd
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	log.Info(ctx, build.String())
	if err = build.Run(); nil != err {
		log.Error(ctx, err.Error())
	}
}
