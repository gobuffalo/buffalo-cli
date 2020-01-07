package resource

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/generatecmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
	"github.com/markbates/safe"
)

var _ generatecmd.Generator = &Generator{}

func (g *Generator) beforeGenerate(ctx context.Context, info here.Info, args []string) error {
	plugs := g.ScopedPlugins()

	for _, p := range plugs {
		b, ok := p.(BeforeGenerator)
		if !ok {
			continue
		}
		err := safe.RunE(func() error {
			return b.BeforeGenerateResource(ctx, info.Dir, args)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) afterGenerate(ctx context.Context, info here.Info, args []string, err error) error {
	plugs := g.ScopedPlugins()

	for _, p := range plugs {
		b, ok := p.(AfterGenerator)
		if !ok {
			continue
		}
		err := safe.RunE(func() error {
			return b.AfterGenerateResource(ctx, info.Dir, args, err)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
func (g *Generator) Generate(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must specify a name for the resource")
	}

	flags := g.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

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
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(Actioner)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			return ag.GenerateResourceActions(ctx, info.Dir, args)
		})
	}
	return nil
}

func (g *Generator) generateActionTests(ctx context.Context, info here.Info, args []string) error {
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(ActionTester)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			return ag.GenerateResourceActionTests(ctx, info.Dir, args)
		})
	}
	return nil
}

func (g *Generator) generateTemplates(ctx context.Context, info here.Info, args []string) error {
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(Templater)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			return ag.GenerateResourceTemplates(ctx, info.Dir, args)
		})
	}
	return nil
}

func (g *Generator) generateTemplateTests(ctx context.Context, info here.Info, args []string) error {
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(TemplateTester)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			return ag.GenerateResourceTemplateTests(ctx, info.Dir, args)
		})
	}
	return nil
}

func (g *Generator) generateModels(ctx context.Context, info here.Info, args []string) error {
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(Modeler)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			return ag.GenerateResourceModels(ctx, info.Dir, args)
		})
	}
	return nil
}

func (g *Generator) generateModelTests(ctx context.Context, info here.Info, args []string) error {
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(ModelTester)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			return ag.GenerateResourceModelTests(ctx, info.Dir, args)
		})
	}
	return nil
}

func (g *Generator) generateMigrations(ctx context.Context, info here.Info, args []string) error {
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(Migrationer)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			return ag.GenerateResourceMigrations(ctx, info.Dir, args)
		})
	}
	return nil
}

func (g *Generator) generateMigrationTests(ctx context.Context, info here.Info, args []string) error {
	for _, p := range g.ScopedPlugins() {
		ag, ok := p.(MigrationTester)
		if !ok {
			continue
		}
		return safe.RunE(func() error {
			return ag.GenerateResourceMigrationTests(ctx, info.Dir, args)
		})
	}
	return nil
}
