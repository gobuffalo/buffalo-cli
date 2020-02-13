package setup

import "context"

type Migrater interface {
	MigrateDB(ctx context.Context, root string, args []string) error
}
