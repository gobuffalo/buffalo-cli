package packr

import (
	"context"

	"github.com/gobuffalo/packr/v2/jam"
)

type Buffalo struct{}

func (b *Buffalo) BeforeBuild(ctx context.Context, args []string) error {
	return jam.Clean()
}

func (b *Buffalo) Package(ctx context.Context, root string) error {
	return jam.Pack(jam.PackOptions{
		Roots: []string{root},
	})
}

func (b Buffalo) Name() string {
	return "packr"
}
