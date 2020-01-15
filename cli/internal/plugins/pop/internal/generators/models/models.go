package models

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/generatecmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/soda"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/genny/v2"
	gmodel "github.com/gobuffalo/pop/v5/genny/model"
)

type Generator struct {
	modelPath string
	modelPkg  string
	structTag string
}

var _ plugins.Plugin = Generator{}

func (Generator) Name() string {
	return "pop/model"
}

var _ plugprint.NamedCommand = Generator{}

func (Generator) CmdName() string {
	return "model"
}

var _ plugprint.Describer = Generator{}

func (Generator) Description() string {
	return "Generate a Pop model"
}

var _ generatecmd.Generator = &Generator{}

// Generate implements generatecmd.Generator and is the entry point for `buffalo generate model`
func (mg *Generator) Generate(ctx context.Context, args []string) error {
	args = append([]string{"generate", "model"}, args...)
	return soda.Main(ctx, args)
}

var _ resource.Modeler = &Generator{}

// GenerateResourceModels implements resource.Modeler and is responsible for generating a model
// during `buffalo generate resource`
func (mg *Generator) GenerateResourceModels(ctx context.Context, root string, args []string) error {
	flags := mg.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}
	args = flags.Args()

	run := genny.WetRunner(context.Background())

	modelPath := mg.modelPath
	if len(modelPath) == 0 {
		modelPath = "models"
	}
	modelPath = filepath.Join(root, modelPath)

	structTag := mg.structTag
	if len(structTag) == 0 {
		structTag = "json"
	}

	atts, err := attrs.ParseArgs(args[1:]...)
	if err != nil {
		return err
	}

	// Mount models generator
	g, err := gmodel.New(&gmodel.Options{
		Name:                   args[0],
		Attrs:                  atts,
		Path:                   modelPath,
		Encoding:               structTag,
		ForceDefaultID:         true,
		ForceDefaultTimestamps: true,
	})
	if err != nil {
		return err
	}

	run.With(g)
	return mg.Generate(ctx, args)
}
