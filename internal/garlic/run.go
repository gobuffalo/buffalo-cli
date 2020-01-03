package garlic

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli"
	"github.com/gobuffalo/here"
	"github.com/markbates/jim"
)

type tasker interface {
	Task() *jim.Task
}

func Run(ctx context.Context, args []string) error {
	info, err := here.Dir(".")
	if err != nil {
		return err
	}

	ip := path.Join(info.Module.Path, "cli")

	bargs, err := buildTags(ctx, info)
	if err != nil {
		return err
	}

	t := &jim.Task{
		Info:      info,
		Args:      args,
		BuildArgs: bargs,
		Pkg:       ip,
		Sel:       "cli",
		Name:      "Buffalo",
	}

	err = jim.Run(ctx, t)
	if err == nil {
		return nil
	}

	if _, ok := err.(tasker); !ok {
		return err
	}

	if err := NewApp(ctx, info.Dir, args); err != nil {
		return err
	}

	err = jim.Run(ctx, t)
	if err == nil {
		return nil
	}

	b, err := cli.New()
	if err != nil {
		return err
	}
	return b.Main(ctx, args)

}

func buildTags(ctx context.Context, info here.Info) ([]string, error) {
	var args []string
	dy := filepath.Join(info.Dir, "database.yml")
	if _, err := os.Stat(dy); err != nil {
		return args, nil
	}

	b, err := ioutil.ReadFile(dy)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(b, []byte("sqlite")) {
		args = append(args, "-tags", "sqlite")
	}
	return args, nil
}
