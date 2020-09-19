package mainfix

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Cmd_Fix(t *testing.T) {
	r := require.New(t)

	dir, err := ioutil.TempDir("", "")
	r.NoError(err)

	f, err := os.Create(filepath.Join(dir, "main.go"))
	r.NoError(err)
	f.WriteString(`package main
	func main() {

	}
	`)
	r.NoError(f.Close())

	f, err = os.Create(filepath.Join(dir, "go.mod"))
	r.NoError(err)
	f.WriteString("module coke")
	r.NoError(f.Close())

	ctx := context.Background()
	var args []string

	fixer := &Cmd{}
	err = fixer.Fix(ctx, dir, args)
	r.NoError(err)

	mainFolder := filepath.Join(dir, "cmd", "coke")

	_, err = os.Stat(mainFolder)
	r.NoError(err)

	fp := filepath.Join(mainFolder, "main.go")
	b, err := ioutil.ReadFile(fp)
	r.NoError(err)
	r.Contains(string(b), `package `)
}

func Test_Cmd_NoMain(t *testing.T) {
	r := require.New(t)

	dir, err := ioutil.TempDir("", "")
	r.NoError(err)

	f, err := os.Create(filepath.Join(dir, "go.mod"))
	r.NoError(err)
	f.WriteString("module coke")
	r.NoError(f.Close())

	fixer := &Cmd{}
	r.NoError(fixer.Fix(context.Background(), dir, []string{}))
}

func Test_Cmd_Exists(t *testing.T) {
	r := require.New(t)

	dir, err := ioutil.TempDir("", "")
	r.NoError(err)

	f, err := os.Create(filepath.Join(dir, "go.mod"))
	r.NoError(err)
	f.WriteString("module coke")
	r.NoError(f.Close())

	r.NoError(os.MkdirAll(filepath.Join(dir, "cmd", "coke"), 0755))
	f, err = os.Create(filepath.Join(dir, "cmd", "coke", "main.go"))
	r.NoError(err)
	f.WriteString(`package main
	func main() {

	}
	`)
	r.NoError(f.Close())

	fixer := &Cmd{}
	r.NoError(fixer.Fix(context.Background(), dir, []string{}))
}
