package buildcmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/cli/plugins/plugprint"
	"github.com/gobuffalo/buffalo-cli/internal/v1/genny/build"
	"github.com/gobuffalo/here/there"
	"github.com/gobuffalo/meta"
	"github.com/spf13/pflag"
)

func (b *BuildCmd) setIO(p plugins.Plugin) {
	if stdin, ok := p.(plugins.StdinSetter); ok {
		stdin.SetStdin(b.Stdin())
	}
	if stdout, ok := p.(plugins.StdoutSetter); ok {
		stdout.SetStdout(b.Stdout())
	}
	if stderr, ok := p.(plugins.StderrSetter); ok {
		stderr.SetStderr(b.Stderr())
	}
}

type BuildCmd struct {
	Parent                 plugins.Plugin
	Plugins                func() plugins.Plugins
	dryRun                 bool
	help                   bool
	skipAssets             bool
	skipTemplateValidation bool
	verbose                bool
	tags                   string
	stdin                  io.Reader
	stdout                 io.Writer
	stderr                 io.Writer
}

func (b *BuildCmd) Stdin() io.Reader {
	if b.stdin == nil {
		return os.Stdin
	}
	return b.stdin
}

func (b *BuildCmd) Stdout() io.Writer {
	if b.stdout == nil {
		return os.Stdout
	}
	return b.stdout
}

func (b *BuildCmd) Stderr() io.Writer {
	if b.stderr == nil {
		return os.Stderr
	}
	return b.stderr
}

func (b *BuildCmd) SetStderr(w io.Writer) {
	b.stderr = w
}

func (b *BuildCmd) SetStdin(r io.Reader) {
	b.stdin = r
}

func (b *BuildCmd) SetStdout(w io.Writer) {
	b.stdout = w
}

func (*BuildCmd) Aliases() []string {
	return []string{"b", "bill", "install"}
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

func (bc *BuildCmd) plugins() plugins.Plugins {
	if bc.Plugins == nil {
		return nil
	}
	return bc.Plugins()
}

func (bc *BuildCmd) builders() plugins.Plugins {
	var plugs plugins.Plugins
	for _, p := range bc.plugins() {
		switch p.(type) {
		case BeforeBuilder:
			plugs = append(plugs, p)
		case AfterBuilder:
			plugs = append(plugs, p)
		}
	}
	return plugs
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

	plugs := bc.plugins()

	for _, p := range plugs {
		bf, ok := p.(BuildFlagger)
		if !ok {
			continue
		}
		for _, f := range bf.BuildFlags() {
			flags.AddGoFlag(f)
		}
	}

	for _, p := range plugs {
		bf, ok := p.(BuildPflagger)
		if !ok {
			continue
		}
		for _, f := range bf.BuildPflags() {
			flags.AddFlag(f)
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
		return plugprint.Print(bc.stdout, bc, nil)
	}

	plugs := bc.plugins()

	builders := bc.builders()
	for _, p := range builders {
		if bb, ok := p.(BeforeBuilder); ok {
			bc.setIO(p)
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
			bc.setIO(p)
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
		bc.setIO(p)
		if err := pkg.Package(ctx, info.Root); err != nil {
			return err
		}
	}

	version := time.Now().Format(time.RFC3339)
	// for _, p := range plugs {
	// 	bv, ok := p.(BuildVersioner)
	// 	if !ok {
	// 		continue
	// 	}
	// 	bc.setIO(p)
	// 	s, err := bv.BuildVersion(ctx, info.Root)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if len(s) == 0 {
	// 		continue
	// 	}
	// 	version = strings.TrimSpace(s)
	// }

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
			bc.setIO(p)
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
