package built

import "context"

type Initer interface {
	BuiltInit(ctx context.Context, args []string) error
}
