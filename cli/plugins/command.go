package plugins

import (
	"context"
	"fmt"
)

type Command interface {
	Name() string
	Main(ctx context.Context, args []string) error
}

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

func (plugs Plugins) Commands() Commands {
	var commands Commands
	for _, p := range plugs {
		c, ok := p.(Command)
		if !ok {
			continue
		}
		commands = append(commands, c)
	}
	return commands
}

func (commands Commands) Find(name string) (Command, error) {
	for _, c := range commands {
		if c.Name() == name {
			return c, nil
		}
	}
	return nil, fmt.Errorf("command %s not found", name)
}
