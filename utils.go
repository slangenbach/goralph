package main

import (
	"os"
	"path/filepath"

	_ "embed"
)

const defaultConfigDir = ".ralph"
const configFile = "config.json"
const prdFile = "prd.json"
const promptFile = "prompt.md"

//go:embed templates/config.json
var defaultConfig string

//go:embed templates/prd.json
var defaultPRD string

//go:embed templates/prompt.md
var defaultPrompt string

func readFile(path string) (string, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return string(f), err
	}

	return string(f), err
}

func initFiles() {
	files := map[string]string{
		configFile: defaultConfig,
		prdFile:    defaultPRD,
		promptFile: defaultPrompt,
	}

	os.Mkdir(defaultConfigDir, 0755)

	for path, content := range files {
		fullPath := filepath.Join(defaultConfigDir, path)
		_, err := os.Stat(fullPath)
		if !os.IsNotExist(err) {
			continue
		}
		os.WriteFile(fullPath, []byte(content), 0644)
	}
}
