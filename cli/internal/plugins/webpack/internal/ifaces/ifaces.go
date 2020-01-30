package ifaces

import (
	"context"
)

type Tooler interface {
	// AssetTool returns the name of the asset tool to use.
	// npm, yarnpkg, etc...
	AssetTool(ctx context.Context, root string) (string, error)
}

type ToolerFn func(ctx context.Context, root string) (string, error)

func (fn ToolerFn) AssetTool(ctx context.Context, root string) (string, error) {
	return fn(ctx, root)
}

func (ToolerFn) PluginName() string {
	return "ToolerFn"
}

type Scripter interface {
	AssetScript(ctx context.Context, root string, name string) (string, error)
}
