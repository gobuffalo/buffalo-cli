package plugins

import "context"

// Generator is an optional interface a plugin can implement
// to be run with `buffalo generate`.
type Generator interface {
	Generate(ctx context.Context, args []string) error
}

// Model is an optional interface a plugin can implement to be
// in charge of generate models during `buffalo generate resource`,
// `buffalo new`, and other generators that need it.
type Model interface {
	Model(ctx context.Context, args []string) error
}

type ModelTest interface {
	ModelTest(ctx context.Context, args []string) error
}

// Actions is an optional interface a plugin can implement to be
// in charge of generate actions during `buffalo generate resource`,
// `buffalo new`, and other generators that need it.
type Actions interface {
	Actions(ctx context.Context, args []string) error
}

type ActionsTest interface {
	ActionsTest(ctx context.Context, args []string) error
}

// Migrations is an optional interface a plugin can implement to be
// in charge of generate migrations during `buffalo generate resource`,
// `buffalo new`, and other generators that need it.
type Migrations interface {
	Migrations(ctx context.Context, args []string) error
}

type MigrationsTest interface {
	MigrationsTest(ctx context.Context, args []string) error
}

// Templates is an optional interface a plugin can implement to be
// in charge of generate templates during `buffalo generate resource`,
// `buffalo new`, and other generators that need it.
type Templates interface {
	Templates(ctx context.Context, args []string) error
}

type TemplatesTest interface {
	TemplatesTest(ctx context.Context, args []string) error
}

// Resource is a full implementation of `buffalo generate resource`
type Resource interface {
	Actions
	ActionsTest
	Migrations
	MigrationsTest
	Model
	ModelTest
	Templates
	TemplatesTest
}
