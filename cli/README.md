# CLI

## Installation

```bash
$ go get github.com/gobuffalo/buffalo-cli/cmd/buffalo
```

## Sub-Command

To be a sub-command of the buffalo tooling, `buffalo xyz`, a plugin must implement [`gobuffalo/buffalo-cli/cli#Command`](https://godoc.org/github.com/gobuffalo/buffalo-cli/cli#Command).

```go
type Xyz struct {
  info here.Info
}
```

## Getting here.Info

Plugins implementing `WithHere` will get called when `cli.Buffalo#Main` is called.

```go
type Xyz struct {
  info here.Info
}

func (xyz *Xyz) WithHereInfo(i here.Info) {
  xyz.info = i
}

func (xyz *Xyz) HereInfo() (here.Info, error) {
  if !xyz.info.IsZero() {
    return xyz.info, nil
  }
  return here.Current()
}
```
