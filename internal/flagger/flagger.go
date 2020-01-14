package flagger

import (
	"flag"
	"fmt"
	"path"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/spf13/pflag"
)

func CleanPflags(p plugins.Plugin, flags []*pflag.Flag) []*pflag.Flag {
	for i, f := range flags {
		f.Name = fmt.Sprintf("%s-%s", path.Base(p.Name()), f.Name)
		f.Shorthand = ""
		flags[i] = f
	}
	return flags
}

func CleanFlags(p plugins.Plugin, flags []*flag.Flag) []*flag.Flag {
	for i, f := range flags {
		f.Name = fmt.Sprintf("%s-%s", path.Base(p.Name()), f.Name)
		flags[i] = f
	}
	return flags
}
