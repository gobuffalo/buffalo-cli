package setup

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop/setup"
	"github.com/gobuffalo/plugins"
	"github.com/markbates/grift/cli"
)

var _ plugins.Plugin = &Setup{}
var _ setup.DBSeeder = &Setup{}

type Setup struct {
}

func (s *Setup) PluginName() string {
	return "grifts/setup"
}

func (s *Setup) SeedDB(ctx context.Context, root string, app []string) error {
	return cli.Run(ctx, []string{"db:seed"})
}
