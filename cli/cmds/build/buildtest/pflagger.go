package buildtest

import "github.com/spf13/pflag"

type Pflagger []*pflag.Flag

func (Pflagger) PluginName() string {
	return "buildtest/flagger"
}

func (f Pflagger) BuildFlags() []*pflag.Flag {
	return []*pflag.Flag(f)
}
