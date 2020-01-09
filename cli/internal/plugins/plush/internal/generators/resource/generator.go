package resource

import (
	"context"
	"fmt"
	"os"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/markbates/pkger"
)

type Generator struct {
}

func (g *Generator) Name() string {
	return "plush/templates"
}

var _ resource.Templater = &Generator{}

// Attrs
// Model
// Name
// Folder
func (g *Generator) GenerateResourceTemplates(ctx context.Context, root string, args []string) error {
	flags := g.Flags()

	if err := flags.Parse(args); err != nil {
		return err
	}

	args = flags.Args()

	fp := pkger.Include("github.com/gobuffalo/buffalo-cli:/cli/internal/plugins/plush/internal/generators/resource/templates")

	err := pkger.Walk(fp, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fmt.Println(">>>TODO cli/internal/plugins/plush/internal/generators/resource/generator.go:41: path ", path)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
