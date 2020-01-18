package flagger

import (
	"flag"
	"fmt"
	"path"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/spf13/pflag"
)

func CleanPflags(p plugins.Plugin, flags []*pflag.Flag) []*pflag.Flag {
	fls := make([]*pflag.Flag, len(flags))
	for i, f := range flags {
		fls[i] = &pflag.Flag{
			Annotations:         f.Annotations,
			Changed:             f.Changed,
			DefValue:            f.DefValue,
			Deprecated:          f.Deprecated,
			Hidden:              f.Hidden,
			Name:                fmt.Sprintf("%s-%s", path.Base(name(p)), f.Name),
			NoOptDefVal:         f.NoOptDefVal,
			ShorthandDeprecated: f.ShorthandDeprecated,
			Usage:               f.Usage,
			Value:               f.Value,
		}
	}
	return fls
}

type cmdName interface {
	CmdName() string
}

func name(p plugins.Plugin) string {
	if c, ok := p.(cmdName); ok {
		return c.CmdName()
	}
	return p.Name()
}

func CleanFlags(p plugins.Plugin, flags []*flag.Flag) []*flag.Flag {
	fls := make([]*flag.Flag, len(flags))
	for i, f := range flags {
		fls[i] = &flag.Flag{
			DefValue: f.DefValue,
			Name:     fmt.Sprintf("%s-%s", path.Base(name(p)), f.Name),
			Usage:    f.Usage,
			Value:    f.Value,
		}
	}
	return fls
}
