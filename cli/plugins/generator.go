package plugins

import (
	"context"
)

// Fixer is an optional interface a plugin can implement
// to be run with `buffalo fix`
type Generator interface {
	Generate(ctx context.Context, args []string) error
}

// // Generate will run the specified generator.
// // 	buffalo generate -h
// // 	buffalo generate pop ...
// func (plugs Plugins) Generate(ctx context.Context, args []string) error {
// 	opts := struct {
// 		help bool
// 	}{}
//
// 	flags := cmdx.NewFlagSet("buffalo generate", cmdx.Stderr(ctx))
// 	flags.BoolVar(&opts.help, "h", false, "print this help")
// 	flags.Parse(args)
//
// 	args = flags.Args()
// 	if opts.help || len(args) == 0 {
// 		sort.Slice(plugs, func(i, j int) bool {
// 			return plugs[i].Name() < plugs[j].Name()
// 		})
//
// 		// stderr := cmdx.Stderr(ctx)
// 		for _, p := range plugs {
// 			if _, ok := p.(Generator); ok {
// 				// fmt.Fprintf(stderr, "%s %s - [%s]\n", flags.Name(), p.Name(), p)
// 			}
// 		}
// 		return nil
// 	}
//
// 	arg := args[0]
// 	if len(args) > 0 {
// 		args = args[1:]
// 	}
//
// 	for _, p := range plugs {
// 		f, ok := p.(Generator)
// 		if !ok {
// 			continue
// 		}
// 		if p.Name() != arg {
// 			continue
// 		}
//
// 		return f.Generate(ctx, args)
// 	}
// 	return fmt.Errorf("unknown generator %s", arg)
// }
