package cli

import (
	"fmt"
	"sort"

	"github.com/getoai/getoai-cli/internal/tools"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for AI tools",
	Long: `Search for AI tools by name or description.

Examples:
  getoai search llm
  getoai search coding
  getoai search chat`,
	Args: cobra.ExactArgs(1),
	Run:  runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) {
	query := args[0]
	results := tools.Search(query)

	if len(results) == 0 {
		fmt.Printf("No tools found matching '%s'\n", query)
		fmt.Println("Use 'getoai list' to see all available tools")
		return
	}

	// Sort by name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	fmt.Printf("\nFound %d tool(s) matching '%s':\n\n", len(results), query)
	fmt.Printf("%-15s %-10s %-10s %s\n", "NAME", "CATEGORY", "STATUS", "DESCRIPTION")
	fmt.Printf("%-15s %-10s %-10s %s\n", "----", "--------", "------", "-----------")

	for _, tool := range results {
		status := "\033[31m✗\033[0m"
		if tool.IsInstalled() {
			status = "\033[32m✓\033[0m"
		}
		fmt.Printf("%-15s %-10s %-10s %s\n", tool.Name, tool.Category, status, tool.Description)
	}
	fmt.Println()
}
