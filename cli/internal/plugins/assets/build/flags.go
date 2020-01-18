package build

import (
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ buildcmd.Pflagger = &Builder{}

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
	flags.BoolVar(&a.Clean, "clean", false, "will delete public/assets before calling webpack")
	flags.BoolVarP(&a.Extract, "extract", "e", false, "extract the assets and put them in a distinct archive")
	flags.BoolVarP(&a.Skip, "skip", "k", false, "skip running webpack and building assets")

	flags.StringVar(&a.ExtractTo, "extract-to", filepath.Join("bin", "assets.zip"), "extract the assets and put them in a distinct archive")
	return flags
}

var _ plugprint.FlagPrinter = &Builder{}

func (a *Builder) PrintFlags(w io.Writer) error {
	flags := a.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}
