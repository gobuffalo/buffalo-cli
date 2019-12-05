package cli

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
	"github.com/gobuffalo/buffalo-cli/internal/v1/cmd"
)

func (b *Buffalo) Main(ctx context.Context, args []string) error {
	flags := cmdx.NewFlagSet("buffalo", b.Stderr)
	flags.BoolVarP(&b.help, "help", "h", false, "print this help")
	flags.Parse(args)

	var cmds plugins.Commands
	for _, p := range b.Plugins {
		if c, ok := p.(plugins.Command); ok {
			cmds = append(cmds, c)
		}
	}

	if len(args) == 0 || (len(flags.Args()) == 0 && b.help) {
		plugs := make(plugins.Plugins, len(cmds))
		for i, c := range cmds {
			plugs[i] = c
		}
		return plugins.Print(b.Stderr, "", b, plugs)
	}

	if c, err := cmds.Find(args[0]); err == nil {
		return c.Main(ctx, args[1:])
	}

	c := cmd.RootCmd
	c.SetArgs(args)
	return c.Execute()
}
