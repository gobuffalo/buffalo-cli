package generator

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/generate"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"
)

var _ generate.Generator = &Generator{}
var _ plugcmd.Aliaser = Generator{}
var _ plugcmd.Namer = Generator{}
var _ plugins.Plugin = Generator{}
var _ plugprint.Describer = Generator{}

type Generator struct{}

func (Generator) PluginName() string {
	return "grifts/generator"
}

func (Generator) CmdName() string {
	return "grift"
}

func (Generator) CmdAliases() []string {
	return []string{"task"}
}

func (Generator) Description() string {
	return "Generate a grift task"
}

func (gen *Generator) Generate(ctx context.Context, root string, args []string) error {
	run := genny.WetRunner(context.Background())

	opts := &Options{}
	opts.Args = args
	g, err := New(opts)
	if err != nil {
		return plugins.Wrap(gen, err)
	}
	run.With(g)

	g, err = gogen.Fmt(root)
	if err != nil {
		return plugins.Wrap(gen, err)
	}
	run.With(g)

	return run.Run()
}
