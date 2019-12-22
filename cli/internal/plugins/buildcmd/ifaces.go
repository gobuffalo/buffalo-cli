package buildcmd

import (
	"context"
	"flag"

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
