package newapp

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp"
	"github.com/gobuffalo/plugins"
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
		}
	}

	return plugs
}

func (g *Generator) Newapp(ctx context.Context, root string, args []string) error {
	if len(args) == 0 {
		return plugins.Wrap(g, fmt.Errorf("missing application name"))
	}

	name := args[0]

	t, err := template.New("package.json").Parse(tmpl)
	if err != nil {
		return err
	}

	data := struct {
		Name string
	}{
		Name: name,
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
	for k, v := range dependencies {
		c.Args = append(c.Args, fmt.Sprintf("%s@%s", k, v))
	}

	c.Stdout = os.Stdout
	if err := c.Run(); err != nil {
		return err
	}

	c = exec.CommandContext(ctx, tool, "add", "--dev")
	for k, v := range devDependencies {
		c.Args = append(c.Args, fmt.Sprintf("%s@%s", k, v))
	}
	c.Stdout = os.Stdout
	return c.Run()
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
