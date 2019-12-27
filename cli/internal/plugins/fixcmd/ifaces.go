package fixcmd

import (
	"context"
)

// Fixer is an optional interface a plugin can implement
// to be run with `buffalo fix`. This should update the application
// to the current version of the plugin.
// The expectation is fixing of only one major revision.
type Fixer interface {
	Fix(ctx context.Context, args []string) error
}
