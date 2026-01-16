package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PromptChoice displays a numbered menu and returns the selected index
func PromptChoice(title string, options []string, descriptions []string) (int, error) {
	fmt.Println()
	fmt.Printf("\033[1m%s\033[0m\n", title)
	fmt.Println()

	// Display options
	for i, option := range options {
		desc := ""
		if i < len(descriptions) && descriptions[i] != "" {
			desc = fmt.Sprintf(" - %s", descriptions[i])
		}
		fmt.Printf("  \033[36m%d)\033[0m %s%s\n", i+1, option, desc)
	}

	fmt.Println()
	fmt.Printf("Enter your choice (1-%d): ", len(options))

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Parse input
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(options) {
		return -1, fmt.Errorf("invalid choice: %s", input)
	}

	return choice - 1, nil
}
