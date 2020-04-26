package newapp

import (
	"context"
	"flag"
	"os/exec"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/spf13/pflag"
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

type NewCommandRunner interface {
	plugins.Plugin
	// GoBuild will be called to build, and execute, the
	// presented context and args.
	// The first plugin to receive this call will be the
	// only to answer it.
	RunNewCommand(ctx context.Context, root string, cmd *exec.Cmd) error
}

type Flagger interface {
	plugins.Plugin
	NewappFlags() []*flag.Flag
}

type Pflagger interface {
	plugins.Plugin
	NewappFlags() []*pflag.Flag
}
