package newapp

import (
	"context"

	"github.com/gobuffalo/plugins"
)

var _ AfterNewapper = afternewapper(nil)
var _ Newapper = newapper(nil)
var _ plugins.Plugin = afternewapper(nil)
var _ plugins.Plugin = newapper(nil)

type newapper func(ctx context.Context, root string, args []string) error

func (fn newapper) PluginName() string {
	return "newapper"
}

func (fn newapper) Newapp(ctx context.Context, root string, args []string) error {
	return fn(ctx, root, args)
}

type afternewapper func(ctx context.Context, root string, args []string, err error) error

func (fn afternewapper) PluginName() string {
	return "afternewapper"
}

func (fn afternewapper) AfterNewapp(ctx context.Context, root string, args []string, err error) error {
	return fn(ctx, root, args, err)
}
