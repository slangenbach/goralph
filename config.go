package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	PRDFilePath      string `json:"prd"`
	ProgressFilePath string `json:"progress"`
	PromptFilePath   string `json:"prompt"`
	Model            string `json:"model"`
	Tools            Tools  `json:"tools"`
	Timeout          int    `json:"timeout"`
	LogLevel         string `json:"loglevel"`
}

type Tools struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

func loadConfig(path string) (Config, error) {
	var config Config

	file, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	return config, err
}
