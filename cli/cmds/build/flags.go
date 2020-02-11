package build

import (
	"io"

	"github.com/gobuffalo/buffalo-cli/v2/internal/flagger"
	"github.com/gobuffalo/plugins/plugflag"
	"github.com/spf13/pflag"
)

func (bc *Cmd) PrintFlags(w io.Writer) error {
	flags := bc.Flags()
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
func (bc *Cmd) Flags() *pflag.FlagSet {
	if bc.flags != nil && bc.flags.Parsed() {
		return bc.flags
	}

	flags := pflag.NewFlagSet(bc.String(), pflag.ContinueOnError)

	flags.BoolVar(&bc.SkipTemplateValidation, "skip-template-validation", false, "skip validating templates")
	flags.BoolVarP(&bc.help, "help", "h", false, "print this help")
	flags.BoolVarP(&bc.Verbose, "verbose", "v", false, "print debugging information")
	flags.BoolVarP(&bc.Static, "static", "s", false, "build a static binary using  --ldflags '-linkmode external -extldflags \"-static\"'")

	flags.StringVar(&bc.LDFlags, "ldflags", "", "set any ldflags to be passed to the go build")
	flags.StringVar(&bc.Mod, "mod", "", "-mod flag for go build")
	flags.StringVarP(&bc.Bin, "output", "o", bc.Bin, "set the name of the binary [default: bin/<module name>]")
	flags.StringVarP(&bc.Environment, "environment", "", "development", "set the environment for the binary")
	flags.StringVarP(&bc.Tags, "tags", "t", "", "compile with specific build tags")

	plugs := bc.ScopedPlugins()

	for _, p := range plugs {
		switch t := p.(type) {
		case Flagger:
			for _, f := range plugflag.Clean(p, t.BuildFlags()) {
				flags.AddGoFlag(f)
			}
		case Pflagger:
			for _, f := range flagger.CleanPflags(p, t.BuildFlags()) {
				flags.AddGoFlag(f)
			}
		}
	}

	bc.flags = flags
	return bc.flags
}
