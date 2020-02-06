package grifts

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/generate"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/grifts/internal/griftgen"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

var _ generate.Generator = Generator{}
var _ plugins.NamedCommand = Generator{}
var _ plugins.Plugin = Generator{}
var _ plugprint.Describer = Generator{}

type Generator struct{}

func (Generator) PluginName() string {
	return "grifts/generator"
}

func (Generator) CmdName() string {
	return "grift"
}

func (Generator) Aliases() []string {
	return []string{"task"}
}

func (Generator) Description() string {
	return "Generate a grift task"
}

func (Generator) Generate(ctx context.Context, root string, args []string) error {
	run := genny.WetRunner(context.Background())

	opts := &griftgen.Options{}
	opts.Args = args
	g, err := griftgen.New(opts)
	if err != nil {
		return err
	}
	run.With(g)

	g, err = gogen.Fmt(root)
	if err != nil {
		return err
	}
	run.With(g)

	return run.Run()
}
