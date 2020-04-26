package newapp

import (
	"context"
	"os/exec"
)

var _ AfterNewapper = afternewapper(nil)
var _ Newapper = newapper(nil)
var _ NewCommandRunner = cmdRunner(nil)

type newapper func(ctx context.Context, root string, name string, args []string) error

func (fn newapper) PluginName() string {
	return "newapper"
}

func (fn newapper) Newapp(ctx context.Context, root string, name string, args []string) error {
	return fn(ctx, root, name, args)
}

type afternewapper func(ctx context.Context, root string, name string, args []string, err error) error

func (fn afternewapper) PluginName() string {
	return "after-newapper"
}

func (fn afternewapper) AfterNewapp(ctx context.Context, root string, name string, args []string, err error) error {
	return fn(ctx, root, name, args, err)
}

type cmdRunner func(ctx context.Context, root string, cmd *exec.Cmd) error

func (fn cmdRunner) PluginName() string {
	return "cmd-runner"
}

func (fn cmdRunner) RunNewCommand(ctx context.Context, root string, cmd *exec.Cmd) error {
	return fn(ctx, root, cmd)
}
