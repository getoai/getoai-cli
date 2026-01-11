package cli

import (
	"fmt"

	"github.com/getoai/getoai-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage getoai configuration",
	Long: `View and manage getoai configuration settings.

Configuration includes:
  - Proxy settings (HTTP/HTTPS)
  - Package manager mirrors (npm, pip, go)
  - Default installation preferences`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run:   runConfigShow,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Available keys:
  http_proxy    - HTTP proxy URL
  https_proxy   - HTTPS proxy URL
  npm_registry  - npm registry URL (e.g., https://registry.npmmirror.com)
  pypi_mirror   - PyPI mirror URL (e.g., https://pypi.tuna.tsinghua.edu.cn/simple)
  go_proxy      - Go module proxy (e.g., https://goproxy.cn,direct)

Examples:
  getoai config set npm_registry https://registry.npmmirror.com
  getoai config set go_proxy https://goproxy.cn,direct`,
	Args: cobra.ExactArgs(2),
	Run:  runConfigSet,
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	Run:   runConfigPath,
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configPathCmd)
	rootCmd.AddCommand(configCmd)
}

func runConfigShow(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		printError(fmt.Sprintf("Failed to load config: %v", err))
		return
	}

	fmt.Println("\nCurrent Configuration:")
	fmt.Println("----------------------")

	if cfg.HttpProxy != "" {
		fmt.Printf("http_proxy:    %s\n", cfg.HttpProxy)
	}
	if cfg.HttpsProxy != "" {
		fmt.Printf("https_proxy:   %s\n", cfg.HttpsProxy)
	}
	if cfg.NpmRegistry != "" {
		fmt.Printf("npm_registry:  %s\n", cfg.NpmRegistry)
	}
	if cfg.PypiMirror != "" {
		fmt.Printf("pypi_mirror:   %s\n", cfg.PypiMirror)
	}
	if cfg.GoProxy != "" {
		fmt.Printf("go_proxy:      %s\n", cfg.GoProxy)
	}
	if cfg.BinPath != "" {
		fmt.Printf("bin_path:      %s\n", cfg.BinPath)
	}

	if cfg.HttpProxy == "" && cfg.HttpsProxy == "" && cfg.NpmRegistry == "" &&
		cfg.PypiMirror == "" && cfg.GoProxy == "" && cfg.BinPath == "" {
		fmt.Println("(No custom configuration set)")
	}

	fmt.Println()
}

func runConfigSet(cmd *cobra.Command, args []string) {
	key := args[0]
	value := args[1]

	cfg, err := config.Load()
	if err != nil {
		printError(fmt.Sprintf("Failed to load config: %v", err))
		return
	}

	switch key {
	case "http_proxy":
		cfg.HttpProxy = value
	case "https_proxy":
		cfg.HttpsProxy = value
	case "npm_registry":
		cfg.NpmRegistry = value
	case "pypi_mirror":
		cfg.PypiMirror = value
	case "go_proxy":
		cfg.GoProxy = value
	case "bin_path":
		cfg.BinPath = value
	default:
		printError(fmt.Sprintf("Unknown config key: %s", key))
		fmt.Println("Available keys: http_proxy, https_proxy, npm_registry, pypi_mirror, go_proxy, bin_path")
		return
	}

	if err := config.Save(cfg); err != nil {
		printError(fmt.Sprintf("Failed to save config: %v", err))
		return
	}

	printSuccess(fmt.Sprintf("Set %s = %s", key, value))
}

func runConfigPath(cmd *cobra.Command, args []string) {
	fmt.Println(config.GetConfigPath())
}
