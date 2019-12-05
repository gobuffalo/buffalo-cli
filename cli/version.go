package cli

import (
	"context"
	"encoding/json"
	"fmt"

	bufcli "github.com/gobuffalo/buffalo-cli"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
)

func (b *Buffalo) Version(ctx context.Context, args []string) error {
	var help bool
	var jsonOutput bool
	flags := cmdx.NewFlagSet("buffalo info", cmdx.Stderr(ctx))
	flags.BoolVarP(&help, "help", "h", false, "print this help")
	flags.BoolVarP(&jsonOutput, "json", "j", false, "Print information in json format")
	if err := flags.Parse(args); err != nil {
		return err
	}

	if help {
		flags.Usage()
		return nil
	}

	if !jsonOutput {
		fmt.Fprintln(cmdx.Stdout(ctx), bufcli.Version)
		return nil
	}

	enc := json.NewEncoder(cmdx.Stdout(ctx))
	enc.SetIndent("", "    ")
	return enc.Encode(map[string]string{
		"version": bufcli.Version,
	})

}
