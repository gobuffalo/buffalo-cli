package install

import (
	"os/exec"

	"github.com/gobuffalo/buffalo-cli/v2/internal/v1/genny/add"
	"github.com/gobuffalo/genny/v2"
)

// New installs plugins and then added them to the config file
func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}

	if err := opts.Validate(); err != nil {
		return gg, err
	}

	aopts := &add.Options{
		App:     opts.App,
		Plugins: opts.Plugins,
	}

	if err := aopts.Validate(); err != nil {
		return gg, err
	}

	g := genny.New()
	for _, p := range opts.Plugins {
		if len(p.GoGet) == 0 {
			continue
		}

		var args []string
		if len(p.Tags) > 0 {
			args = append(args, "-tags", p.Tags.String())
		}
		args = append([]string{p.GoGet}, args...)
		g.Command(exec.Command("go", args...))
	}
	gg.Add(g)

	g, err := add.New(aopts)
	if err != nil {
		return gg, err
	}

	gg.Add(g)

	return gg, nil
}
