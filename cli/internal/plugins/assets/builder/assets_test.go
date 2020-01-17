package builder

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/here"
)

func tempApp(t *testing.T, scripts map[string]string) here.Info {
	t.Helper()
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(filepath.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	sc := packageJSON{
		Scripts: scripts,
	}

	err = json.NewEncoder(f).Encode(sc)
	if err != nil {
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	return here.Info{
		Root: dir,
		Dir:  dir,
	}
}
