package scripts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

// ScriptFor will attempt to find the named script in the
// package.json file of the application.
func ScriptFor(root string, name string) (string, error) {
	scripts := packageJSON{}

	pf, err := os.Open(filepath.Join(root, "package.json"))
	if err != nil {
		return "", err
	}
	defer pf.Close()

	if err = json.NewDecoder(pf).Decode(&scripts); err != nil {
		return "", err
	}

	if s, ok := scripts.Scripts[name]; ok {
		return s, nil
	}
	return "", fmt.Errorf("script %q not found", name)
}

// WebpackBin returns the location of the webpack binary
func WebpackBin(root string) string {
	s := filepath.Join(root, "node_modules", ".bin", "webpack")
	if runtime.GOOS == "windows" {
		s += ".cmd"
	}
	return s
}

// Tool tries to determine whether to use yarn or npm
func Tool(root string) (string, error) {
	if _, err := os.Stat(filepath.Join(root, "yarn.lock")); err == nil {
		return "yarnpkg", nil
	}

	if _, err := os.Stat(filepath.Join(root, "package.json")); err == nil {
		return "npm", nil
	}

	return "", fmt.Errorf("could not determine asset tool from %q", root)
}
