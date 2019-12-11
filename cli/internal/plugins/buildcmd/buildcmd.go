package buildcmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo-cli/internal/v1/genny/build"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here/there"
	"github.com/gobuffalo/meta"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = &BuildCmd{}
var _ plugprint.Aliases = &BuildCmd{}
var _ plugprint.SubCommander = &BuildCmd{}
var _ plugprint.Describer = &BuildCmd{}
var _ plugprint.FlagPrinter = &BuildCmd{}
var _ plugprint.WithPlugins = &BuildCmd{}

type BuildCmd struct {
	plugins.IO
	Parent                 plugins.Plugin
	Plugins                func() plugins.Plugins
	dryRun                 bool
	help                   bool
	skipAssets             bool
	skipTemplateValidation bool
	verbose                bool
	tags                   string
}

func (*BuildCmd) Aliases() []string {
	return []string{"b", "install"}
}

func (b BuildCmd) Name() string {
	return "build"
}

func (b BuildCmd) String() string {
	s := b.Name()
	if b.Parent != nil {
		s = fmt.Sprintf("%s %s", b.Parent.Name(), b.Name())
	}
	return strings.TrimSpace(s)
}

func (BuildCmd) Description() string {
	return "Build the application binary, including bundling of assets (packr & webpack)"
}

func (bc *BuildCmd) SubCommands() plugins.Plugins {
	var plugs plugins.Plugins
	for _, p := range bc.WithPlugins() {
		if _, ok := p.(Builder); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}

func (bc *BuildCmd) WithPlugins() plugins.Plugins {
	var plugs plugins.Plugins
	if bc.Plugins != nil {
		plugs = bc.Plugins()
	}

	var builders plugins.Plugins
	for _, p := range plugs {
		switch p.(type) {
		case Builder:
			builders = append(builders, p)
		case BeforeBuilder:
			builders = append(builders, p)
		case AfterBuilder:
			builders = append(builders, p)
		}
	}
	return builders
}

func (bc *BuildCmd) PrintFlags(w io.Writer) error {
	flags := bc.flagSet(&build.Options{})
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (bc *BuildCmd) flagSet(opts *build.Options) *pflag.FlagSet {
	flags := pflag.NewFlagSet(bc.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)

	flags.BoolVar(&bc.skipTemplateValidation, "skip-template-validation", false, "skip validating templates")
	flags.BoolVarP(&bc.help, "help", "h", false, "print this help")
	flags.BoolVarP(&bc.verbose, "verbose", "v", false, "print debugging information")
	flags.BoolVarP(&opts.Static, "static", "s", false, "build a static binary using  --ldflags '-linkmode external -extldflags \"-static\"'")

	flags.StringVar(&opts.LDFlags, "ldflags", "", "set any ldflags to be passed to the go build")
	flags.StringVar(&opts.Mod, "mod", "", "-mod flag for go build")
	flags.StringVarP(&opts.App.Bin, "output", "o", opts.Bin, "set the name of the binary")
	flags.StringVarP(&opts.Environment, "environment", "", "development", "set the environment for the binary")
	flags.StringVarP(&bc.tags, "tags", "t", "", "compile with specific build tags")

	plugs := bc.WithPlugins()

	for _, p := range plugs {
		switch t := p.(type) {
		case BuildFlagger:
			for _, f := range t.BuildFlags() {
				flags.AddGoFlag(f)
			}
		case BuildPflagger:
			for _, f := range t.BuildFlags() {
				flags.AddFlag(f)
			}
		}
	}

	return flags
}

func (bc *BuildCmd) Main(ctx context.Context, args []string) error {
	info, err := there.Current()
	if err != nil {
		return err
	}

	opts := &build.Options{
		App: meta.New(info.Root),
	}

	flags := bc.flagSet(opts)
	if err = flags.Parse(args); err != nil {
		return err
	}

	if bc.help {
		return plugprint.Print(bc.Stdout(), bc)
	}

	plugs := bc.WithPlugins()

	if len(flags.Args()) > 0 {
		for _, p := range plugs {
			b, ok := p.(Builder)
			if !ok {
				continue
			}
			if p.Name() == args[0] {
				return b.Build(ctx, args)
			}
		}
		return fmt.Errorf("unknown command %q", args[0])
	}

	builders := bc.WithPlugins()
	for _, p := range builders {
		if bb, ok := p.(BeforeBuilder); ok {
			plugins.SetIO(bc, p)
			if err := bb.BeforeBuild(ctx, args); err != nil {
				return err
			}
		}
	}

	if bc.verbose {
		opts.BuildFlags = append(opts.BuildFlags, "-v")
	}

	if len(bc.tags) > 0 {
		opts.Tags = append(opts.Tags, bc.tags)
	}

	if !bc.skipTemplateValidation {
		for _, p := range plugs {
			tv, ok := p.(TemplatesValidator)
			if !ok {
				continue
			}
			plugins.SetIO(bc, p)
			if err := tv.ValidateTemplates(filepath.Join(info.Root, "templates")); err != nil {
				return err
			}
		}
	}

	for _, p := range plugs {
		pkg, ok := p.(Packager)
		if !ok {
			continue
		}
		plugins.SetIO(bc, p)
		if err := pkg.Package(ctx, info.Root); err != nil {
			return err
		}
	}

	version := time.Now().Format(time.RFC3339)
	for _, p := range plugs {
		bv, ok := p.(BuildVersioner)
		if !ok {
			continue
		}
		plugins.SetIO(bc, p)
		s, err := bv.BuildVersion(ctx, info.Root)
		if err != nil {
			return err
		}
		if len(s) == 0 {
			continue
		}
		version = strings.TrimSpace(s)
	}

	fmt.Println("TODO: go build ...", version)
	// opts.GoCommand = bc.Name()
	// clean := build.Cleanup(opts)
	// defer func() {
	// 	if err := clean(run); err != nil {
	// 		log.Fatal("build:clean", err)
	// 	}
	// }()
	//
	// bd, err := build.New(opts)
	// if err != nil {
	// 	return err
	// }
	//
	// // opts.BuildVersion = cmd.buildVersion(opts)
	// // fmt.Println(">>>TODO cli/build.go:106: opts ", opts)
	//
	// if err := run.With(bd); err != nil {
	// 	return err
	// }
	// if err := run.Run(); err != nil {
	// 	return err
	// }

	for _, p := range builders {
		if bb, ok := p.(AfterBuilder); ok {
			plugins.SetIO(bc, p)
			if err := bb.AfterBuild(ctx, args); err != nil {
				return err
			}
		}
	}
	return nil
}

// func (bc *BuildCmd) buildVersion(opts *build.Options) string {
// 	version := opts.BuildTime.Format(time.RFC3339)
// 	vcs := opts.VCS
//
// 	if len(vcs) == 0 {
// 		return version
// 	}
//
// 	if _, err := exec.LookPath(vcs); err != nil {
// 		return version
// 	}
//
// 	var cmd *exec.Cmd
// 	switch vcs {
// 	case "git":
// 		cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
// 	case "bzr":
// 		cmd = exec.Command("bzr", "revno")
// 	default:
// 		return vcs
// 	}
//
// 	out := &bytes.Buffer{}
// 	cmd.Stdout = out
// 	if err := cmd.Run(); err != nil {
// 		return version
// 	}
//
// 	if out.String() == "" {
// 		return version
// 	}
//
// 	return strings.TrimSpace(out.String())
// }
