package test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/pop/v5"
)

type Ref struct {
	here.Info
	TX *pop.Connection
}

func (r *Ref) Context() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "tx", r.TX)
	return ctx
}

func testerRef() (*Ref, error) {
	b := make([]byte, 10)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	dir := hex.EncodeToString(b)
	dir = filepath.Join("internal", "testdata", "temp", dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	cd := &pop.ConnectionDetails{
		URL: fmt.Sprintf("sqlite://%s", filepath.Join(dir, "test.db")),
	}
	if err := cd.Finalize(); err != nil {
		return nil, err
	}
	// json.NewEncoder(os.Stdout).Encode(cd)

	tx, err := pop.NewConnection(cd)
	if err != nil {
		return nil, err
	}
	if err := tx.Open(); err != nil {
		return nil, err
	}

	y := fmt.Sprintf(dbYml, filepath.Join(dir, "test.db"))
	err = writeFile(filepath.Join(dir, "database.yml"), y)
	if err != nil {
		return nil, err
	}

	ref := &Ref{
		Info: here.Info{
			Dir: dir,
		},
		TX: tx,
	}
	return ref, nil
}

func writeSchema(info here.Info, schema string) error {
	fp := filepath.Join(info.Dir, "migrations", "schema.sql")
	return writeFile(fp, schema)
}

func writeFile(fp string, body string) error {
	if err := os.MkdirAll(filepath.Dir(fp), 0777); err != nil {
		return err
	}
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	_, err = f.WriteString(body)
	if err != nil {
		return err
	}
	return f.Close()
}

// type Widget struct {
// 	ID        uuid.UUID
// 	Name      string    `db:"name"`
// 	CreatedAt time.Time `json:"created_at" db:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
// }

const dbYml = `test:
  dialect: "sqlite3"
  database: "%s"
`

const dbEmptySchema = `
CREATE TABLE IF NOT EXISTS "schema_migration" (
"version" TEXT NOT NULL
);
CREATE UNIQUE INDEX "schema_migration_version_idx" ON "schema_migration" (version);
`

const dbWidgetsSchema = `
CREATE TABLE IF NOT EXISTS "schema_migration" (
"version" TEXT NOT NULL
);
CREATE UNIQUE INDEX "schema_migration_version_idx" ON "schema_migration" (version);
CREATE TABLE IF NOT EXISTS "widgets" (
"id" TEXT PRIMARY KEY,
"name" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
`

const dbWidgetsMigration = `
create_table("widgets") {
	t.Column("id", "uuid", {primary: true})
	t.Column("name", "string", {})
	t.Timestamps()
}
`
