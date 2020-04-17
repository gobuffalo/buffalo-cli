package newapp

import "github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/golang/mod"

var _ mod.Replacer = devReplacer(nil)

type devReplacer func(root string) map[string]string

func (devReplacer) PluginName() string {
	return "dev-replacer"
}

func (d devReplacer) ModReplace(root string) map[string]string {
	return d(root)
}
