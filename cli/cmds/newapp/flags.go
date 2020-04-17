package newapp

import (
	"fmt"
	"io"
	"path"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp/presets"
	"github.com/spf13/pflag"
)

func (cmd *Cmd) PrintFlags(w io.Writer) error {
	flags := cmd.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (cmd *Cmd) Flags() *pflag.FlagSet {
	if cmd.flags != nil {
		return cmd.flags
	}
	flags := pflag.NewFlagSet(cmd.PluginName(), pflag.ContinueOnError)
	flags.BoolVarP(&cmd.help, "help", "h", false, "print this help")
	flags.BoolVarP(&cmd.force, "force", "f", false, "delete the existing application first")

	pres := presets.Presets()
	var names []string
	for _, p := range pres {
		names = append(names, path.Base(p))
	}

	flags.StringVarP(&cmd.preset, "preset", "p", "webapp", fmt.Sprintf("preset list of plugins to use %s", names))

	cmd.flags = flags
	return cmd.flags
}
