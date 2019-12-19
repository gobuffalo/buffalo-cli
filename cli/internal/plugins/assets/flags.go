package assets

import (
	"io"
	"io/ioutil"

	"github.com/spf13/pflag"
)

func (a *Builder) BuildFlags() []*pflag.Flag {
	var values []*pflag.Flag
	flags := a.Flags()
	flags.VisitAll(func(f *pflag.Flag) {
		values = append(values, f)
	})
	return values
}

func (a *Builder) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(a.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.BoolVar(&a.Clean, "clean-assets", false, "will delete public/assets before calling webpack")
	flags.BoolVarP(&a.Extract, "extract-assets", "e", false, "extract the assets and put them in a distinct archive")
	flags.BoolVarP(&a.Skip, "skip-assets", "k", false, "skip running webpack and building assets")

	return flags
}

func (a *Builder) PrintFlags(w io.Writer) error {
	flags := a.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}
