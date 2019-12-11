package versioncmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	bufcli "github.com/gobuffalo/buffalo-cli"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = &VersionCmd{}
var _ plugprint.Describer = &VersionCmd{}
var _ plugprint.FlagPrinter = &VersionCmd{}

type VersionCmd struct {
	plugins.IO
	Parent plugins.Plugin
	help   bool
	json   bool
}

func (vc *VersionCmd) PrintFlags(w io.Writer) error {
	flags := vc.flagSet()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
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

func (vc *VersionCmd) flagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet(vc.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.BoolVarP(&vc.help, "help", "h", false, "print this help")
	flags.BoolVarP(&vc.json, "json", "j", false, "Print information in json format")
	return flags
}

func (vc *VersionCmd) Main(ctx context.Context, args []string) error {
	flags := vc.flagSet()
	if err := flags.Parse(args); err != nil {
		return err
	}

	out := vc.Stdout()
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
