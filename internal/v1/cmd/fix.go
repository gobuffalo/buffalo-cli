package cmd

import (
	"fmt"

	bufcli "github.com/gobuffalo/buffalo-cli"
	"github.com/gobuffalo/buffalo-cli/internal/v1/cmd/fix"
	"github.com/spf13/cobra"
)

// fixCmd represents the info command
var fixCmd = &cobra.Command{
	Use:     "fix",
	Aliases: []string{"update"},
	Short:   fmt.Sprintf("Attempt to fix a Buffalo application's API to match version %s", bufcli.Version),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fix.Run()
	},
}

func init() {
	decorate("fix", RootCmd)
	decorate("update", RootCmd)

	fixCmd.Flags().BoolVarP(&fix.YesToAll, "y", "y", false, "update all without asking for confirmation")
	RootCmd.AddCommand(fixCmd)
}
