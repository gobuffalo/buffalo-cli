package builder

import (
	"context"
	"os/exec"
)

type AssetBuilder interface {
	BuildAssets(ctx context.Context, cmd *exec.Cmd) error
}

type Tooler interface {
	// AssetTool returns the name of the asset build tool to use.
	// npm, yarnpkg, etc...
	AssetTool(ctx context.Context, root string) (string, error)
}

type Scripter interface {
	AssetScripter(ctx context.Context, root string, name string) (string, error)
}

type toolerFn func(ctx context.Context, root string) (string, error)

func (fn toolerFn) AssetTool(ctx context.Context, root string) (string, error) {
	return fn(ctx, root)
}

func (toolerFn) Name() string {
	return "toolerFn"
}
