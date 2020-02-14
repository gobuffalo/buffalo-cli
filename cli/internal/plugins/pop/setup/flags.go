package setup

import (
	"io"

	"github.com/gobuffalo/buffalo-cli/v2/internal/flagger"
	"github.com/spf13/pflag"
)

func (setup *Setup) PrintFlags(w io.Writer) error {
	flags := setup.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (setup *Setup) Flags() *pflag.FlagSet {
	if setup.flags != nil {
		return setup.flags
	}

	flags := pflag.NewFlagSet(setup.PluginName(), pflag.ContinueOnError)
	flags.BoolVarP(&setup.dropDB, "drop-db", "d", false, "drop database before creating them")
	flags.BoolVarP(&setup.help, "help", "h", false, "print this help")

	setup.flags = flags
	return setup.flags
}

func (setup *Setup) SetupFlags() []*pflag.Flag {
	return flagger.SetToSlice(setup.Flags())
}
