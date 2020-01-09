package actiontest

import (
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
)

type Generator struct {
	info    here.Info
	testPkg string
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

var _ plugins.Plugin = Generator{}

func (Generator) Name() string {
	return "pop/action-test"
}

var _ plugprint.Describer = Generator{}

func (Generator) Description() string {
	return "Generate a Pop action test"
}
