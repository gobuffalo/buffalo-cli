package generatecmd

import (
	"io"

	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugprint.FlagPrinter = &Command{}

func (bc *Command) PrintFlags(w io.Writer) error {
	flags := bc.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (bc *Command) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(bc.String(), pflag.ContinueOnError)

	flags.BoolVarP(&bc.help, "help", "h", false, "print this help")

	return flags
}
