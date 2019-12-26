package assets

import (
	"context"
	"os/exec"
)

type AssetBuilderContext interface {
	context.Context
	BuildAssets(cmd *exec.Cmd) error
}

type buildContext struct {
	context.Context
	fn func(cmd *exec.Cmd) error
}

func (c *buildContext) BuildAssets(cmd *exec.Cmd) error {
	if c.fn == nil {
		return nil
	}
	return c.fn(cmd)
}

func WithBuilderContext(ctx context.Context, fn func(cmd *exec.Cmd) error) AssetBuilderContext {
	return &buildContext{
		Context: ctx,
		fn:      fn,
	}
}
