package cli

import (
	"context"
	"flag"
	"io"
	"os"

	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
	"github.com/gobuffalo/buffalo-cli/internal/v1/cmd"
)

// Buffalo represents the `buffalo` cli.
type Buffalo struct {
	context.Context
	flags   *flag.FlagSet
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Plugins Plugins
	version bool
	help    bool
}

func New(ctx context.Context) (*Buffalo, error) {
	b := &Buffalo{
		Context: ctx,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
	b.setFlags()
	return b, nil
}

func (b *Buffalo) Flags() *flag.FlagSet {
	if b.flags == nil {
		b.setFlags()
	}
	return b.flags
}

func (b *Buffalo) setFlags() {
	b.flags = flag.NewFlagSet("buffalo", flag.ContinueOnError)
	b.flags.BoolVar(&b.version, "v", false, "display version")
	b.flags.BoolVar(&b.help, "h", false, "display help")
	cmdx.Usage(b, b.flags)
}

func (b *Buffalo) Main(ctx context.Context, args []string) error {
	c := cmd.RootCmd
	c.SetArgs(args)
	return c.Execute()

	// flags := b.Flags()
	// if err := flags.Parse(args); err != nil {
	// 	return err
	// }
	// args = flags.Args()
	//
	// if len(args) == 0 {
	// 	flags.Usage()
	// 	return nil
	// }
	//
	// arg := args[0]
	// if len(args) > 0 {
	// 	args = args[1:]
	// }
	//
	// switch arg {
	// case "fix":
	// 	return b.Plugins.Fix(ctx, args)
	// case "generate":
	// 	return b.Plugins.Generate(ctx, args)
	// }
	//
	// if b.version {
	// 	fmt.Fprintln(b.Stdout, "ssh")
	// 	return nil
	// }
	//
	// if b.help {
	// 	flags.Usage()
	// 	return nil
	// }
	// return nil
}
