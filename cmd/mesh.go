/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/be-io/mesh/client/golang/boost"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/plugin"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/gofrs/flock"
	"github.com/spf13/cobra"
	args "github.com/spf13/pflag"
	"io"
	"os"
	"path/filepath"
	system "runtime"
	"strconv"
	"syscall"
)

func init() {
	plugin.Provide(new(meshCli))
}

type Command interface {
	Home(ctx context.Context) *cobra.Command
}

type CommandFn func(ctx context.Context) *cobra.Command

func (that CommandFn) Home(ctx context.Context) *cobra.Command {
	return that(ctx)
}

func Provide(command Command) {
	commands = append(commands, command.Home)
}

var commands []func(ctx context.Context) *cobra.Command

type meshCli struct {
	Server   bool   `json:"server"`
	Node     bool   `json:"node"`
	Panel    bool   `json:"panel"`
	Operator bool   `json:"operator"`
	Proxy    bool   `json:"proxy"`
	Home     string `json:"home" dft:"${MESH_HOME}/mesh" usage:"Mesh home dir."`
	LOCK     string `json:"lock" dft:"${MESH_HOME}/mesh/${IPH}/mesh.lock" usage:"Mesh main process lock file."`
	Stderr   string `json:"stderr" dft:"${MESH_HOME}/mesh/${IPH}/error.log" usage:"Mesh std error file"`
}

func (that *meshCli) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Name: "mesh", Flags: meshCli{}, Create: func() plugin.Plugin {
		return that
	}}
}

func (that *meshCli) Start(ctx context.Context, runtime plugin.Runtime) {
	log.Catch(runtime.Parse(that))
	boost.RedirectStderrFile(ctx, that.Stderr)
	if plugin.Debug {
		system.SetBlockProfileRate(20)
		system.SetMutexProfileFraction(20)
	}
	if macro.JsonLogFormat.Enable() {
		log.SetFormatter(new(log.Jsonify))
	}
	flag.CommandLine.SetOutput(io.Discard)
	args.CommandLine.SetOutput(io.Discard)
	root := cobra.Command{
		Use:           "Mesh",
		SilenceErrors: true,
		Version:       prsim.Version,
		Example:       "mesh COMMAND",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := os.Stdout.WriteString(cmd.UsageString())
			log.Catch(err)
		},
	}
	root.PersistentFlags().StringVarP(&plugin.Config, "config", "c", "{}", "Mesh configuration input, it can be url/path/body, default is {}.")
	root.PersistentFlags().StringVarP(&plugin.Format, "format", "f", "yaml", "Mesh configuration format, it can be json/yaml, default is yaml.")
	root.PersistentFlags().BoolVarP(&plugin.Debug, "debug", "d", false, "Set the run mode to debug with more information.")
	for _, command := range commands {
		root.AddCommand(command(ctx))
	}
	root.AddCommand(that.Boot(ctx))
	root.AddCommand(that.Shutdown(ctx))
	root.SetVersionTemplate(fmt.Sprintf("Mesh version %s, build %s %s/%s. \n", prsim.Version, prsim.CommitID, prsim.GOOS, prsim.GOARCH))
	if err := root.Execute(); nil != err {
		log.Error(ctx, err.Error())
	}
}

func (that *meshCli) Boot(ctx context.Context) *cobra.Command {
	start := &cobra.Command{
		Use:     "start",
		Version: prsim.Version,
		Short:   "Start mesh with the run mode(server/operator/node/panel).",
		Long:    "Start mesh with the run mode(server/operator/node/panel).",
		Run: func(cmd *cobra.Command, args []string) {
			plugin.OPERATOR.Enable(that.Operator)
			plugin.PANEL.Enable(that.Panel)
			plugin.NODE.Enable(that.Node)
			plugin.SERVER.Enable(that.Server)
			plugin.PROXY.Enable(that.Proxy)
			if plugin.SERVER.Match() {
				_ = os.Setenv("mesh.runtime", "127.0.0.1:8864")
			}
			macro.WithMode(macro.EableSPIFirst)
			//environ, err := macro.Load(prsim.INetwork).Get(macro.MeshSPI).(prsim.Network).GetEnviron(ctx)
			//if nil != err {
			//	log.Error(ctx, err.Error())
			//	return
			//}
			locker, pid, err := that.TryLockEvenIfAbsent(ctx)
			if nil != err {
				log.Error(ctx, err.Error())
				return
			}
			ok, err := locker.TryLock()
			if nil != err {
				log.Error(ctx, err.Error())
				return
			}
			defer func() {
				log.Catch(locker.Unlock())
			}()
			if ok {
				that.WritePID(ctx, pid, os.Getpid())
			} else {
				log.Warn(ctx, "Mesh cant lock %s, is mesh %d bootstrap already. ", that.LOCK, pid)
			}
			_, err = os.Stdout.WriteString(fmt.Sprintf(plugin.Banner, prsim.Version, prsim.CommitID))
			if nil != err {
				log.Catch(err)
			}
			// log.Info(ctx, "Node %s (%s[%s]) has been started. ", environ.NodeId, environ.InstId, environ.InstName)
			container := plugin.LoadC(that.DetermineRunMode())
			container.WaitAny(ctx, locker)
			container.Start(ctx)
			container.Wait(ctx)
		},
	}
	start.PersistentFlags().BoolVarP(&that.Server, "server", "", false, "Run mesh as mesh server.")
	start.PersistentFlags().BoolVarP(&that.Node, "node", "", false, "Run mesh as mesh node.")
	start.PersistentFlags().BoolVarP(&that.Panel, "panel", "", false, "Run mesh as mesh panel.")
	start.PersistentFlags().BoolVarP(&that.Operator, "operator", "", false, "Run mesh as k8s operator.")
	start.PersistentFlags().BoolVarP(&that.Proxy, "proxy", "", false, "Run mesh as network proxy.")
	return start
}

func (that *meshCli) Shutdown(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "stop",
		Version: prsim.Version,
		Short:   "Stop the mesh process if available.",
		Long:    "Stop the mesh process if available.",
		Run: func(cmd *cobra.Command, args []string) {
			locker, pid, err := that.TryLockEvenIfAbsent(ctx)
			if nil != err {
				log.Error(ctx, err.Error())
				return
			}
			ok, err := locker.TryLock()
			if nil != err {
				log.Error(ctx, err.Error())
				return
			}
			if ok {
				log.Warn(ctx, "None available mesh process to shutdown! ")
				return
			}
			if pid < 1 {
				log.Info(ctx, "Mesh has been stop now!")
				return
			}
			if err = syscall.Kill(pid, syscall.SIGKILL); nil != err {
				log.Error(ctx, err.Error())
			}
			log.Info(ctx, "Mesh has been stop now!")
		},
	}
}

func (that *meshCli) Stop(ctx context.Context, runtime plugin.Runtime) {
	locker, pid, err := that.TryLockEvenIfAbsent(ctx)
	if nil != err {
		log.Error(ctx, "Mesh process %d stop, %s", pid, err.Error())
		return
	}
	ok, err := locker.TryRLock()
	if nil != err {
		log.Error(ctx, err.Error())
		return
	}
	if ok {
		that.WritePID(ctx, pid, -1)
	}
	if err = locker.Unlock(); nil != err {
		log.Error(ctx, "Mesh process %d stop, %s", pid, err.Error())
	}
}

func (that *meshCli) TryLockEvenIfAbsent(ctx context.Context) (plugin.Locker, int, error) {
	if err := tool.MakeDir(filepath.Dir(that.LOCK)); nil != err {
		return nil, -1, cause.Error(err)
	}
	if err := tool.MakeFile(that.LOCK); nil != err {
		return nil, -1, cause.Error(err)
	}
	locker := flock.New(that.LOCK)
	text, err := os.ReadFile(that.LOCK)
	if nil != err || "" == string(text) {
		return locker, -1, err
	}
	pid, err := strconv.Atoi(string(text))
	return locker, pid, err
}

func (that *meshCli) WritePID(ctx context.Context, pid int, npid int) {
	if process, _ := os.FindProcess(pid); nil != process {
		if err := process.Signal(syscall.Signal(0)); nil == err {
			log.Warn(ctx, "Mesh has been already bootstrap with pid %d. ", pid)
		}
	}
	if err := os.WriteFile(that.LOCK, []byte(strconv.Itoa(npid)), 0644); nil != err {
		log.Error(ctx, err.Error())
	}
}

func (that *meshCli) DetermineRunMode() plugin.Name {
	if plugin.SERVER.Match() {
		return plugin.ServerCar
	}
	if plugin.NODE.Match() {
		return plugin.NodeCar
	}
	if plugin.PANEL.Match() {
		return plugin.PanelCar
	}
	if plugin.OPERATOR.Match() {
		return plugin.WorkloadCar
	}
	if plugin.PROXY.Match() {
		return plugin.ProxyCar
	}
	return plugin.ServerCar
}
