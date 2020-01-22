package actiontest

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
	"github.com/gobuffalo/here"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = Generator{}
var _ plugprint.Describer = Generator{}
var _ plugprint.FlagPrinter = &Generator{}
var _ resource.ActionTester = &Generator{}
var _ resource.Pflagger = &Generator{}

type Generator struct {
	TestPkg string
	flags   *pflag.FlagSet
	info    here.Info
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
	return "pop/action-test"
}

func (Generator) Description() string {
	return "Generate a Pop action test"
}
