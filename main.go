package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	configPath := flag.String("config", filepath.Join(defaultConfigDir, configFile), "Path to config file")
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

	prd, err := readFile(config.PRDFilePath)
	if err != nil {
		slog.Error("Could not load PRD: ", "err", err)
		os.Exit(1)
	}

	promptTmpl, err := readFile(config.PromptFilePath)
	if err != nil {
		slog.Error("Could not load prompt template", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Timeout)*time.Minute)
	defer cancel()

	slog.Info("Starting implementation")
	for {
		progress, err := readFile(config.ProgressFilePath)
		if err != nil {
			slog.Warn("Could not load progress", "err", err)
		}

		prompt, err := buildPrompt(promptTmpl, prd, progress)
		if err != nil {
			slog.Error("Could not build prompt", "err", err)
			os.Exit(1)
		}

		result, err := runCopilot(ctx, prompt, config.Model, config.Tools, config.LogLevel)
		if err != nil {
			if ctx.Err() != nil {
				slog.Warn("Reached timeout without completion", "timeout", config.Timeout)
				os.Exit(2)
			}
			slog.Error("Could not run Copilot CLI", "err", err)
			os.Exit(1)
		}
		slog.Debug("Here is the result", "result", result)

		if strings.Contains(result, exitCondition) {
			slog.Info("Implementation completed")
			os.Exit(0)
		}
	}
}
