package flect

import (
	"context"
	"fmt"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/plugins"
	"github.com/markbates/pkger"
)

type Initer struct{}

var _ plugins.Plugin = &Initer{}

func (Initer) PluginName() string {
	return "flect"
}

func (fl *Initer) BuiltInit(ctx context.Context, args []string) error {
	f, err := pkger.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to load inflections %s", err)
	}
	defer f.Close()

	err = flect.LoadInflections(f)
	if err != nil {
		return err
	}
	return nil
}
