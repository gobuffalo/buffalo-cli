package git

import (
	"context"
	"os/exec"
)

type commandRunner struct {
	root   string
	cmd    *exec.Cmd
	stdout string
	err    error
}

func (v *commandRunner) PluginName() string {
	return "commandRunner"
}

var _ CommandRunner = &commandRunner{}

func (v *commandRunner) RunGitCommand(ctx context.Context, root string, cmd *exec.Cmd) error {
	v.cmd = cmd
	v.root = root
	if len(v.stdout) > 0 {
		v.cmd.Stdout.Write([]byte(v.stdout))
	}
	return v.err
}
