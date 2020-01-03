package grifts

import (
	"context"
	"os"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/grifts/internal/griftgen"
	"github.com/gobuffalo/buffalo-cli/internal/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/plugins/plugprint"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

type Generator struct{}

var _ plugins.Plugin = Generator{}

func (Generator) Name() string {
	return "grifts/generator"
}

var _ plugins.NamedCommand = Generator{}

func (Generator) CmdName() string {
	return "grift"
}

var _ plugins.Aliases = Buffalo{}

func (Generator) Aliases() []string {
	return []string{"task"}
}

var _ plugprint.Describer = Generator{}

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
