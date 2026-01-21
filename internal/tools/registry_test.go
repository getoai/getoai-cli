package tools

import (
	"testing"

	"github.com/getoai/getoai-cli/internal/installer"
)

func TestAllToolsHaveValidConfiguration(t *testing.T) {
	all := List()

	for _, tool := range all {
		t.Run(tool.Name, func(t *testing.T) {
			// Test 1: Name should not be empty
			if tool.Name == "" {
				t.Error("Tool name is empty")
			}

			// Test 2: Description should not be empty
			if tool.Description == "" {
				t.Errorf("Tool %s has no description", tool.Name)
			}

			// Test 3: Should have at least one installation method
			if len(tool.InstallMethods) == 0 {
				t.Errorf("Tool %s has no installation methods", tool.Name)
			}

			// Test 4: Category should be valid
			validCategories := map[Category]bool{
				CategoryLLM:      true,
				CategoryCoding:   true,
				CategoryUI:       true,
				CategoryUtility:  true,
				CategoryPlatform: true,
				CategoryInfra:    true,
			}
			if !validCategories[tool.Category] {
				t.Errorf("Tool %s has invalid category: %s", tool.Name, tool.Category)
			}

			// Test 5: Website should not be empty
			if tool.Website == "" {
				t.Errorf("Tool %s has no website", tool.Name)
			}

			// Test 6: Verify install method configurations
			for method, config := range tool.InstallMethods {
				if config.Package == "" && len(config.DownloadURLs) == 0 {
					t.Errorf("Tool %s has %s method but no package or download URLs", tool.Name, method)
				}

				// Test 7: Download method should have URLs or package page
				if method == installer.MethodDownload {
					if config.Package == "" && len(config.DownloadURLs) == 0 {
						t.Errorf("Tool %s download method has neither package URL nor download URLs", tool.Name)
					}
				}
			}
		})
	}
}

func TestWindowsToolsHaveCorrectMethods(t *testing.T) {
	windowsTools := []string{
		"ollama", "node", "gh", "docker", "vscode", "docker-compose",
	}

	for _, toolName := range windowsTools {
		t.Run(toolName, func(t *testing.T) {
			tool, ok := Get(toolName)
			if !ok {
				t.Fatalf("Tool %s not found", toolName)
			}

			// Check if tool has Chocolatey or Scoop support
			hasChoco := false
			hasScoop := false

			for method := range tool.InstallMethods {
				if method == installer.MethodChoco {
					hasChoco = true
				}
				if method == installer.MethodScoop {
					hasScoop = true
				}
			}

			if !hasChoco && !hasScoop && toolName != "docker-compose" {
				t.Errorf("Windows tool %s should have Chocolatey or Scoop support", toolName)
			}

			// docker-compose only supports Chocolatey
			if toolName == "docker-compose" && !hasChoco {
				t.Error("docker-compose should have Chocolatey support")
			}
		})
	}
}

func TestDesktopAppsHaveDownloadMethod(t *testing.T) {
	desktopApps := []string{
		"cursor", "lmstudio", "jan", "msty", "cherry-studio",
		"chatbox", "vscode",
	}

	for _, appName := range desktopApps {
		t.Run(appName, func(t *testing.T) {
			tool, ok := Get(appName)
			if !ok {
				t.Fatalf("Desktop app %s not found", appName)
			}

			// Should have download method
			config, hasDownload := tool.InstallMethods[installer.MethodDownload]
			if !hasDownload {
				t.Errorf("Desktop app %s should have download method", appName)
			}

			// Download method should have URLs for at least one platform
			if len(config.DownloadURLs) == 0 && config.Package == "" {
				t.Errorf("Desktop app %s download method has no URLs", appName)
			}

			// Should have AppName for desktop apps
			if tool.AppName == "" {
				t.Errorf("Desktop app %s should have AppName set", appName)
			}
		})
	}
}

func TestPlatformOverrides(t *testing.T) {
	tests := []struct {
		tool     string
		platform string
		want     installer.InstallMethod
	}{
		{tool: "ollama", platform: "windows", want: installer.MethodChoco},
		{tool: "node", platform: "windows", want: installer.MethodChoco},
		{tool: "gh", platform: "windows", want: installer.MethodChoco},
		{tool: "docker", platform: "windows", want: installer.MethodChoco},
	}

	for _, tt := range tests {
		t.Run(tt.tool+"_"+tt.platform, func(t *testing.T) {
			tool, ok := Get(tt.tool)
			if !ok {
				t.Fatalf("Tool %s not found", tt.tool)
			}

			platformConfig, hasPlatformOverride := tool.PlatformOverrides[tt.platform]
			if !hasPlatformOverride {
				t.Errorf("Tool %s should have platform override for %s", tt.tool, tt.platform)
				return
			}

			if _, hasMethod := platformConfig[tt.want]; !hasMethod {
				t.Errorf("Tool %s platform override for %s should include %s method", tt.tool, tt.platform, tt.want)
			}
		})
	}
}

func TestToolSearch(t *testing.T) {
	tests := []struct {
		query string
		want  []string // Expected to be in results
	}{
		{
			query: "ollama",
			want:  []string{"ollama"},
		},
		{
			query: "docker",
			want:  []string{"docker", "docker-compose"},
		},
		{
			query: "code",
			want:  []string{"vscode", "claude-code"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			results := Search(tt.query)

			// Check if all expected tools are in results
			resultMap := make(map[string]bool)
			for _, tool := range results {
				resultMap[tool.Name] = true
			}

			for _, expectedTool := range tt.want {
				if !resultMap[expectedTool] {
					// Build a readable list of tool names for error message
					toolNames := make([]string, len(results))
					for i, tool := range results {
						toolNames[i] = tool.Name
					}
					t.Errorf("Search(%q) should include %s, got: %v", tt.query, expectedTool, toolNames)
				}
			}
		})
	}
}

func TestNoHardcodedFileTypes(t *testing.T) {
	all := List()

	for _, tool := range all {
		t.Run(tool.Name, func(t *testing.T) {
			for method, config := range tool.InstallMethods {
				if method == installer.MethodDownload {
					// FileType should not be set - we use auto-detection
					if config.FileType != "" {
						t.Errorf("Tool %s has hardcoded FileType: %s (should use auto-detection)", tool.Name, config.FileType)
					}
				}
			}
		})
	}
}

func TestDownloadURLsFormat(t *testing.T) {
	all := List()

	for _, tool := range all {
		t.Run(tool.Name, func(t *testing.T) {
			config, hasDownload := tool.InstallMethods[installer.MethodDownload]
			if !hasDownload {
				return // Skip non-download tools
			}

			for platform, url := range config.DownloadURLs {
				// URLs should start with https://
				if url != "" && url[:8] != "https://" && url[:7] != "http://" {
					t.Errorf("Tool %s %s download URL should be https:// or http://, got: %s", tool.Name, platform, url)
				}

				// URL should not have obvious errors
				if url == "https://" || url == "http://" {
					t.Errorf("Tool %s %s download URL is incomplete: %s", tool.Name, platform, url)
				}
			}
		})
	}
}
