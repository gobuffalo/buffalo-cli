package actions

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/resource"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugprint"
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
}

func (Generator) PluginName() string {
	return "pop/action"
}

func (Generator) Description() string {
	return "Generate a Pop action"
}
