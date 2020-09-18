package refresh

import (
	"context"
	"path"
	"path/filepath"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/markbates/refresh/refresh"
)

// Fixer will fix current .buffalo.dev.yml to build cmd/[name]/ instead of
// expecting the main to be in the root folder.
type Fixer struct{}

func (*Fixer) PluginName() string {
	return "refresh/fixer"
}

// Fix changes .buffalo.dev.yml to build cmd/[name] if it exists, otherwise
// it creates the file in the root of the project.
func (f *Fixer) Fix(ctx context.Context, root string, args []string) error {

	info, err := here.Dir(root)
	if err != nil {
		return plugins.Wrap(f, err)
	}

	buildRoot := filepath.Join(root, "cmd", path.Base(info.Module.Path))

	c := &refresh.Configuration{
		BuildTargetPath:    buildRoot,
		IgnoredFolders:     []string{"vendor", "log", "logs", "webpack", "public", "grifts", "tmp", "bin", "node_modules", ".sass-cache"},
		IncludedExtensions: []string{".go", ".mod", ".env"},
		BuildPath:          "tmp",
		BuildDelay:         200,
		BinaryName:         "",
		EnableColors:       true,
		LogName:            "buffalo",
	}

	configPath := filepath.Join(root, ".buffalo.dev.yml")
	c.Load(configPath)

	relative, err := filepath.Rel(root, buildRoot)
	if err != nil {
		return plugins.Wrap(f, err)
	}

	c.BuildTargetPath = "." + string(filepath.Separator) + relative
	return c.Dump(configPath)
}
