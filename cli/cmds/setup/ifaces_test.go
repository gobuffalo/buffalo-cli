package setup

import "context"

var _ BeforeSetuper = beforeSetuper(nil)

type beforeSetuper func(ctx context.Context, root string, args []string) error

func (b beforeSetuper) PluginName() string {
	return "before-setuper"
}

func (b beforeSetuper) BeforeSetup(ctx context.Context, root string, args []string) error {
	return b(ctx, root, args)
}
