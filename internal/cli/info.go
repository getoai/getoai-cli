package cli

import (
	"fmt"

	"github.com/getoai/getoai-cli/internal/tools"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <tool>",
	Short: "Show detailed information about a tool",
	Long:  `Display detailed information about a specific AI tool.`,
	Args:  cobra.ExactArgs(1),
	Run:   runInfo,
}

func runInfo(cmd *cobra.Command, args []string) {
	name := args[0]

	tool, ok := tools.Get(name)
	if !ok {
		printError(fmt.Sprintf("Unknown tool: %s", name))
		fmt.Println("Use 'getoai list' to see all available tools")
		return
	}

	fmt.Println()
	fmt.Printf("Name:        %s\n", tool.Name)
	fmt.Printf("Description: %s\n", tool.Description)
	fmt.Printf("Category:    %s\n", tool.Category)
	fmt.Printf("Website:     %s\n", tool.Website)

	if tool.IsInstalled() {
		fmt.Printf("Status:      \033[32mInstalled\033[0m\n")
		fmt.Printf("Version:     %s\n", tool.GetVersion())
	} else {
		fmt.Printf("Status:      \033[31mNot installed\033[0m\n")
	}

	methods := tool.GetAvailableMethods()
	if len(methods) > 0 {
		fmt.Printf("Install via: ")
		for i, m := range methods {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%s", m)
		}
		fmt.Println()
	}
	fmt.Println()
}
