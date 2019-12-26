package assets

import (
	"context"
	"os/exec"
)

var _ AssetBuilder = &bladeRunner{}

type bladeRunner struct {
	cmd *exec.Cmd
	err error
}

func (bladeRunner) Name() string {
	return "blade"
}

func (b *bladeRunner) BuildAssets(ctx context.Context, cmd *exec.Cmd) error {
	b.cmd = cmd
	return b.err
}
