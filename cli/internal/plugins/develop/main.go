package develop

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"golang.org/x/sync/errgroup"
)

func (cmd *Cmd) SubCommand(ctx context.Context, name string, args []string) error {
	cmds := plugins.Commands(cmd.SubCommands())
	p, err := cmds.Find(name)
	if err != nil {
		return err
	}

	d, ok := p.(Developer)
	if !ok {
		return fmt.Errorf("%s is not a developer", name)
	}

	info, err := cmd.HereInfo()
	if err != nil {
		return err
	}
	return d.Develop(ctx, info.Dir, args)
}

func (cmd *Cmd) Main(ctx context.Context, args []string) error {
	flags := cmd.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

	args = flags.Args()

	if len(args) == 0 && cmd.help {
		ioe := plugins.CtxIO(ctx)
		return plugprint.Print(ioe.Stdout(), cmd)
	}

	if len(args) > 0 {
		return cmd.SubCommand(ctx, args[0], args[1:])
	}

	info, err := cmd.HereInfo()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := &errgroup.Group{}

	for _, p := range cmd.ScopedPlugins() {
		if d, ok := p.(Developer); ok {
			wg.Go(func() error {
				return d.Develop(ctx, info.Dir, args)
			})
		}
	}

	return wg.Wait()
}
