# Plugins

When writing applications extensibility is important. Allowing users to plug into, and extend the functionality of, your application can prove vital to a projects success. When that application is command-line based, the problem of adding that extensibility becomes difficult. Go is a compiled language, and because of this external code can not be loaded at runtime, unlike dynamic languages such as Ruby or Python. Go provides a [`plugin`](https://godoc.org/plugin) package, but it does not work on all platforms.

While working on `BUFFALO`, we decided to add plugin support to our command line tooling by adopting the following strategy:

* Plugins must be named in the format of `buffalo-<plugin-name>`. For example, `buffalo-myplugin`.
* Plugins must be executable and must be available in one of the following places:
  * in the `$BUFFALO_PLUGIN_PATH`
  * if not set, `$GOPATH/bin`, is tried
  * in the `./plugins` folder of your buffalo application
* Plugins must implement an `available` command that prints a JSON response listing the available commands.

This strategy failed spectacularly and has become a source of confusion, bugs, and issues.

* Slow
* Works by finding executables in `PATH` and interrogates them for information
* Hard to development, maintain, test, use
* Caching. :(
* Currently writing plugins requires many dependencies, using `cobra`, and buffalo-centric code and idioms

In addition to those problems with adding plugins in this way, is that the executables, and therefore the plugins themselves, are not versioned control. The `buffalo` command-line tool faced a similar versioning problem.

To solve these problems, and others, I wanted to put the end user in charge of the tooling, to let them decide what happens after the command `buffalo ...` is run. The tooling should be an importable library, that anyone can import and use. It should also be simple to configure and use and plugin registration should be a simple as appending, or pre-pending, to a slice.

```go
package main

import (
  "context"
  "log"
  "os"

  "github.com/buffalo/buffalo-cli/cli"
  "github.com/buffalo/buffalo-cli/plugins"
)

func main() {
  if err := run(); err != nil {
    log.Fatal(err)
  }
}

func run() error {
  buffalo, err := cli.New()
  if err != nil {
    return err
  }

  buffalo.Plugins = append([]plugins.Plugin{
    // prepend your plugins here
  }, buffalo.Plugins...)
  return buffalo.Main(context.Background(), os.Args[1:])
}
```

The end user is now in charge of adding, or removing plugins. This puts control back in the user's hands and locks the versioning of those plugins inside of the `go.mod` file.

With an understanding of how to solve the biggest problems with the current plugin system, the next problem was to design a new plugin system. Before doing so, a set of guidelines were established:

* Everything must be a plugin, including anything that was previously a "hard-coded" sub-command of `buffalo`.
* Plugins must be independent of each other.
* Plugins should be responsible for their own interfaces.
* Interfaces should be 1 or 2 methods, no more.
* Interfaces should use only standard library types.

It was decided to use a minimal interface for becoming a plugin.

```go
type Plugin interface {
  Name() string
}
```

This small interface provides no real functionality, but makes for an easier entry point to the plugin system, allowing plugin developers to quickly see their plugin compiling and working.

## Guidelines

* plugins should have working zero-values
* have an `ifaces.go` file that lists all of the interfaces the plugin listens for
* when naming plugins use a simple, cli friendly name like `assets`, `develop`, etc....
* when naming multiple plugins in a package use a `/` to seperate the name of the package from the name of the plugin like `assets/build`, `assets/develop`, etc...
* if implementing a plugin that is a sub-command of another plugin, implement `plugins.NamedCommand` to specify a command name. the default is `path.Base(<plugin>.Name())`
* if a package has multiple plugins, provide a `Plugins() []plugins.Plugin` function that provides zero value versions of all plugins.

## Interfaces

* interfaces should be standard library only.
* interfaces should be 1-2 methods max
* interfaces should encourage using `context.Context`
* interface methods should narrow in scope and focus

A plugin that is bringing in a 3rd party package, such as the `pop`, `fizz`, or `plush` plugins may offer up their own interfaces that use the 3rd party package types.

## Dependencies

This `plugins` package will **ALWAYS** have **zero** dependencies.

## Working with Other Plugins

Implement `plugins.PluginNeeder` to receive a function that returns a list of all plugins.

When using this with `buffalo-cli/cli` this will be called with a function that contains all of the plugins registered with the cli.

```go
type Xyz struct {
  pluginsFn plugins.PluginFeeder
}

func (xyz *Xyz) WithPlugins(fn plugins.PluginFeeder) {
  xyz.pluginsFn = fn
}
```

## Scoping Plugins

* implement `plugins.PluginScoper` to return a list of plugins that are to be used by the plugin.

```go
type Xyz struct {
  pluginsFn plugins.PluginFeeder
}

func (xyz *Xyz) ScopedPlugins() []plugins.Plugin {
  var plugs []plugins.Plugin
  if xyz.pluginsFn == nil {
    return plugs
  }

  for _, p := range xyz.pluginsFn() {
    switch p.(type) {
    case InterfaceA:
      plugs = append(plugs, p)
    case InterfaceB:
      plugs = append(plugs, p)
    }
  }

  return plugs
}
```

## Flags

Flags should be exported fields on `struct`. If it can be called via a CLI's args, it can be called via API access.

It is recommended to cache flags and provide a simple function that returns the cached flags, if parsed, or build a new flag set.

```go
func (xyz *Xyz) Flags() *flag.FlagSet
func (xyz *Xyz) Flags() *pflag.FlagSet
```

```go
type Xyz struct {
  flags     *flag.FlagSet
}

func (xyz *Xyz) Flags() *flag.FlagSet {
  if xyz.flags != nil && xyz.flags.Parsed() {
    return xyz.flags
  }

  flags := flag.NewFlagSet(xyz.Name(), flag.ContinueOnError)

  // ...

  xyz.flags = flags

  return xyz.flags
}

func (xyz *Xyz) PrintFlags(w io.Writer) error {
  flags := xyz.Flags()
  flags.SetOutput(w)
  flags.PrintDefaults()
  return nil
}
```

### Sub-Command Flags

To have sub-commands of your Plugin, it is recommended to create interfaces to allow other plugins to declare their flags for your plugins.

Examples of interfaces that use [`flag`](https://godoc.org/flag) or [`github.com/spf13/pflag`](https://godoc.org/github.com/spf13/pflag) flag sets.

```go
type Flagger interface {
  plugins.Plugin
  XyzFlags() []*flag.Flag
}

type Pflagger interface {
  plugins.Plugin
  XyzFlags() []*pflag.Flag
}
```

```go
type Xyz struct {
  flags     *flag.FlagSet
}

// Flags returns a defined set of flags for this command.
// It imports flags provided by plugins that use either
// the `Flagger` or `Pflagger` interfaces. Flags provided
// by plugins will have their shorthand ("-x") flag stripped
// and the name ("--some-flag") of the flag will be
// prefixed with the plugin's name ("--xyz-some-flag")
func (xyz *Xyz) Flags() *flag.FlagSet {
  if xyz.flags != nil {
    return xyz.flags
  }

  flags := flag.NewFlagSet(xyz.Name(), flag.ContinueOnError)

  // ...

  for _, p := range xyz.ScopedPlugins() {
    switch t := p.(type) {
    case Flagger:
      for _, f := range plugins.CleanFlags(p, t.XyzFlags()) {
        flags.Var(f.Value, f.Name, f.Usage)
      }
    case Pflagger:
      // do work
    }
  }

  xyz.flags = flags

  return xyz.flags
}
```
