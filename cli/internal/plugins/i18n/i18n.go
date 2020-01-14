package i18n

import (
	"context"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/flect/name"
)

type Generator struct {
	pluginsFn plugins.PluginFeeder
}

var _ plugins.PluginNeeder = &Generator{}

func (g *Generator) WithPlugins(f plugins.PluginFeeder) {
	g.pluginsFn = f
}

func (g *Generator) ScopedPlugins() []plugins.Plugin {
	if g.pluginsFn == nil {
		return []plugins.Plugin{}
	}
	plugs := g.pluginsFn()

	var scoped []plugins.Plugin
	for _, p := range plugs {
		switch p.(type) {
		case NamedWriter:
			scoped = append(scoped, p)
		}
	}

	return scoped
}

var _ plugins.Plugin = &Generator{}

func (g *Generator) Name() string {
	return "i18n"
}

var _ resource.AfterGenerator = &Generator{}

func (g *Generator) AfterGenerateResource(ctx context.Context, root string, args []string, err error) error {
	if err != nil || len(args) == 0 {
		return nil
	}

	model := name.New(args[0])

	fp := filepath.Join(root, "locales", model.Resource().File().String()+".en-us.yaml")
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return err
	}

	if err := t.Execute(f, model); err != nil {
		return err
	}

	flash := &flasher{}
	flash.WithPlugins(g.ScopedPlugins)
	if err := flash.Flash(ctx, root, model); err != nil {
		return err
	}
	return nil
}

const tmpl = `- id: "{{.Singularize.Underscore}}.created.success"
  translation: "{{.Proper}} was successfully created."
- id: "{{.Singularize.Underscore}}.updated.success"
  translation: "{{.Proper}} was successfully updated."
- id: "{{.Singularize.Underscore}}.destroyed.success"
  translation: "{{.Proper}} was successfully destroyed."
`
