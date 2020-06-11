package ci

import (
	"context"

	"github.com/spf13/pflag"
)

type Generator struct {
	provider string
	flags    *pflag.FlagSet
}

func (Generator) PluginName() string {
	return "ci"
}

func (Generator) Description() string {
	return "Generates CI configuration file"
}

func (g Generator) Newapp(ctx context.Context, root string, name string, args []string) error {
	return nil
}
