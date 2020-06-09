package newapp

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/golang/mod"
	"github.com/gobuffalo/plugins"
)

func (cmd *Cmd) modInit(ctx context.Context, root string, name string) error {
	mi := &mod.Initer{}
	var miplugs []plugins.Plugin

	if cmd.pluginsFn != nil {
		miplugs = cmd.pluginsFn()
	}

	fp := os.Getenv("BUFFALO_CLI")
	if len(fp) == 0 {
		err := mi.ModInit(ctx, root, name)
		return plugins.Wrap(cmd, err)
	}

	if _, err := os.Stat(fp); err != nil {
		return plugins.Wrap(cmd, err)
	}

	rel, err := filepath.Rel(root, fp)
	if err != nil {
		return plugins.Wrap(cmd, err)
	}
	rel = filepath.Dir(rel)

	fn := func(root string) map[string]string {
		return map[string]string{
			"github.com/gobuffalo/buffalo-cli/v2": filepath.Join(rel, "/buffalo-cli"),
		}
	}

	miplugs = append(miplugs, devReplacer(fn))
	mi.WithPlugins(func() []plugins.Plugin {
		return miplugs
	})

	pwd, err := os.Getwd()
	if err != nil {
		return plugins.Wrap(cmd, err)
	}
	defer os.Chdir(pwd)
	os.Chdir(root)

	err = mi.ModInit(ctx, root, name)
	return plugins.Wrap(cmd, err)
}
