package actions

import (
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = Generator{}
var _ plugprint.Describer = Generator{}
var _ plugprint.FlagPrinter = &Generator{}
var _ resource.Pflagger = &Generator{}
var _ resource.Actioner = &Generator{}

type Generator struct {
	ModelName    string
	ModelsPkg    string
	ModelsPkgSel string
	flags        *pflag.FlagSet
	info         here.Info
}

func (g *Generator) WithHereInfo(i here.Info) {
	g.info = i
}

func (g *Generator) HereInfo() (here.Info, error) {
	if g.info.IsZero() {
		return here.Current()
	}
	return g.info, nil
}

func (Generator) Name() string {
	return "pop/action"
}

func (Generator) Description() string {
	return "Generate a Pop action"
}
