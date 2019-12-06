package plugprint

// Describer is called by `Print` and can be
// implemented to print a short, single line,
// description of the plugin. `-h`
type Describer interface {
	Description() string
}
