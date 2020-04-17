package buildtest

import "flag"

type Flagger []*flag.Flag

func (Flagger) PluginName() string {
	return "buildtest/flagger"
}

func (f Flagger) BuildFlags() []*flag.Flag {
	return []*flag.Flag(f)
}
