package i18n

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp"
	"github.com/gobuffalo/plugins"
)

var _ plugins.Plugin = &Newapp{}
var _ newapp.Newapper = &Newapp{}

type Newapp struct{}

func (na *Newapp) PluginName() string {
	return "i18n/newapp"
}

func (na *Newapp) Newapp(ctx context.Context, root string, args []string) error {
	f, err := os.Create(filepath.Join(root, "inflections.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")

	mm := map[string]string{
		"singular": "plural",
	}
	return enc.Encode(mm)
}
