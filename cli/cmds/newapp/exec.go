package newapp

import (
	"context"
	"fmt"

	"github.com/gobuffalo/plugins"
)

func Execute(plugs []plugins.Plugin, ctx context.Context, root string, args []string) error {
	var during []Newapper
	var after []AfterNewapper

	for _, p := range plugs {
		switch t := p.(type) {
		case Newapper:
			during = append(during, t)
		case AfterNewapper:
			after = append(after, t)
		}
	}

	var err error
	for _, p := range during {
		fmt.Println(">>>TODO DURING ", p.PluginName())
		if err = p.Newapp(ctx, root, args); err != nil {
			break
		}
	}

	for _, p := range after {
		fmt.Println(">>>TODO AFTER ", p.PluginName())
		if err := p.AfterNewapp(ctx, root, args, err); err != nil {
			return err
		}
	}
	return err
}
