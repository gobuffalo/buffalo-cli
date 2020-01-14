package i18n

import (
	"bytes"
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/flect/name"
	"github.com/stretchr/testify/require"
)

func Test_Flash_Flash(t *testing.T) {
	r := require.New(t)

	bb := &bytes.Buffer{}

	flash := &flasher{
		pluginsFn: plugins.Plugins{
			&namedWriter{w: bb},
		}.ScopedPlugins,
	}

	model := name.New("widgets")

	err := flash.Flash(context.Background(), "testdata", model)
	r.NoError(err)

	wr, err := ioutil.ReadFile(filepath.Join("testdata", "actions", "widgets.go"))
	r.NoError(err)
	r.NotEqual(wr, bb.Bytes())

	res := bb.String()
	r.Contains(res, `T.Translate(c, "widget.updated.success")`)
}
