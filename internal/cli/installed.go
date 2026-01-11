package cli

import (
	"fmt"
	"sort"

	"github.com/getoai/getoai-cli/internal/tools"
	"github.com/spf13/cobra"
)

var installedCmd = &cobra.Command{
	Use:   "installed",
	Short: "List installed AI tools",
	Long:  `Display a list of AI tools that are currently installed on your system.`,
	Run:   runInstalled,
}

func runInstalled(cmd *cobra.Command, args []string) {
	allTools := tools.List()

	var installed []*tools.Tool
	for _, tool := range allTools {
		if tool.IsInstalled() {
			installed = append(installed, tool)
		}
	}

	// Sort by name
	sort.Slice(installed, func(i, j int) bool {
		return installed[i].Name < installed[j].Name
	})

	if len(installed) == 0 {
		fmt.Println("No AI tools installed yet.")
		fmt.Println("Use 'getoai list' to see available tools")
		fmt.Println("Use 'getoai install <tool>' to install one")
		return
	}

	fmt.Println()
	fmt.Printf("%-15s %-10s %s\n", "NAME", "CATEGORY", "VERSION")
	fmt.Printf("%-15s %-10s %s\n", "----", "--------", "-------")

	for _, tool := range installed {
		fmt.Printf("%-15s %-10s %s\n", tool.Name, tool.Category, tool.GetVersion())
	}
	fmt.Println()
}
