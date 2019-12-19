package cli

import "github.com/gobuffalo/buffalo-cli/plugins"

type Aliases interface {
	plugins.Plugin
	Aliases() []string
}
