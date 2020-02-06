package cli

import (
	"fmt"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
)

// Commands is a slice of type `Command`
type Commands []Command

// Find will try and find the given command in the slice
// by it's `Aliases()` or `Name()` method.
// If it can't be found an error is returned.
func (commands Commands) Find(name string) (Command, error) {
	plugs := make([]plugins.Plugin, len(commands))
	for i, c := range commands {
		plugs[i] = c
	}
	p, err := plugins.Commands(plugs).Find(name)
	if err != nil {
		return nil, err
	}
	if c, ok := p.(Command); ok {
		return c, nil
	}
	return nil, fmt.Errorf("command %s not found", name)
}
