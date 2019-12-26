package buildcmd

import (
	"context"
	"flag"
	"os/exec"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/spf13/pflag"
)

// Builder is a sub-command of buffalo build.
// 	buffalo build assets
type Builder interface {
	plugins.Plugin
	Build(ctx context.Context, args []string) error
}

type BeforeBuilder interface {
	plugins.Plugin
	BeforeBuild(ctx context.Context, args []string) error
}

type AfterBuilder interface {
	plugins.Plugin
	AfterBuild(ctx context.Context, args []string, err error) error
}

type Flagger interface {
	plugins.Plugin
	BuildFlags() []*flag.Flag
}

type Pflagger interface {
	plugins.Plugin
	BuildFlags() []*pflag.Flag
}

type TemplatesValidator interface {
	plugins.Plugin
	ValidateTemplates(root string) error
}

type Packager interface {
	plugins.Plugin
	Package(ctx context.Context, root string, files []string) error
}

type PackFiler interface {
	plugins.Plugin
	PackageFiles(ctx context.Context, root string) ([]string, error)
}

type Versioner interface {
	plugins.Plugin
	BuildVersion(ctx context.Context, root string) (string, error)
}

type Importer interface {
	plugins.Plugin
	BuildImports(ctx context.Context, root string) ([]string, error)
}

// BuilderContext can be implemented to capture the `go build` command
// before it is executed. It is up to the BuilderContext to execute, or not,
// the command.
type BuilderContext interface {
	context.Context
	Build(cmd *exec.Cmd) error
}

type buildContext struct {
	context.Context
	fn func(cmd *exec.Cmd) error
}

func (c *buildContext) Build(cmd *exec.Cmd) error {
	if c.fn == nil {
		return nil
	}
	return c.fn(cmd)
}

func WithBuilderContext(ctx context.Context, fn func(cmd *exec.Cmd) error) BuilderContext {
	return &buildContext{
		Context: ctx,
		fn:      fn,
	}
}
