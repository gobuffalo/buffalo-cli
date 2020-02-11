package tester

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/test"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/pop/v5"
)

var _ plugins.Plugin = &Tester{}
var _ test.Argumenter = &Tester{}
var _ test.BeforeTester = &Tester{}

type Tester struct {
}

func (t *Tester) TestArgs(ctx context.Context, root string) ([]string, error) {
	args := []string{"-p", "1"}

	dy := filepath.Join(root, "database.yml")
	if _, err := os.Stat(dy); err != nil {
		return args, nil
	}

	b, err := ioutil.ReadFile(dy)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(b, []byte("sqlite")) {
		args = append(args, "-tags", "sqlite")
	}
	return args, nil
}

func (Tester) PluginName() string {
	return "pop/tester"
}

func (t *Tester) BeforeTest(ctx context.Context, root string, args []string) error {
	if err := pop.AddLookupPaths(root); err != nil {
		return err
	}

	var err error
	db, ok := ctx.Value("tx").(*pop.Connection)
	if !ok {
		if _, err := os.Stat(filepath.Join(root, "database.yml")); err != nil {
			return err
		}

		db, err = pop.Connect("test")
		if err != nil {
			return err
		}
	}
	// drop the test db:
	db.Dialect.DropDB()

	// create the test db:
	if err := db.Dialect.CreateDB(); err != nil {
		return err
	}

	for _, a := range args {
		if a == "--force-migrations" {
			return t.forceMigrations(root, db)
		}
	}

	schema, err := t.findSchema(root)
	if err != nil {
		return err
	}
	if schema == nil {
		return t.forceMigrations(root, db)
	}

	if err = db.Dialect.LoadSchema(schema); err != nil {
		return err
	}
	return nil
}

func (t *Tester) forceMigrations(root string, db *pop.Connection) error {

	ms := filepath.Join(root, "migrations")
	fm, err := pop.NewFileMigrator(ms, db)
	if err != nil {
		return err
	}
	return fm.Up()
}

func (t *Tester) findSchema(root string) (io.Reader, error) {
	ms := filepath.Join(root, "migrations", "schema.sql")
	if f, err := os.Open(ms); err == nil {
		return f, nil
	}

	if dev, err := pop.Connect("development"); err == nil {
		schema := &bytes.Buffer{}
		if err = dev.Dialect.DumpSchema(schema); err == nil {
			return schema, nil
		}
	}

	return nil, nil
}
