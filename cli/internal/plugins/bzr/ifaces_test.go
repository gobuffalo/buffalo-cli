package bzr

import (
	"context"
	"os/exec"
)

var _ CommandRunner = &commandRunner{}

type commandRunner struct {
	cmd    *exec.Cmd
	stdout string
	err    error
}

func (v *commandRunner) PluginName() string {
	return "commandRunner"
}

func (v *commandRunner) RunBzrCommand(ctx context.Context, root string, cmd *exec.Cmd) error {
	v.cmd = cmd
	if len(v.stdout) > 0 {
		v.cmd.Stdout.Write([]byte(v.stdout))
	}
	return v.err
}
