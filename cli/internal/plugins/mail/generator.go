package mail

import (
	"context"
	"os"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/generate"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/mail/internal/mailgen"
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

func (Generator) Name() string {
	return "mail/generator"
}

func (Generator) CmdName() string {
	return "mail"
}

func (Generator) Description() string {
	return "Generate a new mailer for Buffalo"
}

func (Generator) Generate(ctx context.Context, args []string) error {
	run := genny.WetRunner(context.Background())

	opts := &mailgen.Options{
		Args: args,
	}
	gg, err := mailgen.New(opts)
	if err != nil {
		return err
	}
	run.WithGroup(gg)

	pwd, _ := os.Getwd()
	g, err := gogen.Fmt(pwd)
	if err != nil {
		return err
	}
	run.With(g)

	return run.Run()
}
