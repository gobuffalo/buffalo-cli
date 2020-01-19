package flagger

import (
	"flag"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/spf13/pflag"
)

func CleanPflags(p plugins.Plugin, pflags []*pflag.Flag) []*flag.Flag {

	flags := make([]*flag.Flag, len(pflags))
	for i, f := range pflags {
		flags[i] = &flag.Flag{
			// DefValue: f.DefValue,
			Name:  f.Name,
			Usage: f.Usage,
			Value: f.Value,
		}
	}
	return plugins.CleanFlags(p, flags)
}
