package plugins

import (
	"flag"
	"fmt"
	"path"
)

func CleanFlags(p Plugin, flags []*flag.Flag) []*flag.Flag {
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

func name(p Plugin) string {
	type cmdName interface {
		CmdName() string
	}
	if c, ok := p.(cmdName); ok {
		return c.CmdName()
	}
	return p.PluginName()
}
