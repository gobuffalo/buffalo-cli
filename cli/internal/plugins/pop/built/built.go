package built

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/pop/v5"
	"github.com/markbates/pkger"
)

var _ plugins.Plugin = Initer{}

type Initer struct{}

func (Initer) PluginName() string {
	return "pop/built/initer"
}

func (p *Initer) BuiltInit(ctx context.Context, root string, args []string) error {
	f, err := pkger.Open("/database.yml")
	if err != nil {
		return err
	}
	defer f.Close()

	err = pop.LoadFrom(f)
	if err != nil {
		return err
	}
	return nil
}
