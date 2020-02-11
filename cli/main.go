package cli

import (
	"context"
	"fmt"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugfind"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/markbates/safe"
	"github.com/spf13/pflag"
)

// Main is the entry point of the `buffalo` command
func (b *Buffalo) Main(ctx context.Context, root string, args []string) error {
	var help bool
	flags := pflag.NewFlagSet(b.String(), pflag.ContinueOnError)
	flags.BoolVarP(&help, "help", "h", false, "print this help")
	flags.Parse(args)

	pfn := func() []plugins.Plugin {
		return b.Plugins
	}

	plugs := b.Plugins
	for _, p := range plugs {
		switch t := p.(type) {
		case Needer:
			t.WithPlugins(pfn)
		case StdinNeeder:
			if err := t.SetStdin(plugio.Stdin(plugs...)); err != nil {
				return err
			}
		case StdoutNeeder:
			if err := t.SetStdout(plugio.Stdout(plugs...)); err != nil {
				return err
			}
		case StderrNeeder:
			if err := t.SetStderr(plugio.Stderr(plugs...)); err != nil {
				return err
			}
		}
	}

	if len(args) == 0 || (len(flags.Args()) == 0 && help) {
		return plugprint.Print(plugio.Stdout(plugs...), b)
	}

	name := args[0]

	fn := plugfind.Background()
	fn = plugcmd.ByCommander(fn)
	fn = plugcmd.ByNamer(fn)
	fn = plugcmd.ByAliaser(fn)

	p := fn.Find(name, plugs)

	c, ok := p.(Commander)
	if !ok {
		return fmt.Errorf("unknown command %s", name)
	}

	return safe.RunE(func() error {
		return c.Main(ctx, root, args[1:])
	})
}
