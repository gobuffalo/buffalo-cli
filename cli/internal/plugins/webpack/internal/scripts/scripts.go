package scripts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack/internal/ifaces"
	"github.com/gobuffalo/plugins"
)

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

// ScriptFor will attempt to find the named script in the
// package.json file of the application.
func ScriptFor(plug plugins.Plugin, ctx context.Context, root string, name string) (string, error) {

	if pp, ok := plug.(plugins.Scoper); ok {
		for _, p := range pp.ScopedPlugins() {
			if tp, ok := p.(ifaces.Scripter); ok {
				return tp.AssetScript(ctx, root, name)
			}
		}
	}
	scripts := packageJSON{}

	pf, err := os.Open(filepath.Join(root, "package.json"))
	if err != nil {
		return "", err
	}
	defer pf.Close()

	if err = json.NewDecoder(pf).Decode(&scripts); err != nil {
		return "", err
	}

	if s, ok := scripts.Scripts[name]; ok {
		return s, nil
	}
	return "", fmt.Errorf("script %q not found", name)
}

// WebpackBin returns the location of the webpack binary
func WebpackBin(root string) string {
	s := filepath.Join(root, "node_modules", ".bin", "webpack")
	if runtime.GOOS == "windows" {
		s += ".cmd"
	}
	return s
}
