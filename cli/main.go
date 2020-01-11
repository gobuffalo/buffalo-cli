package cli

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/markbates/safe"
	"github.com/spf13/pflag"
)

// Main is the entry point of the `buffalo` command
func (b *Buffalo) Main(ctx context.Context, args []string) error {
	var help bool
	flags := pflag.NewFlagSet(b.String(), pflag.ContinueOnError)
	flags.BoolVarP(&help, "help", "h", false, "print this help")
	flags.Parse(args)

	var cmds Commands
	for _, p := range b.ScopedPlugins() {
		if c, ok := p.(Command); ok {
			cmds = append(cmds, c)
		}
	}

	ioe := plugins.CtxIO(ctx)
	if len(args) == 0 || (len(flags.Args()) == 0 && help) {
		return plugprint.Print(ioe.Stdout(), b)
	}

	name := args[0]
	if c, err := cmds.Find(name); err == nil {
		return safe.RunE(func() error {
			return c.Main(ctx, args[1:])
		})
	}
	return fmt.Errorf("unknown command %s", name)
}
