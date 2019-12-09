package infocmd

import (
	"context"
)

// Informer can be implemented to add more checks
// to `buffalo info`
type Informer interface {
	Name() string
	Info(ctx context.Context, args []string) error
}
