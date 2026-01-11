package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "getoai",
	Short: "One-click installer for AI tools and CLIs",
	Long: `getoai-cli is a cross-platform CLI tool for installing and managing
AI-related tools and command-line programs.

Supported tools include:
  - ollama      : Run LLMs locally
  - claude-code : Claude AI coding assistant
  - aider       : AI pair programming
  - llm         : LLM CLI by Simon Willison
  - open-webui  : Web UI for Ollama
  - cursor      : AI-first code editor
  - and more...

Examples:
  getoai list                    # List all available tools
  getoai install ollama          # Install ollama
  getoai install claude-code     # Install Claude Code CLI
  getoai info aider              # Show info about aider`,
	Version: Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(installedCmd)
}

func printSuccess(msg string) {
	fmt.Printf("\033[32m✓\033[0m %s\n", msg)
}

func printError(msg string) {
	fmt.Printf("\033[31m✗\033[0m %s\n", msg)
}

func printInfo(msg string) {
	fmt.Printf("\033[34mℹ\033[0m %s\n", msg)
}
