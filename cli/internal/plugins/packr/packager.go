package packr

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/build"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/packr/v2/jam"
)

var _ build.BeforeBuilder = &Packager{}
var _ build.Packager = &Packager{}
var _ plugins.Plugin = Packager{}
var _ plugprint.NamedCommand

type Packager struct{}

func (b *Packager) BeforeBuild(ctx context.Context, args []string) error {
	return jam.Clean()
}

func (b *Packager) Package(ctx context.Context, root string, files []string) error {
	if len(files) > 0 {
		fmt.Printf("%s does not support additional files", b.Name())
	}
	return jam.Pack(jam.PackOptions{
		Roots: []string{root},
	})
}

func (b Packager) Name() string {
	return "packr"
}

func (b Packager) CmdName() string {
	return "packr"
}
