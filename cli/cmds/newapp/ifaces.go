package newapp

import (
	"context"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
)

type Stdouter = plugio.Outer
type Stdiner = plugio.Inner
type Stderrer = plugio.Errer

type Newapper interface {
	plugins.Plugin
	Newapp(ctx context.Context, root string, args []string) error
}

type AfterNewapper interface {
	plugins.Plugin
	AfterNewapp(ctx context.Context, root string, args []string, err error) error
}
