package git

import (
	"context"
	"os/exec"
)

type versionRunner struct {
	cmd     *exec.Cmd
	version string
	err     error
}

func (v *versionRunner) Name() string {
	return "versionRunner"
}

var _ VersionRunner = &versionRunner{}

func (v *versionRunner) RunGitVersion(ctx context.Context, cmd *exec.Cmd) (string, error) {
	v.cmd = cmd
	return v.version, v.err
}
