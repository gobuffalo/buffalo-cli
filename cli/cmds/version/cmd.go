package version

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = &Cmd{}
var _ plugio.OutNeeder = &Cmd{}
var _ plugprint.Describer = &Cmd{}
var _ plugprint.FlagPrinter = &Cmd{}

// Cmd is responsible for the `buffalo version` command.
type Cmd struct {
	help   bool
	json   bool
	stdout io.Writer
}

func (c *Cmd) SetStdout(w io.Writer) error {
	c.stdout = w
	return nil
}

func (vc *Cmd) PrintFlags(w io.Writer) error {
	flags := vc.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (vc *Cmd) PluginName() string {
	return "version"
}

func (vc *Cmd) Description() string {
	return "Print the version information"
}

func (i Cmd) String() string {
	return i.PluginName()
}

func (vc *Cmd) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(vc.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.BoolVarP(&vc.help, "help", "h", false, "print this help")
	flags.BoolVarP(&vc.json, "json", "j", false, "Print information in json format")
	return flags
}

func (vc *Cmd) Main(ctx context.Context, root string, args []string) error {
	flags := vc.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

	out := vc.stdout
	if vc.stdout == nil {
		out = os.Stdout
	}
	if vc.help {
		return plugprint.Print(out, vc)
	}

	if !vc.json {
		fmt.Fprintln(out, Version)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	return enc.Encode(map[string]string{
		"version": Version,
	})

}
