package env

import (
	"context"
	"os"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/test"
	"github.com/gobuffalo/plugins"
)

var _ plugins.Plugin = &BeforeTester{}
var _ test.BeforeTester = &BeforeTester{}

type BeforeTester struct{}

func (ebt BeforeTester) PluginName() string {
	return "env/before-tests"
}

func (ebt *BeforeTester) BeforeTest(ctx context.Context, root string, args []string) error {
	return os.Setenv("GO_ENV", "test")
}
