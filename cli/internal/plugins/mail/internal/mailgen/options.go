package mailgen

import (
	"fmt"
)

// Options needed to create a new mailer
type Options struct {
	Args     []string `json:"args"`
	Name     string   `json:"name"`
	SkipInit bool     `json:"skip_init"`
}

// Validate options are useful
func (opts *Options) Validate() error {
	if len(opts.Name) == 0 {
		if len(opts.Args) == 0 {
			return fmt.Errorf("you must supply a name for your mailer")
		}
		opts.Name = opts.Args[0]
	}
	return nil
}
