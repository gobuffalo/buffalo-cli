package cmd

import (
	"encoding/json"
	"os"

	bufcli "github.com/gobuffalo/buffalo-cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var jsonOutput bool

func init() {
	decorate("version", versionCmd)
	versionCmd.Flags().BoolVar(&jsonOutput, "json", false, "Print information in json format")

	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Run: func(c *cobra.Command, args []string) {
		if jsonOutput {
			enc := json.NewEncoder(os.Stderr)
			enc.SetIndent("", "    ")
			enc.Encode(map[string]string{
				"version": bufcli.Version,
			})
			return
		}

		logrus.Infof("Buffalo version is: %s", bufcli.Version)
	},
	// needed to override the root level pre-run func
	PersistentPreRunE: func(c *cobra.Command, args []string) error {
		return nil
	},
}
