package mod

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/plugins"
)

var _ plugins.Plugin = &Initer{}
var _ plugins.Needer = &Initer{}
var _ plugins.Scoper = &Initer{}

type Initer struct {
	pluginsFn plugins.Feeder
}

func (i *Initer) WithPlugins(f plugins.Feeder) {
	i.pluginsFn = f
}

func (i *Initer) ScopedPlugins() []plugins.Plugin {
	if i.pluginsFn == nil {
		return nil
	}

	var plugs []plugins.Plugin
	for _, p := range i.pluginsFn() {
		switch p.(type) {
		case Stderrer:
			plugs = append(plugs, p)
		case Stdiner:
			plugs = append(plugs, p)
		case Stdouter:
			plugs = append(plugs, p)
		case Requirer:
			plugs = append(plugs, p)
		case Replacer:
			plugs = append(plugs, p)
		}
	}

	return plugs
}

func (i Initer) PluginName() string {
	return "go/mod/init"
}

func (i *Initer) ModInit(ctx context.Context, root string, name string) error {
	mr := filepath.Join(root, "go.mod")
	f, err := os.Create(mr)
	if err != nil {
		return plugins.Wrap(i, err)
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("module %s\n", name))

	plugs := i.ScopedPlugins()

	if len(plugs) == 0 {
		return nil
	}

	for _, p := range plugs {
		switch t := p.(type) {
		case Replacer:
			m := t.ModReplace(root)
			for k, v := range m {
				s := fmt.Sprintf("\nreplace %s => %s", k, v)
				f.WriteString(s)
			}
		}
	}
	return nil
}
