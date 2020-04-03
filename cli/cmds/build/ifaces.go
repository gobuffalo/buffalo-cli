package build

import (
	"context"
	"flag"
	"os/exec"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/spf13/pflag"
)

// Builder is a sub-command of buffalo build.
// 	buffalo build webpack
type Builder interface {
	plugins.Plugin
	Build(ctx context.Context, root string, args []string) error
}

type BeforeBuilder interface {
	plugins.Plugin
	BeforeBuild(ctx context.Context, root string, args []string) error
}

type AfterBuilder interface {
	plugins.Plugin
	AfterBuild(ctx context.Context, root string, args []string, err error) error
}

type Flagger interface {
	plugins.Plugin
	BuildFlags() []*flag.Flag
}

type Pflagger interface {
	plugins.Plugin
	BuildFlags() []*pflag.Flag
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

type Runner interface {
	plugins.Plugin
	RunBuild(ctx context.Context, cmd *exec.Cmd) error
}

type Tagger interface {
	plugins.Plugin
	BuildTags(ctx context.Context, root string) ([]string, error)
}

type Namer interface {
	Builder
	CmdName() string
}

type Aliaser interface {
	Builder
	CmdAliases() []string
}

type Stdouter = plugio.Outer
