package generate

import (
	"context"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
)

// Generator is a sub-command of buffalo generate.
// 	buffalo generate model
type Generator interface {
	plugins.Plugin
	Generate(ctx context.Context, root string, args []string) error
}

type Namer interface {
	Generator
	CmdName() string
}

type Aliaser interface {
	Generator
	CmdAliases() []string
}

type Stdouter = plugio.Outer
