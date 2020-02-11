package packr

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build"
	"github.com/gobuffalo/packr/v2/jam"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
)

var _ build.BeforeBuilder = &Packager{}
var _ build.Packager = &Packager{}
var _ plugcmd.Namer = &Packager{}
var _ plugins.Plugin = &Packager{}

type Packager struct{}

func (b *Packager) BeforeBuild(ctx context.Context, root string, args []string) error {
	return jam.Clean()
}

func (b *Packager) Package(ctx context.Context, root string, files []string) error {
	if len(files) > 0 {
		fmt.Printf("%s does not support additional files\n", b.PluginName())
		for _, f := range files {
			fmt.Printf("\t> %s\n", f)
		}
	}
	return jam.Pack(jam.PackOptions{
		Roots: []string{root},
	})
}

func (b Packager) PluginName() string {
	return "packr"
}

func (b Packager) CmdName() string {
	return "packr"
}
