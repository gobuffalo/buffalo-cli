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
