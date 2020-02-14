package info

import (
	"io"

	"github.com/spf13/pflag"
)

func (cmd *Cmd) PrintFlags(w io.Writer) error {
	flags := cmd.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (cmd *Cmd) Flags() *pflag.FlagSet {
	if cmd.flags != nil {
		return cmd.flags
	}

	flags := pflag.NewFlagSet(cmd.PluginName(), pflag.ContinueOnError)
	flags.BoolVarP(&cmd.help, "help", "h", false, "print this help")

	cmd.flags = flags
	return cmd.flags
}
