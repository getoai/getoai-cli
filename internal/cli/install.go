package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/getoai/getoai-cli/internal/config"
	"github.com/getoai/getoai-cli/internal/installer"
	"github.com/getoai/getoai-cli/internal/platform"
	"github.com/getoai/getoai-cli/internal/tools"
	"github.com/getoai/getoai-cli/internal/util"
)

var installCmd = &cobra.Command{
	Use:   "install <tool> [tools...]",
	Short: "Install AI tools",
	Long: `Install one or more AI tools.

When multiple installation methods are available, you'll be prompted
to choose your preferred method. Use --method to skip the prompt.

Examples:
  getoai install ollama
  getoai install claude-code aider
  getoai install ollama --method brew
  getoai install ollama --method docker`,
	Args: cobra.MinimumNArgs(1),
	Run:  runInstall,
}

var installMethod string
var skipDepsCheck bool

func init() {
	installCmd.Flags().StringVarP(&installMethod, "method", "m", "", "Installation method (brew, npm, pip, script, go, docker)")
	installCmd.Flags().BoolVar(&skipDepsCheck, "skip-deps", false, "Skip dependency check")
}

func runInstall(cmd *cobra.Command, args []string) {
	// Apply config (proxy settings, etc.)
	cfg, _ := config.Load()
	if cfg != nil {
		cfg.ApplyEnv()
	}

	for _, toolName := range args {
		installTool(toolName)
		fmt.Println()
	}
}

// promptMethodSelection shows an interactive menu for selecting install method
func promptMethodSelection(toolName string, availableMethods []installer.InstallMethod) (installer.InstallMethod, error) {
	if len(availableMethods) == 0 {
		return "", fmt.Errorf("no methods available")
	}

	if len(availableMethods) == 1 {
		return availableMethods[0], nil
	}

	// Build menu options
	options := make([]string, len(availableMethods))
	descriptions := make([]string, len(availableMethods))

	for i, method := range availableMethods {
		options[i] = string(method)
		descriptions[i] = installer.GetMethodDescription(method)
	}

	title := fmt.Sprintf("Multiple installation methods available for %s", toolName)
	choice, err := util.PromptChoice(title, options, descriptions)

	if err != nil {
		return "", fmt.Errorf("selection cancelled or invalid: %w", err)
	}

	return availableMethods[choice], nil
}

func installTool(name string) {
	tool, ok := tools.Get(name)
	if !ok {
		printError(fmt.Sprintf("Unknown tool: %s", name))
		suggestSimilar(name)
		return
	}

	if tool.IsInstalled() {
		printInfo(fmt.Sprintf("%s is already installed (version: %s)", name, tool.GetVersion()))
		return
	}

	availableMethods := tool.GetAvailableMethods()
	if len(availableMethods) == 0 {
		// Check if there are missing dependencies that can be installed
		missingDeps := getMissingDependencies(tool)
		if len(missingDeps) > 0 {
			// Ask user if they want to install dependencies
			if promptInstallDependencies(missingDeps) {
				// Install dependencies
				for _, dep := range missingDeps {
					fmt.Printf("\n")
					installTool(dep)
				}
				// Refresh platform detection after installing dependencies
				platform.Refresh()
				// Retry getting available methods
				availableMethods = tool.GetAvailableMethods()
				if len(availableMethods) > 0 {
					fmt.Printf("\nDependencies installed. Continuing with %s installation...\n\n", name)
				}
			}
		}

		// Check again after potential dependency installation
		if len(availableMethods) == 0 {
			printError(fmt.Sprintf("No installation method available for %s on this system", name))
			showMissingDependencies(tool)
			fmt.Printf("  Visit %s for manual installation\n", tool.Website)
			return
		}
	}

	var method installer.InstallMethod
	var err error

	// If --method flag is specified, use it (backward compatibility)
	if installMethod != "" {
		method = installer.InstallMethod(installMethod)
		found := false
		for _, m := range availableMethods {
			if m == method {
				found = true
				break
			}
		}
		if !found {
			printError(fmt.Sprintf("Method '%s' not available for %s", installMethod, name))
			fmt.Printf("  Available methods: %v\n", availableMethods)
			return
		}
	} else {
		// If multiple methods available, show interactive menu
		if len(availableMethods) > 1 {
			method, err = promptMethodSelection(name, availableMethods)
			if err != nil {
				printError(fmt.Sprintf("Method selection failed: %v", err))
				return
			}
			fmt.Printf("\nSelected installation method: \033[32m%s\033[0m\n\n", method)
		} else {
			// Single method available, use it automatically
			method = availableMethods[0]
		}
	}

	// Check dependencies
	if !skipDepsCheck {
		checkDependencies(method)
	}

	spinner := util.NewSpinner(fmt.Sprintf("Installing %s using %s...", name, method))
	spinner.Start()

	if err := tool.Install(method); err != nil {
		spinner.Error(fmt.Sprintf("Failed to install %s: %v", name, err))
		return
	}

	// Verify installation
	if tool.IsInstalled() {
		version := tool.GetVersion()
		if version == "N/A" || version == "not installed" {
			// Desktop app without version info
			spinner.Success(fmt.Sprintf("%s installed successfully!", name))
		} else {
			spinner.Success(fmt.Sprintf("%s installed successfully! (version: %s)", name, version))
		}
	} else {
		// For desktop apps (with AppName), show different message
		if tool.AppName != "" {
			spinner.Info(fmt.Sprintf("%s installation completed", name))
			fmt.Println("  Desktop app installed, you may need to restart Finder or reboot to see it")
		} else {
			spinner.Info(fmt.Sprintf("%s installation completed, but command not found in PATH", name))
			showPathHint(method)
		}
	}
}

func suggestSimilar(name string) {
	results := tools.Search(name)
	if len(results) > 0 && len(results) <= 5 {
		names := make([]string, len(results))
		for i, t := range results {
			names[i] = t.Name
		}
		fmt.Printf("  Did you mean: %s?\n", strings.Join(names, ", "))
	}
	fmt.Println("  Use 'getoai list' to see all available tools")
	fmt.Println("  Use 'getoai search <keyword>' to search for tools")
}

func checkDependencies(method installer.InstallMethod) {
	p := platform.Detect()

	switch method {
	case installer.MethodNpm:
		if !p.HasNpm {
			printInfo("npm not found. Install Node.js first:")
			fmt.Println("  macOS:   brew install node")
			fmt.Println("  Linux:   apt install nodejs npm / dnf install nodejs npm")
			fmt.Println("  Windows: winget install OpenJS.NodeJS")
		}
	case installer.MethodPip:
		if !p.HasPip && !p.HasPip3 {
			printInfo("pip not found. Install Python first:")
			fmt.Println("  macOS:   brew install python")
			fmt.Println("  Linux:   apt install python3-pip / dnf install python3-pip")
			fmt.Println("  Windows: winget install Python.Python.3")
		}
	case installer.MethodGo:
		if !p.HasGo {
			printInfo("go not found. Install Go first:")
			fmt.Println("  macOS:   brew install go")
			fmt.Println("  Linux:   apt install golang / dnf install golang")
			fmt.Println("  Windows: winget install GoLang.Go")
		}
	case installer.MethodDocker:
		if !p.HasDocker {
			printInfo("docker not found. Install Docker first:")
			fmt.Println("  Visit: https://docs.docker.com/get-docker/")
		}
	case installer.MethodBrew:
		if !p.HasBrew {
			printInfo("Homebrew not found. Install Homebrew first:")
			fmt.Println("  /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
		}
	}
}

func showMissingDependencies(tool *tools.Tool) {
	p := platform.Detect()
	var missing []string

	for method := range tool.InstallMethods {
		switch method {
		case installer.MethodNpm:
			if !p.HasNpm {
				missing = append(missing, "npm (install Node.js)")
			}
		case installer.MethodPip:
			if !p.HasPip && !p.HasPip3 {
				missing = append(missing, "pip (install Python)")
			}
		case installer.MethodGo:
			if !p.HasGo {
				missing = append(missing, "go (install Go)")
			}
		case installer.MethodDocker:
			if !p.HasDocker {
				missing = append(missing, "docker")
			}
		case installer.MethodBrew:
			if !p.HasBrew {
				missing = append(missing, "brew (Homebrew)")
			}
		}
	}

	if len(missing) > 0 {
		fmt.Printf("  Missing dependencies: %s\n", strings.Join(missing, ", "))
	}
}

func showPathHint(method installer.InstallMethod) {
	switch method {
	case installer.MethodGo:
		fmt.Println("  Add ~/go/bin to your PATH:")
		fmt.Println("    export PATH=$PATH:~/go/bin")
	case installer.MethodPip:
		fmt.Println("  The binary might be in ~/.local/bin")
		fmt.Println("  Add it to your PATH if needed:")
		fmt.Println("    export PATH=$PATH:~/.local/bin")
	case installer.MethodNpm:
		fmt.Println("  You may need to restart your shell")
	default:
		fmt.Println("  You may need to restart your shell or add the binary to your PATH")
	}
}

// dependencyMap maps install methods to the tool names that provide them
var dependencyMap = map[installer.InstallMethod]struct {
	toolName string
	desc     string
}{
	installer.MethodNpm:    {toolName: "node", desc: "Node.js"},
	installer.MethodPip:    {toolName: "", desc: "Python"}, // Python installation is complex, skip auto-install
	installer.MethodGo:     {toolName: "", desc: "Go"},     // Go installation is complex, skip auto-install
	installer.MethodDocker: {toolName: "docker", desc: "Docker"},
}

// getMissingDependencies returns a list of tool names that can be installed to satisfy missing dependencies
func getMissingDependencies(tool *tools.Tool) []string {
	p := platform.Detect()
	var deps []string
	seen := make(map[string]bool)

	for method := range tool.InstallMethods {
		var missing bool
		var depInfo struct {
			toolName string
			desc     string
		}

		switch method {
		case installer.MethodNpm:
			if !p.HasNpm {
				missing = true
				depInfo = dependencyMap[method]
			}
		case installer.MethodDocker:
			if !p.HasDocker {
				missing = true
				depInfo = dependencyMap[method]
			}
		}

		if missing && depInfo.toolName != "" && !seen[depInfo.toolName] {
			// Check if the dependency tool is available
			if depTool, ok := tools.Get(depInfo.toolName); ok {
				if len(depTool.GetAvailableMethods()) > 0 {
					deps = append(deps, depInfo.toolName)
					seen[depInfo.toolName] = true
				}
			}
		}
	}

	return deps
}

// promptInstallDependencies asks the user if they want to install missing dependencies
func promptInstallDependencies(deps []string) bool {
	if len(deps) == 0 {
		return false
	}

	var depDescs []string
	for _, dep := range deps {
		if info, ok := dependencyMap[getMethodForDep(dep)]; ok {
			depDescs = append(depDescs, info.desc)
		} else {
			depDescs = append(depDescs, dep)
		}
	}

	fmt.Printf("\n\033[33m!\033[0m Missing dependencies: %s\n", strings.Join(depDescs, ", "))
	fmt.Printf("  Install required dependencies first? [Y/n] ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	return response == "" || response == "y" || response == "yes"
}

// getMethodForDep returns the install method that a dependency tool provides
func getMethodForDep(toolName string) installer.InstallMethod {
	switch toolName {
	case "node":
		return installer.MethodNpm
	case "docker":
		return installer.MethodDocker
	default:
		return ""
	}
}
