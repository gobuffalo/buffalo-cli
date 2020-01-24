package build

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

func (b *bladeRunner) BuildAssets(ctx context.Context, root string, cmd *exec.Cmd) error {
	b.cmd = cmd
	return b.err
}

var _ Tooler = &tooler{}

type tooler struct {
	root string
	tool string
	err  error
}

func (tooler) Name() string {
	return "tooler"
}

func (tool *tooler) AssetTool(ctx context.Context, root string) (string, error) {
	tool.root = root
	return tool.tool, tool.err
}
