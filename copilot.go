package main

import (
	"context"
	"log/slog"
	"os/exec"
	"strings"
	"text/template"
)

const exitCondition = "<promise>COMPLETE</promise>"

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

func runCopilot(ctx context.Context, prompt string, model string, tools Tools, logLevel string) (string, error) {
	args := []string{"--prompt", prompt, "--model", model, "--log-level", logLevel, "--share", "--silent"}
	args = buildToolArgs(args, tools)

	cmd := exec.CommandContext(ctx, "copilot", args...)

	slog.Debug("Running Copilot", "args", cmd.Args)
	result, err := cmd.CombinedOutput()

	if err != nil {
		return string(result), err
	}

	return string(result), err
}
