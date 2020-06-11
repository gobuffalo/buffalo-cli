package docker

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp"
	"github.com/gobuffalo/here"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
	"github.com/spf13/pflag"
)

var _ newapp.Newapper = &Generator{}

type Generator struct {
	style string
	flags *pflag.FlagSet

	BuffaloVersion string
	Tool           string
	WebPack        bool
	Name           string
}

func (Generator) PluginName() string {
	return "docker"
}

func (Generator) Description() string {
	return "Generates Dockerfile"
}

func (g Generator) Newapp(ctx context.Context, root string, name string, args []string) error {

	f, err := os.Create(filepath.Join(root, "Dockerfile"))
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl, err := g.buildTemplate()
	if err != nil {
		return err
	}

	version, err := g.imageTag()
	if err != nil {
		return err
	}

	info, err := here.Dir(root)
	if err != nil {
		return err
	}

	g.WebPack = g.hasWebpack(root)
	g.Tool = g.tool(root)
	g.BuffaloVersion = version
	g.Name = info.Name

	err = tmpl.Execute(f, g)
	if err != nil {
		return err
	}

	return nil
}

func (g Generator) hasWebpack(root string) bool {
	if _, err := os.Stat(filepath.Join(root, "webpack.config.js")); os.IsNotExist(err) {
		return false
	}

	return true
}

func (g Generator) tool(root string) string {
	_, err := os.Stat(filepath.Join(root, "yarn.lock"))
	if os.IsNotExist(err) {
		return "npm"
	}

	return "yarn"
}

func (Generator) imageTag() (string, error) {
	info, err := here.Package("github.com/gobuffalo/buffalo")
	if err != nil {
		return "", err
	}

	parts := strings.Split(info.Module.Dir, "/")
	version := parts[len(parts)-1]

	return version, nil
}

func (g Generator) buildTemplate() (*template.Template, error) {
	file, err := g.templateFile()
	if err != nil {
		return nil, err
	}

	t, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	template, err := template.New("Dockerfile").Parse(string(t))
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (g Generator) templateFile() (pkging.File, error) {
	td := pkger.Include("github.com/gobuffalo/buffalo-cli/v2:/cli/internal/plugins/docker/templates")

	file := "Dockerfile.multistage"
	if g.style == "standard" {
		file = "Dockerfile.standard"
	}

	return pkger.Open(filepath.Join(td, file))
}
