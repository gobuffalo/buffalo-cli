package coreapp

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/clifix"
	"github.com/gobuffalo/plugins"
)

func Plugins() []plugins.Plugin {
	plugs := []plugins.Plugin{
		&Generator{},
	}

	plugs = append(plugs, clifix.Plugins()...)

	return plugs
}
