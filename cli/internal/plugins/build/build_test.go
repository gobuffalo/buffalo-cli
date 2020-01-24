package build

import (
	"path/filepath"
	"testing"

	"github.com/gobuffalo/here"
)

func newRef(t *testing.T, root string) here.Info {
	t.Helper()

	info := here.Info{
		Dir:        root,
		ImportPath: "github.com/markbates/coke",
		Name:       "coke",
		Root:       root,
		Module: here.Module{
			Path:  "github.com/markbates/coke",
			Main:  true,
			Dir:   root,
			GoMod: filepath.Join(root, "go.mod"),
		},
	}

	return info
}
