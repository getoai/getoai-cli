package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/getoai/getoai-cli/internal/installer"
	"github.com/getoai/getoai-cli/internal/tools"
	"github.com/getoai/getoai-cli/internal/util"
)

var uninstallCmd = &cobra.Command{
	Use:     "uninstall <tool> [tools...]",
	Aliases: []string{"remove", "rm"},
	Short:   "Uninstall AI tools",
	Long: `Uninstall one or more AI tools from your system.

Examples:
  getoai uninstall chatgpt-cli
  getoai uninstall aider llm
  getoai rm ollama`,
	Args: cobra.MinimumNArgs(1),
	Run:  runUninstall,
}

var forceUninstall bool

func init() {
	uninstallCmd.Flags().BoolVarP(&forceUninstall, "force", "f", false, "Skip confirmation prompt")
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) {
	for _, toolName := range args {
		uninstallTool(toolName)
	}
}

func uninstallTool(name string) {
	tool, ok := tools.Get(name)
	if !ok {
		printError(fmt.Sprintf("Unknown tool: %s", name))
		return
	}

	if !tool.IsInstalled() {
		printInfo(fmt.Sprintf("%s is not installed", name))
		return
	}

	if !forceUninstall {
		fmt.Printf("Are you sure you want to uninstall %s? [y/N] ", name)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			printInfo("Uninstall canceled")
			return
		}
	}

	spinner := util.NewSpinner(fmt.Sprintf("Uninstalling %s...", name))
	spinner.Start()

	var uninstallErr error

	// Handle docker-compose installations specially
	if tool.IsDockerComposeInstall() {
		installDir := tool.GetComposeInstallDir()
		if installDir != "" {
			spinner.Stop()
			dockerInst := installer.NewDockerInstaller()
			uninstallErr = dockerInst.UninstallCompose(installDir)
			if uninstallErr == nil {
				printSuccess(fmt.Sprintf("%s stopped successfully", name))
				return
			}
		}
	} else {
		// Try to find the appropriate uninstaller for non-compose tools
		for method := range tool.InstallMethods {
			inst, err := installer.GetInstaller(method)
			if err != nil {
				continue
			}
			if err := inst.Uninstall(getUninstallPackage(tool, method)); err == nil {
				uninstallErr = nil
				break
			} else {
				uninstallErr = err
			}
		}
	}

	if uninstallErr != nil {
		spinner.Error(fmt.Sprintf("Failed to uninstall %s: %v", name, uninstallErr))
		fmt.Printf("  You may need to manually uninstall from %s\n", tool.Website)
		return
	}

	// Verify uninstallation
	if !tool.IsInstalled() {
		spinner.Success(fmt.Sprintf("%s uninstalled successfully", name))
	} else {
		spinner.Info(fmt.Sprintf("%s uninstall completed, but command still found in PATH", name))
		fmt.Println("  You may need to restart your shell")
	}
}

func getUninstallPackage(tool *tools.Tool, method installer.InstallMethod) string {
	if config, ok := tool.InstallMethods[method]; ok {
		return config.Package
	}
	return tool.Name
}
