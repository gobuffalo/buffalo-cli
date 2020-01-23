package generate

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
)

// Generator is a sub-command of buffalo generate.
// 	buffalo generate model
type Generator interface {
	plugins.Plugin
	Generate(ctx context.Context, root string, args []string) error
}
