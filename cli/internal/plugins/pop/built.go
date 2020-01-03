package pop

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/internal/plugins"
	"github.com/gobuffalo/pop/v5"
	"github.com/markbates/pkger"
)

type Built struct{}

var _ plugins.Plugin = Built{}

func (Built) Name() string {
	return "pop/built"
}

func (p *Built) BuiltInit(ctx context.Context, args []string) error {
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
