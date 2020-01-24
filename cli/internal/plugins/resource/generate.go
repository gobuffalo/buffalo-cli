package resource

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/generate"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/markbates/safe"
	"github.com/spf13/pflag"
)

var _ generate.Generator = &Generator{}
var _ plugins.Aliases = Generator{}
var _ plugins.Plugin = &Generator{}
var _ plugins.PluginNeeder = &Generator{}
var _ plugins.PluginScoper = &Generator{}

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
	pluginsFn plugins.PluginFeeder
}

func (g *Generator) WithPlugins(f plugins.PluginFeeder) {
	g.pluginsFn = f
}

func (g Generator) Name() string {
	return "resource"
}

func (g Generator) Aliases() []string {
	return []string{"r"}
}

func (g *Generator) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if g.pluginsFn != nil {
		plugs = g.pluginsFn()
	}

	pm := map[string]bool{}

	var builders []plugins.Plugin

	for _, p := range plugs {
		_, ok := p.(BeforeGenerator)
		if !ok {
			continue
		}
		builders = append(builders, p)
	}

	for _, p := range plugs {
		switch p.(type) {
		case ResourceGenerator:
			if pm["ResourceGenerator"] {
				continue
			}
			pm["ResourceGenerator"] = true
			break
		case Actioner:
			if pm["Actioner"] {
				continue
			}
			pm["Actioner"] = true
		case ActionTester:
			if pm["ActionTester"] {
				continue
			}
			pm["ActionTester"] = true
		case Modeler:
			if pm["Modeler"] {
				continue
			}
			pm["Modeler"] = true
		case ModelTester:
			if pm["ModelTester"] {
				continue
			}
			pm["ModelTester"] = true
		case Templater:
			if pm["Templater"] {
				continue
			}
			pm["Templater"] = true
		case TemplateTester:
			if pm["TemplateTester"] {
				continue
			}
			pm["TemplateTester"] = true
		default:
			continue
		}

		builders = append(builders, p)
	}

	for _, p := range plugs {
		_, ok := p.(AfterGenerator)
		if !ok {
			continue
		}
		builders = append(builders, p)
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
			fmt.Printf("[Resource] BeforeGenerator %s\n", p.Name())
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
			fmt.Printf("[Resource] AfterGenerator %s\n", p.Name())
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

	if g.help {
		ioe := plugins.CtxIO(ctx)
		return plugprint.Print(ioe.Stdout(), g)
	}

	plugs := g.ScopedPlugins()

	if err := g.beforeGenerate(ctx, root, args); err != nil {
		return g.afterGenerate(ctx, root, args, err)
	}

	for _, p := range plugs {
		gr, ok := p.(ResourceGenerator)
		if !ok {
			continue
		}
		err := safe.RunE(func() error {
			fmt.Printf("[Resource] ResourceGenerator %s\n", p.Name())
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
			fmt.Printf("[Resource] Actioner %s\n", p.Name())
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
			fmt.Printf("[Resource] ActionTester %s\n", p.Name())
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
			fmt.Printf("[Resource] Templater %s\n", p.Name())
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
			fmt.Printf("[Resource] TemplateTester %s\n", p.Name())
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
			fmt.Printf("[Resource] Modeler %s\n", p.Name())
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
			fmt.Printf("[Resource] ModelTester %s\n", p.Name())
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
			fmt.Printf("[Resource] Migrationer %s\n", p.Name())
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
			fmt.Printf("[Resource] MigrationTester %s\n", p.Name())
			return ag.GenerateResourceMigrationTests(ctx, root, args)
		})
	}
	return nil
}
