package grifts

import (
	"context"
	"os"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/generate"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/grifts/internal/griftgen"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

var _ generate.Generator = Generator{}
var _ plugins.NamedCommand = Generator{}
var _ plugins.Plugin = Generator{}
var _ plugprint.Describer = Generator{}

type Generator struct{}

func (Generator) Name() string {
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

func (Generator) Generate(ctx context.Context, args []string) error {
	run := genny.WetRunner(context.Background())

	opts := &griftgen.Options{}
	opts.Args = args
	g, err := griftgen.New(opts)
	if err != nil {
		return err
	}
	run.With(g)

	pwd, _ := os.Getwd()
	g, err = gogen.Fmt(pwd)
	if err != nil {
		return err
	}
	run.With(g)

	return run.Run()
}
