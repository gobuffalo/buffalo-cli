package build

import (
	"flag"
	"testing"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

type flagValue string

func (f flagValue) String() string {
	return string(f)
}

func (f flagValue) Type() string {
	return string(f)
}

func (f flagValue) Set(value string) error {
	return nil
}

func Test_Cmd_Flags(t *testing.T) {
	r := require.New(t)

	var plugs plugins.Plugins

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	flags := bc.Flags()

	var values []*pflag.Flag
	flags.VisitAll(func(f *pflag.Flag) {
		values = append(values, f)
	})

	count := len(values)

	r.True(count > 0)

	plugs = append(plugs, &buildFlagger{
		flags: []*flag.Flag{
			{
				Name:     "my-flag",
				DefValue: "unset",
				Value:    flagValue(""),
			},
		},
	})

	bc = &Cmd{}
	bc.WithPlugins(plugs.ScopedPlugins)
	flags = bc.Flags()

	values = []*pflag.Flag{}
	flags.VisitAll(func(f *pflag.Flag) {
		values = append(values, f)
	})
	r.Equal(count+1, len(values))

	count = len(values)

	plugs = append(plugs, &buildPflagger{
		flags: []*pflag.Flag{
			{
				Name:     "your-flag",
				DefValue: "unset",
				Value:    flagValue(""),
			},
		},
	})

	bc = &Cmd{}
	bc.WithPlugins(plugs.ScopedPlugins)
	flags = bc.Flags()

	values = []*pflag.Flag{}
	flags.VisitAll(func(f *pflag.Flag) {
		values = append(values, f)
	})
	r.Equal(count+1, len(values))
}
