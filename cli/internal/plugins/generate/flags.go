package generate

import (
	"io"

	"github.com/spf13/pflag"
)

func (bc *Cmd) PrintFlags(w io.Writer) error {
	flags := bc.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (bc *Cmd) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(bc.String(), pflag.ContinueOnError)

	flags.BoolVarP(&bc.help, "help", "h", false, "print this help")

	return flags
}
