package cli

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/gobuffalo/pop"
)

// New is WIP!! DON'T USE IT!!!
func (b *Buffalo) New(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("buffalo new", flag.ContinueOnError)

	opts := struct {
		API         bool
		Force       bool
		DryRun      bool
		Verbose     bool
		SkipPop     bool
		SkipWebpack bool
		SkipYarn    bool
		DBType      string
		CI          string
		VCS         string
	}{}

	flags.BoolVar(&opts.API, "api", false, "skip all front-end code and configure for an API server")
	flags.BoolVar(&opts.Force, "force", false, "delete and remake if the app already exists")
	flags.BoolVar(&opts.Force, "f", false, "delete and remake if the app already exists")
	flags.BoolVar(&opts.DryRun, "d", false, "dry run")
	flags.BoolVar(&opts.Verbose, "v", false, "verbosely print out the go get commands")
	flags.BoolVar(&opts.SkipPop, "skip-pop", false, "skips adding pop/soda to your app")
	flags.BoolVar(&opts.SkipWebpack, "skip-webpack", false, "skips adding Webpack to your app")
	flags.BoolVar(&opts.SkipYarn, "skip-yarn", false, "use npm instead of yarn for frontend dependencies management")
	flags.StringVar(&opts.DBType, "db-type", "postgres", fmt.Sprintf("specify the type of database you want to use [%s]", strings.Join(pop.AvailableDialects, ", ")))
	flags.StringVar(&opts.CI, "ci-provider", "none", "specify the type of ci file you would like buffalo to generate [none, travis, gitlab-ci]")
	flags.StringVar(&opts.VCS, "vcs", "git", "specify the Version control system you would like to use [none, git, bzr]")

	if err := flags.Parse(args); err != nil {
		return err
	}

	return nil
}
