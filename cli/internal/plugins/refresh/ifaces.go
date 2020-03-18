package refresh

import (
	"context"

	"github.com/gobuffalo/plugins/plugio"
)

type Tagger interface {
	RefreshTags(ctx context.Context, root string) ([]string, error)
}

type Stdouter = plugio.Outer
