package cmds

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/develop"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/fix"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/generate"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/info"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/resource"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/test"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/version"
	"github.com/gobuffalo/plugins"
)

func Plugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, build.Plugins()...)
	plugs = append(plugs, develop.Plugins()...)
	plugs = append(plugs, fix.Plugins()...)
	plugs = append(plugs, generate.Plugins()...)
	plugs = append(plugs, info.Plugins()...)
	plugs = append(plugs, resource.Plugins()...)
	plugs = append(plugs, test.Plugins()...)
	plugs = append(plugs, version.Plugins()...)
	return plugs
}
