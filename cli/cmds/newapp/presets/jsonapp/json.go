package jsonapp

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/clifix"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/bzr"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/fizz"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/flect"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/git"
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

	plugs = append(plugs, clifix.Plugins()...)
	plugs = append(plugs, fizz.Plugins()...)
	plugs = append(plugs, flect.Plugins()...)
	plugs = append(plugs, golang.Plugins()...)
	plugs = append(plugs, grifts.Plugins()...)
	plugs = append(plugs, mail.Plugins()...)
	plugs = append(plugs, pop.Plugins()...)
	plugs = append(plugs, refresh.Plugins()...)
	plugs = append(plugs, soda.Plugins()...)

	root, _ := os.Getwd()
	if _, err := os.Stat(filepath.Join(root, ".git")); err == nil {
		plugs = append(plugs, git.Plugins()...)
	}

	if _, err := os.Stat(filepath.Join(root, ".bzr")); err == nil {
		plugs = append(plugs, bzr.Plugins()...)
	}

	return plugs
}
