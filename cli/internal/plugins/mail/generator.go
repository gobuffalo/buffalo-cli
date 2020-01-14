package mail

import (
	"context"
	"os"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/mail/internal/mailgen"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

type Generator struct{}

var _ plugins.Plugin = Generator{}

func (Generator) Name() string {
	return "mail/generator"
}

var _ plugins.NamedCommand = Generator{}

func (Generator) CmdName() string {
	return "mail"
}

var _ plugprint.Describer = Generator{}

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
