package versioncmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	bufcli "github.com/gobuffalo/buffalo-cli"
	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
)

type VersionCmd struct {
	Parent plugins.Plugin
	stdout io.Writer
	help   bool
	json   bool
}

func (vc *VersionCmd) SetStdout(w io.Writer) {
	vc.stdout = w
}

func (vc *VersionCmd) Name() string {
	return "version"
}

func (vc *VersionCmd) Description() string {
	return "Print the version information"
}

func (i VersionCmd) String() string {
	s := i.Name()
	if i.Parent != nil {
		s = fmt.Sprintf("%s %s", i.Parent.Name(), i.Name())
	}
	return strings.TrimSpace(s)
}

func (vc *VersionCmd) Main(ctx context.Context, args []string) error {
	flags := cmdx.NewFlagSet(vc.String())
	flags.BoolVarP(&vc.help, "help", "h", false, "print this help")
	flags.BoolVarP(&vc.json, "json", "j", false, "Print information in json format")
	if err := flags.Parse(args); err != nil {
		return err
	}

	out := vc.stdout
	if out == nil {
		out = os.Stdout
	}
	if vc.help {
		return cmdx.Print(out, vc, nil, flags)
	}

	if !vc.json {
		fmt.Fprintln(out, bufcli.Version)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	return enc.Encode(map[string]string{
		"version": bufcli.Version,
	})

}
