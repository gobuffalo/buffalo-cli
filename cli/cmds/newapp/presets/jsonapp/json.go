package jsonapp

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp/presets/coreapp"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/clifix"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/ci"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/docker"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/fizz"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/flect"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/golang"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/grifts"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/mail"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/refresh"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/soda"
	"github.com/gobuffalo/plugins"
)

func Plugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, coreapp.Plugins()...)

	plugs = append(plugs, clifix.Plugins()...)
	plugs = append(plugs, fizz.Plugins()...)
	plugs = append(plugs, flect.Plugins()...)
	plugs = append(plugs, golang.Plugins()...)
	plugs = append(plugs, grifts.Plugins()...)
	plugs = append(plugs, mail.Plugins()...)
	plugs = append(plugs, pop.Plugins()...)
	plugs = append(plugs, refresh.Plugins()...)
	plugs = append(plugs, soda.Plugins()...)
	plugs = append(plugs, docker.Plugins()...)
	plugs = append(plugs, ci.Plugins()...)

	return plugs
}
