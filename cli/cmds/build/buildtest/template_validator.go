package buildtest

type TemplatesValidator func(root string) error

func (TemplatesValidator) PluginName() string {
	return "buildtest/templatesValidator"
}

func (t TemplatesValidator) ValidateTemplates(root string) error {
	return t(root)
}
