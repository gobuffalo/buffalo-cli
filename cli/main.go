package cli

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
	"github.com/gobuffalo/buffalo-cli/internal/v1/cmd"
)

func (b *Buffalo) Main(ctx context.Context, args []string) error {
	flags := cmdx.NewFlagSet("buffalo", b.Stderr)
	flags.BoolVarP(&b.help, "help", "h", false, "print this help")
	flags.Parse(args)

	cmds := b.Plugins.Commands()

	if len(args) == 0 || b.help {
		return b.usage(b.Stderr, cmds)
	}

	if c, err := cmds.Find(args[0]); err == nil {
		return c.Main(ctx, args[1:])
	}

	c := cmd.RootCmd
	c.SetArgs(args)
	return c.Execute()
}

func (b *Buffalo) usage(w io.Writer, cmds plugins.Commands) error {
	fmt.Fprintln(w, strings.TrimSpace(usageTmpl))

	const ac = "\nAvailable Commands:\n"
	if len(cmds) == 0 {
		return nil
	}
	fmt.Fprint(w, ac)

	scmds := make(plugins.Commands, len(cmds))
	copy(scmds, cmds)
	sort.Sort(scmds)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, c := range scmds {
		fmt.Fprintf(tw, "  %s\t%s\n", c.Name(), plugins.Description(c))
	}
	tw.Flush()
	return nil
}

const usageTmpl = `
Build Buffalo applications with ease

Usage:
  buffalo [command]

`
