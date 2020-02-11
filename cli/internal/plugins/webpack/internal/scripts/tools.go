package scripts

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack/internal/ifaces"
	"github.com/gobuffalo/plugins"
)

// Tool tries to determine whether to use yarn or npm
func Tool(plug plugins.Plugin, ctx context.Context, root string) (string, error) {
	if pp, ok := plug.(plugins.Scoper); ok {
		for _, p := range pp.ScopedPlugins() {
			if tp, ok := p.(ifaces.Tooler); ok {
				return tp.AssetTool(ctx, root)
			}
		}
	}

	if _, err := os.Stat(filepath.Join(root, "yarn.lock")); err == nil {
		return "yarnpkg", nil
	}

	if _, err := os.Stat(filepath.Join(root, "package.json")); err == nil {
		return "npm", nil
	}

	return "", fmt.Errorf("could not determine asset tool from %q", root)
}
