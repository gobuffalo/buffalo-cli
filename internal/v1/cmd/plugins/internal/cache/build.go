package cache

import (
	"os"

	"github.com/gobuffalo/buffalo-cli/internal/v1/plugins"
	"github.com/spf13/cobra"
)

// Cmd rebuilds the plugins cache
var Cmd = &cobra.Command{
	Use:   "build",
	Short: "rebuilds the plugins cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		os.RemoveAll(plugins.CachePath)
		_, err := plugins.Available()
		return err
	},
}
