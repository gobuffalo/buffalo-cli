package mainfix

import (
	"context"
	"os"
	"path"
	"path/filepath"

	"github.com/gobuffalo/here"
)

// Cmd takes care of moving existing main.go into
// cmd/[name]/main.go where its expected by the build command.
type Cmd struct{}

func (*Cmd) PluginName() string {
	return "main/fixer"
}

func (fixer *Cmd) Fix(ctx context.Context, root string, args []string) error {
	info, err := here.Dir(root)
	if err != nil {
		return err
	}

	fp := filepath.Join(root, "cmd", path.Base(info.Module.Path))
	if _, err := os.Stat(filepath.Join(fp, "main.go")); err == nil {
		return nil
	}

	if _, err := os.Stat(fp); err != nil {
		err := os.MkdirAll(fp, 0777)
		if err != nil {
			return err
		}
	}

	mainPath := filepath.Join(root, "main.go")
	if _, err := os.Stat(mainPath); err == nil {
		err := os.Rename(mainPath, filepath.Join(fp, "main.go"))
		if err != nil {
			return err
		}
	}

	return nil
}
