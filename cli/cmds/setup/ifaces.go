package setup

import (
	"context"
	"flag"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/spf13/pflag"
)

type BeforeSetuper interface {
	plugins.Plugin
	BeforeSetup(ctx context.Context, root string, args []string) error
}

type AfterSetuper interface {
	plugins.Plugin
	AfterSetup(ctx context.Context, root string, args []string, err error) error
}

type Setuper interface {
	plugins.Plugin
	Setup(ctx context.Context, root string, args []string) error
}

type Flagger interface {
	plugins.Plugin
	SetupFlags() []*flag.Flag
}

type Pflagger interface {
	plugins.Plugin
	SetupFlags() []*pflag.Flag
}

type Namer interface {
	Setuper
	CmdName() string
}

type Aliaser interface {
	Setuper
	CmdAliases() []string
}

type Stdouter = plugio.Outer
