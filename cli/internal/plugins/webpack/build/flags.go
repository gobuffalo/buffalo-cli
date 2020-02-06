package build

import (
	"io"
	"io/ioutil"
	"path/filepath"

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
	if a.flags != nil && a.flags.Parsed() {
		return a.flags
	}

	flags := pflag.NewFlagSet(a.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.BoolVar(&a.Clean, "clean", false, "will delete public/webpack before calling webpack")
	flags.BoolVarP(&a.Extract, "extract", "e", false, "extract the webpack and put them in a distinct archive")
	flags.BoolVarP(&a.Skip, "skip", "k", false, "skip running webpack and building webpack")

	flags.StringVar(&a.ExtractTo, "extract-to", filepath.Join("bin", "webpack.zip"), "extract the webpack and put them in a distinct archive")

	a.flags = flags
	return a.flags
}

func (a *Builder) PrintFlags(w io.Writer) error {
	flags := a.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}
