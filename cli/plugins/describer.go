package plugins

type Describer interface {
	Description() string
}

func Description(p Plugin) string {
	if d, ok := p.(Describer); ok {
		return d.Description()
	}
	return ""
}
