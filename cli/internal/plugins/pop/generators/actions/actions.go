package actions

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/resource"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = Generator{}
var _ plugprint.Describer = Generator{}
var _ plugprint.FlagPrinter = &Generator{}
var _ resource.Actioner = &Generator{}
var _ resource.Pflagger = &Generator{}

type Generator struct {
	modelName    string
	modelsPkg    string
	modelsPkgSel string
	flags        *pflag.FlagSet
}

func (Generator) PluginName() string {
	return "pop/action"
}

func (Generator) Description() string {
	return "Generate a Pop action"
}
