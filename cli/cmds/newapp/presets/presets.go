package presets

func Presets() map[string]string {
	return map[string]string{
		"core": "github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp/presets/coreapp",
		"json": "github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp/presets/jsonapp",
		"web":  "github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp/presets/webapp",
	}
}
