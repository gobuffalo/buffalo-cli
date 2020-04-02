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
	SkipActionTests    bool
	SkipActions        bool
	SkipMigrationTests bool
	SkipMigrations     bool
	SkipModelTests     bool
	SkipModels         bool
	SkipTemplateTests  bool
	SkipTemplates      bool

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
		case ResourceGenerator:
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
		}
	}

	return builders
}

func (g *Generator) beforeGenerate(ctx context.Context, root string, args []string) error {
	plugs := g.ScopedPlugins()

	for _, p := range plugs {
		b, ok := p.(BeforeGenerator)
		if !ok {
			continue
		}
		err := safe.RunE(func() error {
			fmt.Printf("[Resource] BeforeGenerator %s\n", p.PluginName())
			return b.BeforeGenerateResource(ctx, root, args)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) addResource(root string, n string) error {
	fp := filepath.Join(root, "actions", "app.go")

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}

	pres := name.New(n)
	stmt := fmt.Sprintf("app.Resource(\"/%s\", %sResource{})", pres.URL(), pres.Resource())

	gf, err := gogen.AddInsideBlock(genny.NewFileB(fp, b), "if app == nil {", stmt)
	if err != nil {
		return err
	}

	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(gf.String())
	if err != nil {
		return err
	}

	return nil
}
func (g *Generator) afterGenerate(ctx context.Context, root string, args []string, err error) error {
	plugs := g.ScopedPlugins()

	if err == nil && len(args) > 0 {
		if err := g.addResource(root, args[0]); err != nil {
			return err
		}
	}

	for _, p := range plugs {
		b, ok := p.(AfterGenerator)
		if !ok {
			continue
		}
		err := safe.RunE(func() error {
			fmt.Printf("[Resource] AfterGenerator %s\n", p.PluginName())
			return b.AfterGenerateResource(ctx, root, args, err)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Generate implements generate.Generator and is the entry point for `buffalo generate resource`
func (g *Generator) Generate(ctx context.Context, root string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must specify a name for the resource")
	}

	flags := g.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

	args = flags.Args()

	plugs := g.ScopedPlugins()

	if g.help {
		return plugprint.Print(plugio.Stdout(plugs...), g)
	}

	if err := g.beforeGenerate(ctx, root, args); err != nil {
		return g.afterGenerate(ctx, root, args, err)
	}

	for _, p := range plugs {
		gr, ok := p.(ResourceGenerator)
		if !ok {
			continue
		}
		err := safe.RunE(func() error {
			fmt.Printf("[Resource] ResourceGenerator %s\n", p.PluginName())
			return gr.GenerateResource(ctx, root, args)
		})
		return g.afterGenerate(ctx, root, args, err)
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

	for _, step := range steps {
		if err := step(ctx, root, args); err != nil {
			return err
		}
	}

	return g.afterGenerate(ctx, root, args, nil)

}

func (g *Generator) generateActions(ctx context.Context, root string, args []string) error {
	if g.SkipActions {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(Actioner)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] Actioner %s\n", p.PluginName())
			return ag.GenerateResourceActions(ctx, root, args)
		})
	}
	return nil
}

func (g *Generator) generateActionTests(ctx context.Context, root string, args []string) error {
	if g.SkipActionTests {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(ActionTester)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] ActionTester %s\n", p.PluginName())
			return ag.GenerateResourceActionTests(ctx, root, args)
		})
	}
	return nil
}

func (g *Generator) generateTemplates(ctx context.Context, root string, args []string) error {
	if g.SkipTemplates {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(Templater)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] Templater %s\n", p.PluginName())
			return ag.GenerateResourceTemplates(ctx, root, args)
		})
	}
	return nil
}

func (g *Generator) generateTemplateTests(ctx context.Context, root string, args []string) error {
	if g.SkipTemplateTests {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(TemplateTester)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] TemplateTester %s\n", p.PluginName())
			return ag.GenerateResourceTemplateTests(ctx, root, args)
		})
	}
	return nil
}

func (g *Generator) generateModels(ctx context.Context, root string, args []string) error {
	if g.SkipModels {
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
	if g.SkipModelTests {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(ModelTester)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] ModelTester %s\n", p.PluginName())
			return ag.GenerateResourceModelTests(ctx, root, args)
		})
	}
	return nil
}

func (g *Generator) generateMigrations(ctx context.Context, root string, args []string) error {
	if g.SkipMigrations {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(Migrationer)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] Migrationer %s\n", p.PluginName())
			return ag.GenerateResourceMigrations(ctx, root, args)
		})
	}
	return nil
}

func (g *Generator) generateMigrationTests(ctx context.Context, root string, args []string) error {
	if g.SkipMigrationTests {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(MigrationTester)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] MigrationTester %s\n", p.PluginName())
			return ag.GenerateResourceMigrationTests(ctx, root, args)
		})
	}
	return nil
}
