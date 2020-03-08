package buildtest

type BuildArger func(args []string) []string

func (BuildArger) PluginName() string {
	return "buildtest/build-arger"
}

func (b BuildArger) GoBuildArgs(args []string) []string {
	return b(args)
}
