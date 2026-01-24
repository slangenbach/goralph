package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

type Config struct {
	PRD           string `json:"prd"`
	Progress      string `json:"progress"`
	Prompt        string `json:"prompt"`
	MaxIterations int    `json:"maxIterations"`
	Tools         Tools  `json:"tools"`
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

func readFile(path string) (string, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return string(f), err
	}

	return string(f), err
}

func loadPRD(path string) (string, error) {
	var prd string

	file, err := os.ReadFile(path)
	if err != nil {
		return prd, err
	}

	err = json.Unmarshal(file, &prd)
	if err != nil {
		return prd, err
	}

	return prd, err
}

func buildPrompt(promptTemplate string, prd string, progress string) (string, error) {
	type PromptVars struct {
		prd      string
		progress string
	}
	var prompt strings.Builder

	tmpl, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		return prompt.String(), err
	}

	err = tmpl.Execute(&prompt, PromptVars{prd, progress})
	if err != nil {
		return prompt.String(), err
	}

	return prompt.String(), err
}

func main() {
	configPath := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Loading config failed: %v", err)
	}

	fmt.Printf("Loaded config: %+v\n", config)

	prd, err := loadPRD(config.PRD)
	if err != nil {
		log.Fatalf("Could not load PRD: %v", err)
	}

	progress, err := readFile(config.Progress)
	if err != nil {
		log.Printf("Could not load progress: %v", err)
	}

	prompTmpl, err := readFile(config.Prompt)
	if err != nil {
		log.Printf("Could not load prompt template: %v", err)
	}

	prompt, err := buildPrompt(prompTmpl, prd, progress)
	if err != nil {
		log.Printf("Could not build prompt: %v", err)
	}

	fmt.Print(prompt)

}
