package build

import (
	"io"

	"github.com/gobuffalo/buffalo-cli/v2/internal/flagger"
	"github.com/gobuffalo/plugins/plugflag"
	"github.com/spf13/pflag"
)

func (cmd *Cmd) PrintFlags(w io.Writer) error {
	flags := cmd.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

// Flags returns a defined set of flags for this command.
// It imports flags provided by plugins that use either
// the `Flagger` or `Pflagger` interfaces. Flags provided
// by plugins will have their shorthand ("-x") flag stripped
// and the name ("--some-flag") of the flag will be
// prefixed with the plugin's name ("--xyz-some-flag")
func (cmd *Cmd) Flags() *pflag.FlagSet {
	if cmd.flags != nil {
		return cmd.flags
	}

	flags := pflag.NewFlagSet(cmd.PluginName(), pflag.ContinueOnError)

	flags.BoolVarP(&cmd.help, "help", "h", false, "print this help")
	flags.BoolVarP(&cmd.verbose, "verbose", "v", false, "print debugging information")
	flags.BoolVarP(&cmd.static, "static", "s", false, "build a static binary using  --ldflags '-linkmode external -extldflags \"-static\"'")

	flags.StringVar(&cmd.ldFlags, "ldflags", "", "set any ldflags to be passed to the go build")
	flags.StringVar(&cmd.mod, "mod", "", "-mod flag for go build")
	flags.StringVarP(&cmd.bin, "output", "o", cmd.bin, "set the name of the binary [default: bin/<module name>]")
	flags.StringVarP(&cmd.environment, "environment", "", "development", "set the environment for the binary")
	flags.StringVarP(&cmd.tags, "tags", "t", "", "compile with specific build tags")

	plugs := cmd.ScopedPlugins()

	for _, p := range plugs {
		switch t := p.(type) {
		case Flagger:
			for _, f := range plugflag.Clean(p, t.BuildFlags()) {
				flags.AddGoFlag(f)
			}
		case Pflagger:
			for _, f := range flagger.CleanPflags(p, t.BuildFlags()) {
				flags.AddFlag(f)
			}
		}
	}

	cmd.flags = flags
	return cmd.flags
}
