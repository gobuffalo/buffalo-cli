package resource

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/generate"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/gobuffalo/here"
	"github.com/markbates/safe"
)

var _ generate.Generator = &Generator{}

func (g *Generator) beforeGenerate(ctx context.Context, info here.Info, args []string) error {
	plugs := g.ScopedPlugins()

	for _, p := range plugs {
		b, ok := p.(BeforeGenerator)
		if !ok {
			continue
		}
		err := safe.RunE(func() error {
			fmt.Printf("[Resource] BeforeGenerator %s\n", p.Name())
			return b.BeforeGenerateResource(ctx, info.Dir, args)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) addResource(info here.Info, n string) error {
	fp := filepath.Join(info.Dir, "actions", "app.go")

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
func (g *Generator) afterGenerate(ctx context.Context, info here.Info, args []string, err error) error {
	plugs := g.ScopedPlugins()

	if err == nil && len(args) > 0 {
		if err := g.addResource(info, args[0]); err != nil {
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
			return b.AfterGenerateResource(ctx, info.Dir, args, err)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Generate implements generate.Generator and is the entry point for `buffalo generate resource`
func (g *Generator) Generate(ctx context.Context, args []string) error {
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

	info, err := g.HereInfo()
	if err != nil {
		return err
	}

	plugs := g.ScopedPlugins()

	if err := g.beforeGenerate(ctx, info, args); err != nil {
		return g.afterGenerate(ctx, info, args, err)
	}

	for _, p := range plugs {
		gr, ok := p.(ResourceGenerator)
		if !ok {
			continue
		}
		err := safe.RunE(func() error {
			fmt.Printf("[Resource] ResourceGenerator %s\n", p.Name())
			return gr.GenerateResource(ctx, info.Dir, args)
		})
		return g.afterGenerate(ctx, info, args, err)
	}

	type step func(context.Context, here.Info, []string) error

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
		if err := step(ctx, info, args); err != nil {
			return err
		}
	}

	return g.afterGenerate(ctx, info, args, nil)

}

func (g *Generator) generateActions(ctx context.Context, info here.Info, args []string) error {
	if g.skipActions {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(Actioner)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] Actioner %s\n", p.Name())
			return ag.GenerateResourceActions(ctx, info.Dir, args)
		})
	}
	return nil
}

func (g *Generator) generateActionTests(ctx context.Context, info here.Info, args []string) error {
	if g.skipActionTests {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(ActionTester)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] ActionTester %s\n", p.Name())
			return ag.GenerateResourceActionTests(ctx, info.Dir, args)
		})
	}
	return nil
}

func (g *Generator) generateTemplates(ctx context.Context, info here.Info, args []string) error {
	if g.skipTemplates {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(Templater)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] Templater %s\n", p.Name())
			return ag.GenerateResourceTemplates(ctx, info.Dir, args)
		})
	}
	return nil
}

func (g *Generator) generateTemplateTests(ctx context.Context, info here.Info, args []string) error {
	if g.skipTemplateTests {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(TemplateTester)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] TemplateTester %s\n", p.Name())
			return ag.GenerateResourceTemplateTests(ctx, info.Dir, args)
		})
	}
	return nil
}

func (g *Generator) generateModels(ctx context.Context, info here.Info, args []string) error {
	if g.skipModels {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(Modeler)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] Modeler %s\n", p.Name())
			return ag.GenerateResourceModels(ctx, info.Dir, args)
		})
	}
	return nil
}

func (g *Generator) generateModelTests(ctx context.Context, info here.Info, args []string) error {
	if g.skipModelTests {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(ModelTester)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] ModelTester %s\n", p.Name())
			return ag.GenerateResourceModelTests(ctx, info.Dir, args)
		})
	}
	return nil
}

func (g *Generator) generateMigrations(ctx context.Context, info here.Info, args []string) error {
	if g.skipMigrations {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(Migrationer)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] Migrationer %s\n", p.Name())
			return ag.GenerateResourceMigrations(ctx, info.Dir, args)
		})
	}
	return nil
}

func (g *Generator) generateMigrationTests(ctx context.Context, info here.Info, args []string) error {
	if g.skipMigrationTests {
		return nil
	}
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(MigrationTester)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			fmt.Printf("[Resource] MigrationTester %s\n", p.Name())
			return ag.GenerateResourceMigrationTests(ctx, info.Dir, args)
		})
	}
	return nil
}
