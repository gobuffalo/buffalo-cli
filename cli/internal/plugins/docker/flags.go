package docker

import (
	"io"
	"io/ioutil"

	"github.com/gobuffalo/buffalo-cli/v2/internal/flagger"
	"github.com/spf13/pflag"
)

func (g *Generator) NewappFlags() []*pflag.Flag {
	return flagger.SetToSlice(g.Flags())
}

func (g *Generator) Flags() *pflag.FlagSet {
	if g.flags != nil {
		return g.flags
	}

	flg := pflag.NewFlagSet(g.PluginName(), pflag.ContinueOnError)
	flg.SetOutput(ioutil.Discard)
	flg.StringVarP(&g.style, "style", "", "multi", "specify the type of Docker file to generate [multi, standard]")

	g.flags = flg
	return g.flags
}

func (g *Generator) PrintFlags(w io.Writer) error {
	flags := g.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}
