package fix

import (
	"fmt"

	bufcli "github.com/gobuffalo/buffalo-cli"
	"github.com/gobuffalo/meta"
)

// Check interface for runnable checker functions
type Check func(*Runner) error

// Runner will run all compatible checks
type Runner struct {
	App      meta.App
	Warnings []string
}

// Run all compatible checks
func Run() error {
	fmt.Printf("! This updater will attempt to update your application to Buffalo version: %s\n", bufcli.Version)
	if !ask("Do you wish to continue?") {
		fmt.Println("~~~ cancelling update ~~~")
		return nil
	}

	r := &Runner{
		App:      meta.New("."),
		Warnings: []string{},
	}

	defer func() {
		if len(r.Warnings) == 0 {
			return
		}

		fmt.Println("\n\n----------------------------")
		fmt.Printf("!!! (%d) Warnings Were Found !!!\n\n", len(r.Warnings))
		for _, w := range r.Warnings {
			fmt.Printf("[WARNING]: %s\n", w)
		}
	}()

	for _, c := range checks {
		if err := c(r); err != nil {
			return err
		}
	}
	return nil
}
