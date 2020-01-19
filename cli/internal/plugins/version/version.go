package version

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	bufcli "github.com/gobuffalo/buffalo-cli"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/spf13/pflag"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Cmd{},
	}
}

// Cmd is responsible for the `buffalo version` command.
type Cmd struct {
	help bool
	json bool
}

var _ plugprint.FlagPrinter = &Cmd{}

func (vc *Cmd) PrintFlags(w io.Writer) error {
	flags := vc.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

var _ plugins.Plugin = &Cmd{}

func (vc *Cmd) Name() string {
	return "version"
}

var _ plugprint.Describer = &Cmd{}

func (vc *Cmd) Description() string {
	return "Print the version information"
}

func (i Cmd) String() string {
	return i.Name()
}

func (vc *Cmd) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(vc.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.BoolVarP(&vc.help, "help", "h", false, "print this help")
	flags.BoolVarP(&vc.json, "json", "j", false, "Print information in json format")
	return flags
}

func (vc *Cmd) Main(ctx context.Context, args []string) error {
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
