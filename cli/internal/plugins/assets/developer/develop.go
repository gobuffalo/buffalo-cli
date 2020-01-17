package developer

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/assets/scripts"
	"github.com/gobuffalo/buffalo-cli/plugins"
)

type Developer struct {
}

func (d *Developer) Name() string {
	return "assets/develop"
}

func (d *Developer) CmdName() string {
	return "assets"
}

func (d *Developer) Develop(ctx context.Context, root string, args []string) error {
	tool, err := scripts.Tool(root)
	if err != nil {
		return err
	}

	// TODO: move to a "before" interface
	// // make sure that the node_modules folder is properly "installed"
	// if _, err := os.Stat(filepath.Join(root, "node_modules")); err != nil {
	// 	cmd := plugins.Cmd(ctx, tool, "install")
	// 	if err := cmd.Run(); err != nil {
	// 		return err
	// 	}
	// }

	cmd := plugins.Cmd(ctx, scripts.WebpackBin(root), "--watch")

	if _, err := scripts.ScriptFor(root, "dev"); err == nil {
		cmd = plugins.Cmd(ctx, tool, "run", "dev")
	}

	return cmd.Run()
}
