package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	bufcli "github.com/gobuffalo/buffalo-cli"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
)

type versionCmd struct {
	Buffalo *Buffalo
	help    bool
	json    bool
}

func (vc *versionCmd) Name() string {
	return "version"
}

func (vc *versionCmd) Description() string {
	return "Print the version information"
}

func (vc versionCmd) String() string {
	s := fmt.Sprintf("%s %s", vc.Buffalo, vc.Name())
	return strings.TrimSpace(s)
}

func (vc *versionCmd) Main(ctx context.Context, args []string) error {
	flags := cmdx.NewFlagSet(vc.String())
	flags.BoolVarP(&vc.help, "help", "h", false, "print this help")
	flags.BoolVarP(&vc.json, "json", "j", false, "Print information in json format")
	if err := flags.Parse(args); err != nil {
		return err
	}

	if vc.help {
		return cmdx.Print(vc.Buffalo.Stdout, vc, nil, flags)
	}

	if !vc.json {
		fmt.Fprintln(vc.Buffalo.Stdout, bufcli.Version)
		return nil
	}

	enc := json.NewEncoder(vc.Buffalo.Stdout)
	enc.SetIndent("", "    ")
	return enc.Encode(map[string]string{
		"version": bufcli.Version,
	})

}
