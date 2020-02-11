package mail

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/generate"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/mail/internal/mailgen"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"
)

var _ plugcmd.Namer = Generator{}
var _ generate.Generator = Generator{}
var _ plugins.Plugin = Generator{}
var _ plugprint.Describer = Generator{}

type Generator struct{}

func (Generator) PluginName() string {
	return "mail/generator"
}

func (Generator) CmdName() string {
	return "mail"
}

func (Generator) Description() string {
	return "Generate a new mailer for Buffalo"
}

func (Generator) Generate(ctx context.Context, root string, args []string) error {
	run := genny.WetRunner(context.Background())

	opts := &mailgen.Options{
		Args: args,
	}
	gg, err := mailgen.New(opts)
	if err != nil {
		return err
	}
	run.WithGroup(gg)

	g, err := gogen.Fmt(root)
	if err != nil {
		return err
	}
	run.With(g)

	return run.Run()
}
