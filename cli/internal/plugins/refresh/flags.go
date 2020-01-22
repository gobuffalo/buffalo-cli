package refresh

import (
	"io"

	"github.com/spf13/pflag"
)

func (dev *Developer) Flags() *pflag.FlagSet {
	if dev.flags != nil {
		return dev.flags
	}
	flags := pflag.NewFlagSet(dev.Name(), pflag.ContinueOnError)
	flags.BoolVarP(&dev.help, "help", "h", false, "print this help")
	flags.BoolVarP(&dev.Debug, "debug", "d", false, "turn on delve debugging")
	flags.StringVar(&dev.Config, "config", "", "use a specific config file")
	dev.flags = flags
	return dev.flags
}

func (dev *Developer) DevelopFlags() []*pflag.Flag {
	var values []*pflag.Flag
	flags := dev.Flags()
	flags.VisitAll(func(f *pflag.Flag) {
		values = append(values, f)
	})
	return values
}

func (dev *Developer) PrintFlags(w io.Writer) error {
	flags := dev.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}
