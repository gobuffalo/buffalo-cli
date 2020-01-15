package git

import (
	"context"
	"os/exec"
)

type commandRunner struct {
	cmd    *exec.Cmd
	stdout string
	err    error
}

func (v *commandRunner) Name() string {
	return "commandRunner"
}

var _ CommandRunner = &commandRunner{}

func (v *commandRunner) RunGitCommand(ctx context.Context, cmd *exec.Cmd) error {
	v.cmd = cmd
	if len(v.stdout) > 0 {
		v.cmd.Stdout.Write([]byte(v.stdout))
	}
	return v.err
}
