package packr

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/build"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/packr/v2/jam"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Buffalo{},
	}
}

type Buffalo struct{}

var _ build.BeforeBuilder = &Buffalo{}

func (b *Buffalo) BeforeBuild(ctx context.Context, args []string) error {
	return jam.Clean()
}

var _ build.Packager = &Buffalo{}

func (b *Buffalo) Package(ctx context.Context, root string, files []string) error {
	if len(files) > 0 {
		fmt.Printf("%s does not support additional files", b.Name())
	}
	return jam.Pack(jam.PackOptions{
		Roots: []string{root},
	})
}

var _ plugins.Plugin = Buffalo{}

func (b Buffalo) Name() string {
	return "packr"
}

var _ plugprint.NamedCommand

func (b Buffalo) CmdName() string {
	return "buffalo-packr"
}
