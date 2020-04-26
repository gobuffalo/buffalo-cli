package newapp

import (
	"io"
	"io/ioutil"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp"
	"github.com/gobuffalo/buffalo-cli/v2/internal/flagger"
	"github.com/spf13/pflag"
)

var _ newapp.Pflagger = &Generator{}

func (a *Generator) NewappFlags() []*pflag.Flag {
	return flagger.SetToSlice(a.Flags())
}

func (a *Generator) Flags() *pflag.FlagSet {
	if a.flags != nil {
		return a.flags
	}

	flags := pflag.NewFlagSet(a.PluginName(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.StringVarP(&a.tool, "tool", "t", "yarnpkg", "asset tool to install dependencies")

	a.flags = flags
	return a.flags
}

func (a *Generator) PrintFlags(w io.Writer) error {
	flags := a.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}
