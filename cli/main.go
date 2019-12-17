package cli

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/internal/v1/cmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/spf13/pflag"
)

func (b *Buffalo) Main(ctx context.Context, args []string) error {
	var help bool
	flags := pflag.NewFlagSet(b.String(), pflag.ContinueOnError)
	flags.BoolVarP(&help, "help", "h", false, "print this help")
	flags.Parse(args)

	var cmds Commands
	for _, p := range b.Plugins() {
		if c, ok := p.(Command); ok {
			cmds = append(cmds, c)
		}
	}

	ioe := plugins.CtxIO(ctx)
	if len(args) == 0 || (len(flags.Args()) == 0 && help) {
		return plugprint.Print(ioe.Stdout(), b)
	}
	if c, err := cmds.Find(args[0]); err == nil {
		return c.Main(ctx, args[1:])
	}

	c := cmd.RootCmd
	c.SetArgs(args)
	return c.Execute()
}
