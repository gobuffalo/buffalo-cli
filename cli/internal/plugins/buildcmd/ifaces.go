package buildcmd

import (
	"context"
	"flag"

	"github.com/spf13/pflag"
)

type BeforeBuilder interface {
	BeforeBuild(ctx context.Context, args []string) error
}

type AfterBuilder interface {
	AfterBuild(ctx context.Context, args []string) error
}

type BuildFlagger interface {
	BuildFlags() []*flag.Flag
}

type BuildPflagger interface {
	BuildPflags() []*pflag.Flag
}

type TemplatesValidator interface {
	ValidateTemplates(root string) error
}

type Packager interface {
	Package(ctx context.Context, root string) error
}

type BuildVersioner interface {
	BuildVersion(ctx context.Context, root string) (string, error)
}
