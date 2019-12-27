package versioncmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	bufcli "github.com/gobuffalo/buffalo-cli"
	"github.com/gobuffalo/buffalo-cli/internal/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/plugins/plugprint"
	"github.com/spf13/pflag"
)

// VersionCmd is responsible for the `buffalo version` command.
type VersionCmd struct {
	help bool
	json bool
}

var _ plugprint.FlagPrinter = &VersionCmd{}

func (vc *VersionCmd) PrintFlags(w io.Writer) error {
	flags := vc.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

var _ plugins.Plugin = &VersionCmd{}

func (vc *VersionCmd) Name() string {
	return "version"
}

var _ plugprint.Describer = &VersionCmd{}

func (vc *VersionCmd) Description() string {
	return "Print the version information"
}

func (i VersionCmd) String() string {
	return i.Name()
}

func (vc *VersionCmd) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(vc.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.BoolVarP(&vc.help, "help", "h", false, "print this help")
	flags.BoolVarP(&vc.json, "json", "j", false, "Print information in json format")
	return flags
}

func (vc *VersionCmd) Main(ctx context.Context, args []string) error {
	flags := vc.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

	ioe := plugins.CtxIO(ctx)
	out := ioe.Stdout()
	if vc.help {
		return plugprint.Print(out, vc)
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
