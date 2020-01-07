package actiongen

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
)

type Generator struct{}

var _ plugins.Plugin = Generator{}

func (Generator) Name() string {
	return "pop/action"
}

var _ plugprint.NamedCommand = Generator{}

func (Generator) CmdName() string {
	return "action"
}

var _ plugprint.Describer = Generator{}

func (Generator) Description() string {
	return "Generate a Pop action"
}

// var _ generatecmd.Generator = &Generator{}
//
// func (mg *Generator) Generate(ctx context.Context, args []string) error {
// 	args = append([]string{"generate", "action"}, args...)
// 	return nil
// }

var _ resource.Actioner = &Generator{}

func (mg *Generator) GenerateResourceActions(ctx context.Context, root string, args []string) error {
	fmt.Println(">>>TODO cli/internal/plugins/pop/internal/actiongen/actiongen.go:43: root ", root)
	fmt.Println(">>>TODO cli/internal/plugins/pop/internal/actiongen/actiongen.go:43: args ", args)
	return nil
}
