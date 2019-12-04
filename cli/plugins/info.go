package plugins

import "context"

// Informer can be implemented to add more checks
// to `buffalo info`
type Informer interface {
	Info(ctx context.Context, args []string) error
}

func (plugs Plugins) Info(ctx context.Context, args []string) error {
	for _, p := range plugs {
		i, ok := p.(Informer)
		if !ok {
			continue
		}
		if err := i.Info(ctx, args); err != nil {
			return err
		}
	}
	return nil
}
