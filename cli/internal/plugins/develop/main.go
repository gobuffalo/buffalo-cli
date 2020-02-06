package develop

import (
	"context"
	"fmt"
	"strings"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
	"golang.org/x/sync/errgroup"
)

func (cmd *Cmd) SubCommand(ctx context.Context, root string, name string, args []string) error {
	cmds := plugins.Commands(cmd.SubCommands())
	p, err := cmds.Find(name)
	if err != nil {
		return err
	}

	d, ok := p.(Developer)
	if !ok {
		return fmt.Errorf("%s is not a developer", name)
	}

	return d.Develop(ctx, root, args)
}

func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	if len(args) == 1 && args[0] == "-h" {
		ioe := plugins.CtxIO(ctx)
		return plugprint.Print(ioe.Stdout(), cmd)
	}

	if len(args) > 0 {
		for _, n := range args {
			if strings.HasPrefix(n, "-") {
				continue
			}
			return cmd.SubCommand(ctx, root, n, args[1:])
		}

	}

	flags := cmd.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

	args = flags.Args()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := &errgroup.Group{}

	for _, p := range cmd.ScopedPlugins() {
		if d, ok := p.(Developer); ok {
			wg.Go(func() error {
				return d.Develop(ctx, root, args)
			})
		}
	}

	return wg.Wait()
}
