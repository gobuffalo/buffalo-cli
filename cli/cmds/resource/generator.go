package resource

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/generate"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/markbates/safe"
	"github.com/spf13/pflag"
)

var _ plugcmd.Aliaser = Generator{}
var _ generate.Generator = &Generator{}
var _ plugins.Plugin = &Generator{}
var _ plugins.Needer = &Generator{}
var _ plugins.Scoper = &Generator{}

type Generator struct {
	skipActionTests    bool
	skipActions        bool
	skipMigrationTests bool
	skipMigrations     bool
	skipModelTests     bool
	skipModels         bool
	skipTemplateTests  bool
	skipTemplates      bool

	flags     *pflag.FlagSet
	help      bool
	pluginsFn plugins.Feeder
}

func (g *Generator) WithPlugins(f plugins.Feeder) {
	g.pluginsFn = f
}

func (g Generator) PluginName() string {
	return "resource"
}

func (g Generator) CmdAliases() []string {
	return []string{"r"}
}

func (g *Generator) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if g.pluginsFn != nil {
		plugs = g.pluginsFn()
	}

	var builders []plugins.Plugin

	for _, p := range plugs {
		switch p.(type) {
		case BeforeGenerator:
			builders = append(builders, p)
		case Stdouter:
			builders = append(builders, p)
		case Actioner:
			builders = append(builders, p)
		case ActionTester:
			builders = append(builders, p)
		case Modeler:
			builders = append(builders, p)
		case ModelTester:
			builders = append(builders, p)
		case Templater:
			builders = append(builders, p)
		case TemplateTester:
			builders = append(builders, p)
		case AfterGenerator:
			builders = append(builders, p)
		case Migrationer:
			builders = append(builders, p)
		case MigrationTester:
			builders = append(builders, p)
		case Flagger:
			builders = append(builders, p)
		case Pflagger:
			builders = append(builders, p)
		}
	}

	return builders
}

func (g *Generator) beforeGenerate(ctx context.Context, root string, args []string) error {
	plugs := g.ScopedPlugins()

	for _, p := range plugs {
		if b, ok := p.(BeforeGenerator); ok {
			if err := b.BeforeGenerateResource(ctx, root, args); err != nil {
				return plugins.Wrap(p, err)
			}
		}
	}

	return nil
}

func (g *Generator) addResource(root string, n string) error {
	fp := filepath.Join(root, "actions", "app.go")

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return plugins.Wrap(g, err)
	}

	pres := name.New(n)
	stmt := fmt.Sprintf("app.Resource(\"/%s\", %sResource{})", pres.URL(), pres.Resource())

	gf, err := gogen.AddInsideBlock(genny.NewFileB(fp, b), "if app == nil {", stmt)
	if err != nil {
		return plugins.Wrap(g, err)
	}

	f, err := os.Create(fp)
	if err != nil {
		return plugins.Wrap(g, err)
	}
	defer f.Close()

	_, err = f.WriteString(gf.String())
	if err != nil {
		return plugins.Wrap(g, err)
	}

	return nil
}

func (g *Generator) afterGenerate(ctx context.Context, root string, args []string, err error) error {
	plugs := g.ScopedPlugins()

	if err == nil && len(args) > 0 {
		if err := g.addResource(root, args[0]); err != nil {
			return plugins.Wrap(g, err)
		}
	}

	for _, p := range plugs {
		if b, ok := p.(AfterGenerator); ok {
			if err := b.AfterGenerateResource(ctx, root, args, err); err != nil {
				return plugins.Wrap(b, err)
			}
		}
	}

	return nil
}

// Generate implements generate.Generator and is the entry point for `buffalo generate resource`
func (g *Generator) Generate(ctx context.Context, root string, args []string) error {
	if len(args) == 0 {
		err := fmt.Errorf("you must specify a name for the resource")
		return plugins.Wrap(g, err)
	}

	flags := g.Flags()
	if err := flags.Parse(args); err != nil {
		return plugins.Wrap(g, err)
	}

	args = flags.Args()

	plugs := g.ScopedPlugins()

	if g.help {
		return plugprint.Print(plugio.Stdout(plugs...), g)
	}

	err := g.run(ctx, root, args)
	return g.afterGenerate(ctx, root, args, err)
}

func (g *Generator) run(ctx context.Context, root string, args []string) error {
	if err := g.beforeGenerate(ctx, root, args); err != nil {
		return plugins.Wrap(g, err)
	}

	type step func(context.Context, string, []string) error

	steps := []step{
		g.generateActionTests,
		g.generateActions,
		g.generateMigrationTests,
		g.generateMigrations,
		g.generateModelTests,
		g.generateModels,
		g.generateTemplateTests,
		g.generateTemplates,
	}

	for i, step := range steps {
		if err := step(ctx, root, args); err != nil {
			return fmt.Errorf("(%d) %w", i, plugins.Wrap(g, err))
		}
	}

	return nil
}

func (g *Generator) generateActions(ctx context.Context, root string, args []string) error {
	if g.skipActions {
		return nil
	}

	for _, p := range g.ScopedPlugins() {
		if ag, ok := p.(Actioner); ok {
			if err := ag.GenerateResourceActions(ctx, root, args); err != nil {
				return plugins.Wrap(ag, err)
			}
		}
	}

	return nil
}

func (g *Generator) generateActionTests(ctx context.Context, root string, args []string) error {
	if g.skipActionTests {
		return nil
	}

	for _, p := range g.ScopedPlugins() {
		if ag, ok := p.(ActionTester); ok {
			if err := ag.GenerateResourceActionTests(ctx, root, args); err != nil {
				return plugins.Wrap(ag, err)
			}
		}
	}

	return nil
}

func (g *Generator) generateTemplates(ctx context.Context, root string, args []string) error {
	if g.skipTemplates {
		return nil
	}

	for _, p := range g.ScopedPlugins() {
		if ag, ok := p.(Templater); ok {
			if err := ag.GenerateResourceTemplates(ctx, root, args); err != nil {
				return plugins.Wrap(ag, err)
			}
		}
	}

	return nil
}

func (g *Generator) generateTemplateTests(ctx context.Context, root string, args []string) error {
	if g.skipTemplateTests {
		return nil
	}

	for _, p := range g.ScopedPlugins() {
		if ag, ok := p.(TemplateTester); ok {
			if err := ag.GenerateResourceTemplateTests(ctx, root, args); err != nil {
				return plugins.Wrap(ag, err)
			}
		}
	}

	return nil
}

func (g *Generator) generateModels(ctx context.Context, root string, args []string) error {
	if g.skipModels {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(Modeler)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] Modeler %s\n", p.PluginName())
			return ag.GenerateResourceModels(ctx, root, args)
		})
	}
	return nil
}

func (g *Generator) generateModelTests(ctx context.Context, root string, args []string) error {
	if g.skipModelTests {
		return nil
	}

	for _, p := range g.ScopedPlugins() {
		if ag, ok := p.(ModelTester); ok {
			if err := ag.GenerateResourceModelTests(ctx, root, args); err != nil {
				return plugins.Wrap(ag, err)
			}
		}
	}

	return nil
}

func (g *Generator) generateMigrations(ctx context.Context, root string, args []string) error {
	if g.skipMigrations {
		return nil
	}

	for _, p := range g.ScopedPlugins() {
		if ag, ok := p.(Migrationer); ok {
			if err := ag.GenerateResourceMigrations(ctx, root, args); err != nil {
				return plugins.Wrap(ag, err)
			}
		}
	}

	return nil
}

func (g *Generator) generateMigrationTests(ctx context.Context, root string, args []string) error {
	if g.skipMigrationTests {
		return nil
	}

	for _, p := range g.ScopedPlugins() {
		if ag, ok := p.(MigrationTester); ok {
			if err := ag.GenerateResourceMigrationTests(ctx, root, args); err != nil {
				return plugins.Wrap(ag, err)
			}
		}
	}

	return nil
}
