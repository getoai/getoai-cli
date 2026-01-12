package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/getoai/getoai-cli/internal/tools"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available AI tools",
	Long:  `Display a list of all AI tools that can be installed using getoai.`,
	Run:   runList,
}

var (
	listCategory string
	listGrouped  bool
)

func init() {
	listCmd.Flags().StringVarP(&listCategory, "category", "c", "", "Filter by category (llm, coding, ui, utility, platform, infra)")
	listCmd.Flags().BoolVarP(&listGrouped, "group", "g", false, "Group tools by category")
}

func runList(cmd *cobra.Command, args []string) {
	var toolList []*tools.Tool

	if listCategory != "" {
		toolList = tools.ListByCategory(tools.Category(listCategory))
	} else {
		toolList = tools.List()
	}

	if len(toolList) == 0 {
		fmt.Println("No tools found.")
		return
	}

	fmt.Println()

	if listGrouped && listCategory == "" {
		// Group display mode
		runListGrouped(toolList)
	} else {
		// Flat display mode
		runListFlat(toolList)
	}

	fmt.Printf("\n\033[1mTotal: %d tools available\033[0m\n", tools.Count())
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  getoai install <tool>   Install a tool")
	fmt.Println("  getoai info <tool>      Show tool details")
	fmt.Println("  getoai search <keyword> Search for tools")
	fmt.Println("  getoai list -g          Group by category")
}

func runListGrouped(toolList []*tools.Tool) {
	// Group tools by category
	grouped := make(map[tools.Category][]*tools.Tool)
	for _, tool := range toolList {
		grouped[tool.Category] = append(grouped[tool.Category], tool)
	}

	// Sort tools within each category
	for cat := range grouped {
		sort.Slice(grouped[cat], func(i, j int) bool {
			return grouped[cat][i].Name < grouped[cat][j].Name
		})
	}

	// Display by category order
	categories := tools.GetCategories()
	for _, cat := range categories {
		catTools := grouped[cat]
		if len(catTools) == 0 {
			continue
		}

		// Category header
		catName := tools.GetCategoryName(cat)
		fmt.Printf("\033[1;36m%s\033[0m (%d)\n", catName, len(catTools))
		fmt.Printf("%-16s %-6s %-12s %s\n", "NAME", "STATUS", "INSTALL VIA", "DESCRIPTION")
		fmt.Printf("%-16s %-6s %-12s %s\n", "────────────────", "──────", "────────────", "─────────────────────────────────────")

		for _, tool := range catTools {
			printToolRow(tool, false)
		}
		fmt.Println()
	}
}

func runListFlat(toolList []*tools.Tool) {
	// Sort by name
	sort.Slice(toolList, func(i, j int) bool {
		return toolList[i].Name < toolList[j].Name
	})

	// Header
	fmt.Printf("%-16s %-9s %-6s %-12s %s\n", "NAME", "CATEGORY", "STATUS", "INSTALL VIA", "DESCRIPTION")
	fmt.Printf("%-16s %-9s %-6s %-12s %s\n", "────────────────", "─────────", "──────", "────────────", "─────────────────────────────────────")

	for _, tool := range toolList {
		printToolRow(tool, true)
	}
}

func printToolRow(tool *tools.Tool, showCategory bool) {
	var status string
	if tool.IsInstalled() {
		status = "\033[32m✓\033[0m"
	} else {
		status = "\033[31m✗\033[0m"
	}

	// Get install methods
	methods := tool.GetAvailableMethods()
	var methodStr string
	var methodLen int
	if len(methods) > 0 {
		methodStrs := make([]string, len(methods))
		for i, m := range methods {
			methodStrs[i] = string(m)
		}
		methodStr = strings.Join(methodStrs, ",")
		if len(methodStr) > 12 {
			methodStr = methodStr[:11] + "…"
		}
		methodLen = len(methodStr)
	} else {
		methodStr = "\033[33m-\033[0m"
		methodLen = 1
	}

	// Truncate description if too long
	desc := tool.Description
	if len(desc) > 40 {
		desc = desc[:37] + "..."
	}

	// Calculate padding for method column (12 chars wide)
	paddingLen := 12 - methodLen
	if paddingLen < 0 {
		paddingLen = 0
	}
	methodPadding := strings.Repeat(" ", paddingLen)

	// Print with fixed widths
	if showCategory {
		fmt.Printf("%-16s %-9s %s      %s%s %s\n", tool.Name, tool.Category, status, methodStr, methodPadding, desc)
	} else {
		fmt.Printf("%-16s %s      %s%s %s\n", tool.Name, status, methodStr, methodPadding, desc)
	}
}
