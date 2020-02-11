package models

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/generate"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/resource"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"
	gmodel "github.com/gobuffalo/pop/v5/genny/model"
	"github.com/gobuffalo/pop/v5/soda/cmd"
	"github.com/spf13/pflag"
)

var _ plugcmd.Namer = Generator{}
var _ generate.Generator = &Generator{}
var _ plugins.Plugin = Generator{}
var _ plugprint.Describer = Generator{}
var _ plugprint.FlagPrinter = &Generator{}
var _ resource.Modeler = &Generator{}
var _ resource.Pflagger = &Generator{}

type Generator struct {
	ModelPath string
	ModelPkg  string
	StructTag string

	flags *pflag.FlagSet
}

func (Generator) PluginName() string {
	return "pop/model"
}

func (Generator) CmdName() string {
	return "model"
}

func (Generator) Description() string {
	return "Generate a Pop model"
}

// Generate implements generate.Generator and is the entry point for `buffalo generate model`
func (mg *Generator) Generate(ctx context.Context, root string, args []string) error {
	args = append([]string{"generate", "model"}, args...)
	cmd.RootCmd.SetArgs(args)
	return cmd.RootCmd.Execute()
}

// GenerateResourceModels implements resource.Modeler and is responsible for generating a model
// during `buffalo generate resource`
func (mg *Generator) GenerateResourceModels(ctx context.Context, root string, args []string) error {
	flags := mg.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}
	args = flags.Args()

	run := genny.WetRunner(context.Background())

	modelPath := mg.ModelPath
	if len(modelPath) == 0 {
		modelPath = "models"
	}
	modelPath = filepath.Join(root, modelPath)

	structTag := mg.StructTag
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
	return mg.Generate(ctx, root, args)
}
