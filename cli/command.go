package cli

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/plugins"
)

// Command represents a plugin that can be
// used as a full sub-command. Like Go program's the
// `Main` method is called to run that command.
type Command interface {
	plugins.Plugin
	Main(ctx context.Context, args []string) error
}

// Commands is a slice of type `Command`
type Commands []Command

// Len is the number of elements in the collection.
func (commands Commands) Len() int {
	return len(commands)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (commands Commands) Less(i int, j int) bool {
	return commands[i].Name() < commands[j].Name()
}

// Swap swaps the elements with indexes i and j.
func (commands Commands) Swap(i int, j int) {
	commands[i], commands[j] = commands[j], commands[i]
}

// Find will try and find the given command in the slice
// by it's `Aliases()` or `Name()` method.
// If it can't be found an error is returned.
func (commands Commands) Find(name string) (Command, error) {
	for _, c := range commands {
		names := []string{c.Name()}
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
