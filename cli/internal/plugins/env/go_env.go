package env

import (
	"context"
	"os"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/test"
	"github.com/gobuffalo/plugins"
)

var _ plugins.Plugin = &GoEnv{}
var _ test.BeforeTester = &GoEnv{}
var _ test.AfterTester = &GoEnv{}

//GoEnv Sets GO_ENV before tests run
type GoEnv struct{}

//PluginName for BeforeTestEnv
func (ebt GoEnv) PluginName() string {
	return "env/tests"
}

//BeforeTest should be invoked before tests run to set the GO_ENV variable
func (ebt *GoEnv) BeforeTest(ctx context.Context, root string, args []string) error {
	return os.Setenv("GO_ENV", "test")
}

//AfterTest should be invoked after tests run to reset GO_ENV variable
func (ebt *GoEnv) AfterTest(ctx context.Context, root string, args []string, err error) error {
	return os.Setenv("GO_ENV", "")
}
