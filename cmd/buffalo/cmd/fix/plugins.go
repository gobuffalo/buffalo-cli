package fix

import (
	"context"
	"fmt"
	"os"

	"github.com/gobuffalo/bufcli/genny/plugins/install"
	"github.com/gobuffalo/bufcli/plugins"
	"github.com/gobuffalo/bufcli/plugins/plugdeps"
	"github.com/gobuffalo/genny"
	"github.com/markbates/errx"
)

// Plugins will fix plugins between releases
func Plugins(r *Runner) error {
	fmt.Println("~~~ Cleaning plugins cache ~~~")
	os.RemoveAll(plugins.CachePath)
	plugs, err := plugdeps.List(r.App)
	if err != nil && (errx.Unwrap(err) != plugdeps.ErrMissingConfig) {
		return err
	}

	run := genny.WetRunner(context.Background())
	gg, err := install.New(&install.Options{
		App:     r.App,
		Plugins: plugs.List(),
	})

	run.WithGroup(gg)

	fmt.Println("~~~ Reinstalling plugins ~~~")
	return run.Run()
}
