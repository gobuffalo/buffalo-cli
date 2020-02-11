package built

type background string

func (b background) PluginName() string {
	return string(b)
}
