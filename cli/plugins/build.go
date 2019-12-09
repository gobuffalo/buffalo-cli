package plugins

import "context"

type BeforeBuilder interface {
	BeforeBuild(ctx context.Context, args []string) error
}

type AfterBuilder interface {
	AfterBuild(ctx context.Context, args []string) error
}
