package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/yomorun/yomo/cli"
	"github.com/yomorun/yomo/serverless"
)

// Init is an optional function invoked during the initialization phase
func Init() error {
	// init doc
	return nil
}

// Description returns a description of the function
func Description() string {
	return `
	Analyze the user input, clearly identify the 'yomo' subcommands that meet the user's input requirements, and provide detailed usage instructions. If there are no suitable subcommands, provide usage instructions for the 'yomo' root command

- Initialize a YoMo Serverless LLM Function, use the 'init' subcommand
- Build the YoMo Stream Function, use the 'build' subcommand
- Run a YoMo Serverless LLM Function, use the 'run' subcommand
-	Run a YoMo-Zipper, use the 'serve' subcommand. First, create the 'zipper.yml' file, and then run 'yomo serve -c zipper.yml' to start the service
- Get the version of YoMo, use the 'version' subcommand
`
}

// InputSchema defines the argument structure for LLM Function Calling. It
// utilizes jsonschema tags to detail the definition. For jsonschema in Go,
// see https://github.com/invopop/jsonschema.
func InputSchema() any {
	return &LLMArguments{}
}

// LLMArguments defines the arguments for the LLM Function Calling. These
// arguments are combined to form a prompt automatically.
type LLMArguments struct {
	Command string `json:"command" jsonschema:"description=yomo CLI subcommand, eg: init, build, run, serve, version"`
}

// Handler orchestrates the core processing logic of this function.
// - ctx.ReadLLMArguments() parses LLM Function Calling Arguments (skip if none).
// - ctx.WriteLLMResult() sends the retrieval result back to LLM.
func Handler(ctx serverless.Context) {
	var p LLMArguments
	// deserilize the arguments from llm tool_call response
	ctx.ReadLLMArguments(&p)

	var cmd string
	// simple command matching logic
	switch {
	case strings.Contains(p.Command, "init"):
		cmd = "init"
	case strings.Contains(p.Command, "build"):
		cmd = "build"
	case strings.Contains(p.Command, "run"):
		cmd = "run"
	case strings.Contains(p.Command, "serve"):
		cmd = "serve"
	case strings.Contains(p.Command, "version"):
		cmd = "version"
	default:
		cmd = "yomo"
	}

	doc, err := cli.Doc(cmd)
	if err != nil {
		ctx.WriteLLMResult(fmt.Sprintf("Error get document for command '%s': %v", p.Command, err))
		slog.Error("yomo-cli-mcp", "command", p.Command, "error", err)
		return
	}
	ctx.WriteLLMResult(doc)
	slog.Info("yomo-cli-mcp", "command", p.Command, "doc", doc)
}
