package flagger

import (
	"fmt"
	"path"

	"github.com/gobuffalo/plugins"
	"github.com/spf13/pflag"
)

func CleanPflags(p plugins.Plugin, pflags []*pflag.Flag) []*pflag.Flag {
	flags := make([]*pflag.Flag, len(pflags))
	for i, f := range pflags {
		flags[i] = &pflag.Flag{
			Name:        fmt.Sprintf("%s-%s", path.Base(name(p)), f.Name),
			Usage:       fmt.Sprintf("[%s] %s", p.PluginName(), f.Usage),
			Value:       f.Value,
			DefValue:    f.DefValue,
			NoOptDefVal: f.NoOptDefVal,
		}
	}
	return flags
}

// SetToSlice takes a flag set and returns a slice
// of the flags
func SetToSlice(set *pflag.FlagSet) []*pflag.Flag {
	var flags []*pflag.Flag
	set.VisitAll(func(f *pflag.Flag) {
		flags = append(flags, f)
	})
	return flags
}

func name(p plugins.Plugin) string {
	type cmdName interface {
		CmdName() string
	}
	if c, ok := p.(cmdName); ok {
		return c.CmdName()
	}
	return p.PluginName()
}
