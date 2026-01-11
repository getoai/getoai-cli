package cli

import (
	"fmt"

	"github.com/getoai/getoai-cli/internal/tools"
	"github.com/getoai/getoai-cli/internal/util"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [tool...]",
	Short: "Update installed AI tools",
	Long: `Update one or more installed AI tools to their latest versions.
If no tool is specified, updates all installed tools.

Examples:
  getoai update ollama
  getoai update                  # Update all installed tools
  getoai update aider llm`,
	Run: runUpdate,
}

var updateAll bool

func init() {
	updateCmd.Flags().BoolVarP(&updateAll, "all", "a", false, "Update all installed tools")
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) {
	if len(args) == 0 || updateAll {
		updateAllTools()
		return
	}

	for _, toolName := range args {
		updateTool(toolName)
	}
}

func updateAllTools() {
	allTools := tools.List()
	var installed []*tools.Tool

	for _, tool := range allTools {
		if tool.IsInstalled() {
			installed = append(installed, tool)
		}
	}

	if len(installed) == 0 {
		printInfo("No AI tools installed to update")
		return
	}

	fmt.Printf("Updating %d installed tools...\n\n", len(installed))

	for _, tool := range installed {
		updateTool(tool.Name)
	}
}

func updateTool(name string) {
	tool, ok := tools.Get(name)
	if !ok {
		printError(fmt.Sprintf("Unknown tool: %s", name))
		return
	}

	if !tool.IsInstalled() {
		printInfo(fmt.Sprintf("%s is not installed", name))
		return
	}

	spinner := util.NewSpinner(fmt.Sprintf("Updating %s...", name))
	spinner.Start()

	methods := tool.GetAvailableMethods()
	if len(methods) == 0 {
		spinner.Error(fmt.Sprintf("No update method available for %s", name))
		return
	}

	// For most package managers, reinstalling updates to latest version
	err := tool.Install(methods[0])
	if err != nil {
		spinner.Error(fmt.Sprintf("Failed to update %s: %v", name, err))
		return
	}

	spinner.Success(fmt.Sprintf("%s updated to %s", name, tool.GetVersion()))
}
