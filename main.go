package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type Config struct {
	PRD           string `json:"prd"`
	Progress      string `json:"progress"`
	Prompt        string `json:"prompt"`
	Model         string `json:"model"`
	MaxIterations int    `json:"maxIterations"`
	Tools         Tools  `json:"tools"`
	LogLevel      string `json:"loglevel"`
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

func buildPrompt(promptTemplate string, prd string, progress string) (string, error) {
	type PromptVars struct {
		PRD      string
		PROGRESS string
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

func buildToolArgs(args []string, tools Tools) []string {
	for _, tool := range tools.Allow {
		args = append(args, "--allow-tool", tool)
	}

	for _, tool := range tools.Deny {
		args = append(args, "--deny-tool", tool)
	}

	return args
}

func runCopilot(prompt string, model string, tools Tools, logLevel string) (string, error) {
	args := []string{"--prompt", prompt, "--model", model, "--log-level", logLevel, "--share", "--silent"}
	args = buildToolArgs(args, tools)

	cmd := exec.Command("copilot", args...)

	log.Printf("Running Copilot: %v", cmd.Args)
	result, err := cmd.CombinedOutput()

	if err != nil {
		return string(result), err
	}

	return string(result), err

}

func main() {
	configPath := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	log.Printf("Config: %+v\n", config)

	prd, err := readFile(config.PRD)
	if err != nil {
		log.Fatalf("Could not load PRD: %v", err)
	}

	prompTmpl, err := readFile(config.Prompt)
	if err != nil {
		log.Printf("Could not load prompt template: %v", err)
	}

	for i := 0; i < config.MaxIterations; i++ {

		progress, err := readFile(config.Progress)
		if err != nil {
			log.Printf("Could not load progress: %v", err)
		}

		prompt, err := buildPrompt(prompTmpl, prd, progress)
		if err != nil {
			log.Printf("Could not build prompt: %v", err)
		}
		log.Printf("Built prompt: %v", prompt)

		result, err := runCopilot(prompt, config.Model, config.Tools, config.LogLevel)
		if err != nil {
			log.Printf("Could not run Copilot CLI: %v", err)
		}
		log.Printf("Here is the result: %v", result)

		if strings.Contains(result, "<promise>COMPLETE</promise>") {
			os.Exit(0)
		}
}
}
