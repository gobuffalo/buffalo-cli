package pop

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/generatecmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/internal/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/plugins/plugprint"
)

type ModelGen struct{}

var _ plugins.Plugin = ModelGen{}

func (ModelGen) Name() string {
	return "pop/model"
}

var _ plugprint.NamedCommand = ModelGen{}

func (ModelGen) CmdName() string {
	return "model"
}

var _ plugprint.Describer = ModelGen{}

func (ModelGen) Description() string {
	return "Generate a Pop model"
}

var _ generatecmd.Generator = &ModelGen{}

func (mg *ModelGen) Generate(ctx context.Context, args []string) error {
	args = append([]string{"generate", "model"}, args...)
	return Cmd{}.Main(ctx, args)
}

var _ resource.Modeler = &ModelGen{}

func (mg *ModelGen) GenerateResourceModels(ctx context.Context, root string, args []string) error {
	return mg.Generate(ctx, args)
}
