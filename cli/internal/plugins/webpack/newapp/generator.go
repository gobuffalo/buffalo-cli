package newapp

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/markbates/pkger"
	"github.com/spf13/pflag"
)

var _ newapp.Newapper = &Generator{}

type Generator struct {
	pluginsFn plugins.Feeder
	flags     *pflag.FlagSet
	tool      string
}

func (Generator) PluginName() string {
	return "webpack/newapp"
}

func (Generator) CmdName() string {
	return "webpack"
}

func (g *Generator) WithPlugins(f plugins.Feeder) {
	g.pluginsFn = f
}

func (g *Generator) ScopedPlugins() []plugins.Plugin {
	if g.pluginsFn == nil {
		return nil
	}

	var plugs []plugins.Plugin
	for _, p := range g.pluginsFn() {
		switch p.(type) {
		case Stdouter:
			plugs = append(plugs, p)
		case Stdiner:
			plugs = append(plugs, p)
		case Stderrer:
			plugs = append(plugs, p)
		case CommandRunner:
			plugs = append(plugs, p)
		}
	}

	return plugs
}

func (g *Generator) Newapp(ctx context.Context, root string, name string, args []string) error {
	if err := g.packageJSON(ctx, root, name); err != nil {
		return plugins.Wrap(g, err)
	}

	if err := g.copyTemplates(ctx, root, name); err != nil {
		return plugins.Wrap(g, err)
	}

	return nil
}

func (g *Generator) copyTemplates(ctx context.Context, root string, name string) error {
	troot := pkger.Include("github.com/gobuffalo/buffalo-cli/v2:/cli/internal/plugins/webpack/newapp/templates")

	err := pkger.Walk(troot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		p := strings.TrimPrefix(path, troot)
		p = strings.ReplaceAll(p, "-dot-", ".")
		p = strings.ReplaceAll(p, ".tmpl", "")

		if info.IsDir() {
			if err := os.MkdirAll(filepath.Join(root, p), 0755); err != nil {
				return fmt.Errorf("%s: %w", p, err)
			}
			return nil
		}

		src, err := pkger.Open(path)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		defer src.Close()

		dest, err := os.Create(filepath.Join(root, p))
		if err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		defer src.Close()

		if _, err = io.Copy(dest, src); err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("%s: %w", troot, err)
	}
	return nil
}

func (g *Generator) packageJSON(ctx context.Context, root string, name string) error {
	t, err := template.New("package.json").Parse(tmpl)
	if err != nil {
		return err
	}

	data := map[string]string{
		"Name": name,
	}

	f, err := os.Create(filepath.Join(root, "package.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := t.Execute(f, data); err != nil {
		return err
	}

	tool := g.tool

	c := exec.CommandContext(ctx, tool, "add")
	if err := g.installDeps(ctx, c, root, dependencies); err != nil {
		return err
	}

	c = exec.CommandContext(ctx, tool, "add", "--dev")
	if err := g.installDeps(ctx, c, root, devDependencies); err != nil {
		return err
	}

	return nil
}

func (g *Generator) installDeps(ctx context.Context, c *exec.Cmd, root string, deps map[string]string) error {
	plugs := g.ScopedPlugins()

	c.Stdout = plugio.Stdout(plugs...)
	c.Stdin = plugio.Stdin(plugs...)
	c.Stderr = plugio.Stderr(plugs...)

	for k, v := range deps {
		c.Args = append(c.Args, fmt.Sprintf("%s@%s", k, v))
	}

	for _, p := range plugs {
		if cr, ok := p.(CommandRunner); ok {
			if err := cr.RunWebpackCommand(ctx, root, c); err != nil {
				return err
			}
		}
	}

	if err := c.Run(); err != nil {
		return fmt.Errorf("%v: %w", c.Args, err)
	}
	return nil
}

var dependencies = map[string]string{
	"bootstrap":                     "latest",
	"popper.js":                     "latest",
	"@fortawesome/fontawesome-free": "latest",
	"jquery":                        "latest",
	"jquery-ujs":                    "latest",
}

var devDependencies = map[string]string{
	"@babel/cli":                    "latest",
	"@babel/core":                   "latest",
	"@babel/preset-env":             "latest",
	"babel-loader":                  "latest",
	"copy-webpack-plugin":           "latest",
	"css-loader":                    "latest",
	"expose-loader":                 "latest",
	"file-loader":                   "latest",
	"gopherjs-loader":               "latest",
	"mini-css-extract-plugin":       "latest",
	"node-sass":                     "latest",
	"npm-install-webpack-plugin":    "latest",
	"sass-loader":                   "latest",
	"style-loader":                  "latest",
	"terser-webpack-plugin":         "latest",
	"url-loader":                    "latest",
	"webpack":                       "latest",
	"webpack-clean-obsolete-chunks": "latest",
	"webpack-cli":                   "latest",
	"webpack-livereload-plugin":     "latest",
	"webpack-manifest-plugin":       "latest",
}

const tmpl = `
{
  "name": "{{.Name}}",
  "version": "1.0.0",
  "main": "index.js",
  "license": "MIT"
}
`
