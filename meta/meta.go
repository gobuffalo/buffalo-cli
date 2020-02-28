package meta

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func IsBuffalo(mod string) bool {
	if !strings.HasPrefix(mod, "go.mod") {
		mod = filepath.Join(mod, "go.mod")
	}

	b, err := ioutil.ReadFile(mod)
	if err != nil {
		return false
	}

	return bytes.Contains(b, []byte("github.com/gobuffalo/buffalo"))
}
