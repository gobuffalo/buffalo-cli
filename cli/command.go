package cli

import (
	"fmt"
)

// Commands is a slice of type `Command`
type Commands []Command

// Find will try and find the given command in the slice
// by it's `Aliases()` or `Name()` method.
// If it can't be found an error is returned.
func (commands Commands) Find(name string) (Command, error) {
	for _, c := range commands {
		names := []string{c.Name()}
		if a, ok := c.(NamedCommand); ok {
			names = append(names, a.CmdName())
		}
		if a, ok := c.(Aliases); ok {
			names = append(names, a.Aliases()...)
		}
		for _, n := range names {
			if n == name {
				return c, nil
			}
		}
	}
	return nil, fmt.Errorf("command %s not found", name)
}
