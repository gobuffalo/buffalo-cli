package build

import (
	"flag"
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build/buildtest"
	"github.com/gobuffalo/plugins"
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

	sflags := []*flag.Flag{
		{
			Name:     "my-flag",
			DefValue: "unset",
			Value:    flagValue(""),
		},
	}
	plugs = append(plugs, buildtest.Flagger(sflags))

	bc = &Cmd{}
	bc.WithPlugins(func() []plugins.Plugin {
		return plugs
	})
	flags = bc.Flags()

	values = []*pflag.Flag{}
	flags.VisitAll(func(f *pflag.Flag) {
		values = append(values, f)
	})
	r.Equal(count+1, len(values))

	count = len(values)

	pflags := []*pflag.Flag{
		{
			Name:     "your-flag",
			DefValue: "unset",
			Value:    flagValue(""),
		},
	}
	plugs = append(plugs, buildtest.Pflagger(pflags))

	bc = &Cmd{}
	bc.WithPlugins(func() []plugins.Plugin {
		return plugs
	})
	flags = bc.Flags()

	values = []*pflag.Flag{}
	flags.VisitAll(func(f *pflag.Flag) {
		values = append(values, f)
	})
	r.Equal(count+1, len(values))
}
