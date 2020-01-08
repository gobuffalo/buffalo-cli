package modelgen

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/generatecmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/soda"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
)

type Generator struct{}

var _ plugins.Plugin = Generator{}

func (Generator) Name() string {
	return "pop/model"
}

var _ plugprint.NamedCommand = Generator{}

func (Generator) CmdName() string {
	return "model"
}

var _ plugprint.Describer = Generator{}

func (Generator) Description() string {
	return "Generate a Pop model"
}

var _ generatecmd.Generator = &Generator{}

func (mg *Generator) Generate(ctx context.Context, args []string) error {
	args = append([]string{"generate", "model"}, args...)
	return soda.Main(ctx, args)
}

var _ resource.Modeler = &Generator{}

func (mg *Generator) GenerateResourceModels(ctx context.Context, root string, args []string) error {
	return mg.Generate(ctx, args)
}
