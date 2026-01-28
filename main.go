package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	_ "embed"
)

const exitCondition = "<promise>COMPLETE</promise>"

//go:embed templates/config.json
var defaultConfig string

//go:embed templates/prd.json
var defaultPRD string

//go:embed templates/prompt.md
var defaultPrompt string

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

func createLogger(logLevel string) {
	var level slog.Level

	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))
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

	slog.Debug("Running Copilot", "args", cmd.Args)
	result, err := cmd.CombinedOutput()

	if err != nil {
		return string(result), err
	}

	return string(result), err

}

func initFiles() {
	files := map[string]string{
		"config.json": defaultConfig,
		"prd.json":    defaultPRD,
		"prompt.md":   defaultPrompt,
	}

	os.Mkdir(".ralph", 0755)
	for path, content := range files {
		fullPath := filepath.Join(".ralph", path)
		_, err := os.Stat(fullPath)
		if !os.IsNotExist(err) {
			continue
		}
		os.WriteFile(fullPath, []byte(content), 0644)
	}
}

func main() {
	configPath := flag.String("config", ".ralph/config.json", "Path to config file")
	doInitFiles := flag.Bool("init", false, "Generate sample configs")
	flag.Parse()

	if *doInitFiles {
		initFiles()
		slog.Info("Successfully initialized config files")
		os.Exit(0)
	}

	config, err := loadConfig(*configPath)
	createLogger(config.LogLevel)
	if err != nil {
		slog.Error("Could not load config", "err", err)
		os.Exit(1)
	}

	prd, err := readFile(config.PRD)
	if err != nil {
		slog.Error("Could not load PRD: ", "err", err)
		os.Exit(1)
	}

	promptTmpl, err := readFile(config.Prompt)
	if err != nil {
		slog.Error("Could not load prompt template", "err", err)
		os.Exit(1)
	}

	for i := 0; i < config.MaxIterations; i++ {
		slog.Debug("Running iteration", "iter", i)

		progress, err := readFile(config.Progress)
		if err != nil {
			slog.Warn("Could not load progress", "err", err)
		}

		prompt, err := buildPrompt(promptTmpl, prd, progress)
		if err != nil {
			slog.Error("Could not build prompt", "err", err)
			os.Exit(1)
		}
		slog.Debug("Built prompt", "prompt", prompt)

		result, err := runCopilot(prompt, config.Model, config.Tools, config.LogLevel)
		if err != nil {
			slog.Error("Could not run Copilot CLI", "err", err)
			os.Exit(1)
		}
		slog.Debug("Here is the result", "result", result)

		if strings.Contains(result, exitCondition) {
			os.Exit(0)
		}
	}

	slog.Warn("Reached max iterations without completion", "iterations", config.MaxIterations)
	os.Exit(2)
}
