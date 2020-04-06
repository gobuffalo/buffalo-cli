package build

import (
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/internal/flagger"
	"github.com/spf13/pflag"
)

func (a *Builder) BuildFlags() []*pflag.Flag {
	return flagger.SetToSlice(a.Flags())
}

func (a *Builder) Flags() *pflag.FlagSet {
	if a.flags != nil {
		return a.flags
	}

	flags := pflag.NewFlagSet(a.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.BoolVar(&a.clean, "clean", false, "will delete public/webpack before calling webpack")
	flags.BoolVarP(&a.extract, "extract", "e", false, "extract the webpack and put them in a distinct archive")
	flags.BoolVarP(&a.skip, "skip", "k", false, "skip running webpack and building webpack")

	flags.StringVar(&a.extractTo, "extract-to", filepath.Join("bin", "webpack.zip"), "extract the webpack and put them in a distinct archive")

	a.flags = flags
	return a.flags
}

func (a *Builder) PrintFlags(w io.Writer) error {
	flags := a.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}
