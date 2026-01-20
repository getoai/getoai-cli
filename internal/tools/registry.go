package tools

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/getoai/getoai-cli/internal/installer"
	"github.com/getoai/getoai-cli/internal/platform"
)

type Category string

const (
	CategoryLLM      Category = "llm"
	CategoryCoding   Category = "coding"
	CategoryUI       Category = "ui"
	CategoryUtility  Category = "utility"
	CategoryPlatform Category = "platform"
	CategoryInfra    Category = "infra"
)

type Tool struct {
	Name        string
	Description string
	Category    Category
	Website     string
	Command     string // command to check if installed
	AppName     string // for desktop apps: app name (e.g., "chatbox.app" on macOS)

	// Installation options by method
	InstallMethods map[installer.InstallMethod]InstallConfig

	// Platform-specific overrides
	PlatformOverrides map[string]map[installer.InstallMethod]InstallConfig
}

type InstallConfig struct {
	Package string   // package name or URL
	Args    []string // additional arguments

	// Docker-specific options
	DockerPorts   []string          // port mappings, e.g. ["3000:3000", "8080:80"]
	DockerEnv     map[string]string // environment variables
	DockerVolumes []string          // volume mappings
	DockerName    string            // container name
	DockerCompose string            // docker-compose repo URL (for complex apps)

	// Download-specific options (for desktop apps)
	DownloadURLs map[string]string // platform-specific download URLs: "darwin", "linux", "windows"
	FileType     string            // file type: "dmg", "pkg", "deb", "appimage", "exe", "msi"
}

var registry = map[string]*Tool{}

func init() {
	// Register all tools
	registerOllama()
	registerClaudeCode()
	registerOpenAICLI()
	registerAider()
	registerLLM()
	registerOpenWebUI()
	registerChatGPTCLI()
	registerCursor()
	registerLMStudio()
	// CLI tools
	registerGPTEngineer()
	registerAutoGPT()
	registerFabric()
	registerShellGPT()
	registerMods()
	registerTGPT()
	registerGlow()
	registerGitHubCopilotCLI()
	registerTabby()
	registerLocalAI()
	registerAnythingLLM()
	registerJan()
	registerMSDD()
	// Self-hosted platforms
	registerOneAPI()
	registerNewAPI()
	registerDify()
	registerFastGPT()
	registerLobeChat()
	registerChatGPTNextWeb()
	registerFlowise()
	registerLangflow()
	registerQuivr()
	registerPrivateGPT()
	registerLibreChat()
	registerMaxKB()
	registerRAGFlow()
	registerDBGPT()
	registerChatGLM()
	registerChatWoot()
	// AI Infra
	registerVLLM()
	registerTextGenWebUI()
	registerComfyUI()
	registerSDWebUI()
	registerKoboldCpp()
	registerLlamaCpp()
	registerXinference()
	registerSGLang()
	// Development tools
	registerNvm()
	registerNode()
	registerDocker()
	registerDockerCompose()
	registerGitHubCLI()
	// Desktop AI Apps
	registerCherryStudio()
	registerChatbox()
	registerTypingMind()
	// Additional CLI tools
	registerCodexCLI()
	registerGeminiCLI()
	registerAIChat()
	registerGoingChat()
	registerOpenInterpreter()
	// More AI Coding Tools
	registerOpenCode()
	registerWindsurf()
	registerTabnine()
	registerSupermaven()
	registerPieces()
	registerCody()
	registerQodo()
	registerReplit()
	// Developer Tools (Open Source / Free Community Edition)
	registerVSCode()
	registerIntelliJIDEA()
	registerPyCharm()
	// Terminal Tools (Open Source)
	registerWarp()
	registerITerm2()
	registerAlacritty()
	registerKitty()
	// API Tools (Open Source / Free)
	registerPostman()
	registerInsomnia()
	// Database Tools (Open Source / Free)
	registerTablePlus()
	registerDBeaver()
	// Productivity Tools (Free / Free for Personal Use)
	registerRaycast()
	registerOrbStack()
	registerFig()
}

func registerOllama() {
	Register(&Tool{
		Name:        "ollama",
		Description: "Run large language models locally",
		Category:    CategoryLLM,
		Website:     "https://ollama.ai",
		Command:     "ollama",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDownload: {Package: "https://ollama.ai/download"},
			installer.MethodBrew:     {Package: "ollama"},
			installer.MethodScript:   {Package: "https://ollama.ai/install.sh"},
			installer.MethodChoco:    {Package: "ollama"},
			installer.MethodScoop:    {Package: "ollama"},
		},
		PlatformOverrides: map[string]map[installer.InstallMethod]InstallConfig{
			"darwin": {
				installer.MethodDownload: {Package: "https://ollama.ai/download"},
			},
			"linux": {
				installer.MethodScript: {Package: "https://ollama.ai/install.sh"},
			},
			"windows": {
				installer.MethodChoco: {Package: "ollama"},
			},
		},
	})
}

func registerClaudeCode() {
	Register(&Tool{
		Name:        "claude-code",
		Description: "Claude AI coding assistant CLI",
		Category:    CategoryCoding,
		Website:     "https://claude.ai",
		Command:     "claude",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodNpm: {Package: "@anthropic-ai/claude-code"},
		},
	})
}

func registerOpenAICLI() {
	Register(&Tool{
		Name:        "openai-cli",
		Description: "OpenAI official command-line interface",
		Category:    CategoryUtility,
		Website:     "https://platform.openai.com",
		Command:     "openai",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip: {Package: "openai"},
		},
	})
}

func registerAider() {
	Register(&Tool{
		Name:        "aider",
		Description: "AI pair programming in your terminal",
		Category:    CategoryCoding,
		Website:     "https://aider.chat",
		Command:     "aider",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip:  {Package: "aider-chat"},
			installer.MethodBrew: {Package: "aider"},
		},
	})
}

func registerLLM() {
	Register(&Tool{
		Name:        "llm",
		Description: "Access LLMs from the command line by Simon Willison",
		Category:    CategoryUtility,
		Website:     "https://llm.datasette.io",
		Command:     "llm",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip:  {Package: "llm"},
			installer.MethodBrew: {Package: "llm"},
		},
	})
}

func registerOpenWebUI() {
	Register(&Tool{
		Name:        "open-webui",
		Description: "User-friendly WebUI for LLMs (Ollama compatible)",
		Category:    CategoryUI,
		Website:     "https://openwebui.com",
		Command:     "open-webui",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip: {Package: "open-webui"},
			installer.MethodDocker: {
				Package:     "ghcr.io/open-webui/open-webui:main",
				DockerName:  "open-webui",
				DockerPorts: []string{"3000:8080"},
				DockerVolumes: []string{
					"open-webui-data:/app/backend/data",
				},
			},
		},
	})
}

func registerChatGPTCLI() {
	Register(&Tool{
		Name:        "chatgpt-cli",
		Description: "ChatGPT in your terminal",
		Category:    CategoryUtility,
		Website:     "https://github.com/kardolus/chatgpt-cli",
		Command:     "chatgpt",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodGo:   {Package: "github.com/kardolus/chatgpt-cli/cmd/chatgpt"},
			installer.MethodBrew: {Package: "kardolus/chatgpt-cli/chatgpt-cli"},
		},
	})
}

func registerCursor() {
	Register(&Tool{
		Name:        "cursor",
		Description: "AI-first code editor built on VS Code",
		Category:    CategoryCoding,
		Website:     "https://cursor.sh",
		Command:     "cursor",
		AppName:     "Cursor.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "cursor", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://cursor.sh",
				DownloadURLs: map[string]string{
					"darwin":  "https://downloader.cursor.sh/mac/universal",
					"linux":   "https://downloader.cursor.sh/linux/appImage/x64",
					"windows": "https://downloader.cursor.sh/windows/nsis/x64",
				},
				},
		},
	})
}
func registerLMStudio() {
	Register(&Tool{
		Name:        "lmstudio",
		Description: "Discover, download, and run local LLMs",
		Category:    CategoryUI,
		Website:     "https://lmstudio.ai",
		Command:     "",
		AppName:     "LM Studio.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "lm-studio", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://lmstudio.ai",
				DownloadURLs: map[string]string{
					"darwin":  "https://releases.lmstudio.ai/darwin/arm64/latest",
					"linux":   "https://releases.lmstudio.ai/linux/x86/latest",
					"windows": "https://releases.lmstudio.ai/windows/x86/latest",
				},
				},
		},
	})
}

func Register(tool *Tool) {
	registry[tool.Name] = tool
}

func Get(name string) (*Tool, bool) {
	tool, ok := registry[name]
	return tool, ok
}

func List() []*Tool {
	tools := make([]*Tool, 0, len(registry))
	for _, tool := range registry {
		tools = append(tools, tool)
	}
	return tools
}

func ListByCategory(cat Category) []*Tool {
	var tools []*Tool
	for _, tool := range registry {
		if tool.Category == cat {
			tools = append(tools, tool)
		}
	}
	return tools
}

func (t *Tool) IsInstalled() bool {
	// Check if installed via docker-compose (check if containers are running)
	if t.IsDockerComposeInstall() {
		installDir := t.GetComposeInstallDir()
		if installDir == "" {
			return false
		}
		// Check if containers are actually running
		return t.IsComposeRunning(installDir)
	}

	// Check if installed via docker container (single container)
	if t.IsDockerContainerInstalled() {
		return true
	}

	// Check desktop apps (by AppName)
	if t.AppName != "" {
		return t.IsDesktopAppInstalled()
	}

	// Check command in PATH
	if t.Command == "" {
		return false
	}
	return installer.CheckInstalled(t.Command)
}

// IsDesktopAppInstalled checks if a desktop app is installed
func (t *Tool) IsDesktopAppInstalled() bool {
	p := platform.Detect()

	switch p.OS {
	case "darwin":
		// Check /Applications
		appPath := fmt.Sprintf("/Applications/%s", t.AppName)
		if _, err := os.Stat(appPath); err == nil {
			return true
		}
		// Also check ~/Applications
		homeDir, _ := os.UserHomeDir()
		userAppPath := fmt.Sprintf("%s/Applications/%s", homeDir, t.AppName)
		if _, err := os.Stat(userAppPath); err == nil {
			return true
		}
		return false

	case "linux":
		// Check common installation locations
		locations := []string{
			fmt.Sprintf("/usr/share/applications/%s.desktop", t.Name),
			fmt.Sprintf("/usr/local/share/applications/%s.desktop", t.Name),
		}
		homeDir, _ := os.UserHomeDir()
		if homeDir != "" {
			locations = append(locations,
				fmt.Sprintf("%s/.local/share/applications/%s.desktop", homeDir, t.Name),
				fmt.Sprintf("%s/.local/bin/%s.appimage", homeDir, t.Name),
			)
		}
		for _, loc := range locations {
			if _, err := os.Stat(loc); err == nil {
				return true
			}
		}
		return false

	case "windows":
		// Check Program Files
		programFiles := os.Getenv("ProgramFiles")
		if programFiles != "" {
			appPath := fmt.Sprintf("%s\\%s", programFiles, t.AppName)
			if _, err := os.Stat(appPath); err == nil {
				return true
			}
		}
		// Check Program Files (x86)
		programFilesX86 := os.Getenv("ProgramFiles(x86)")
		if programFilesX86 != "" {
			appPath := fmt.Sprintf("%s\\%s", programFilesX86, t.AppName)
			if _, err := os.Stat(appPath); err == nil {
				return true
			}
		}
		return false
	}

	return false
}

// IsComposeRunning checks if docker-compose containers are running
func (t *Tool) IsComposeRunning(installDir string) bool {
	// Find compose file
	locations := []string{
		"docker/docker-compose.yaml",
		"docker/docker-compose.yml",
		"docker-compose.yaml",
		"docker-compose.yml",
	}

	var composeFile string
	for _, loc := range locations {
		path := fmt.Sprintf("%s/%s", installDir, loc)
		if _, err := os.Stat(path); err == nil {
			composeFile = path
			break
		}
	}

	if composeFile == "" {
		return false
	}

	// Check running containers using --status=running filter
	out, err := installer.RunCommandSilent("docker", "compose", "-f", composeFile, "ps", "--status=running", "-q")
	if err != nil {
		return false
	}

	// Filter out warning lines (only keep container IDs which are hex strings)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Container IDs are 12 or 64 character hex strings
		if len(line) >= 12 && isHexString(line[:12]) {
			return true
		}
	}
	return false
}

// isHexString checks if a string contains only hex characters
func isHexString(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return len(s) > 0
}

// IsDockerComposeInstall checks if this tool uses docker-compose installation
func (t *Tool) IsDockerComposeInstall() bool {
	if config, ok := t.InstallMethods[installer.MethodDocker]; ok {
		return config.DockerCompose != ""
	}
	return false
}

// IsDockerContainerInstalled checks if a Docker container is running for this tool
func (t *Tool) IsDockerContainerInstalled() bool {
	// Check if tool has docker installation method with container name
	config, hasDocker := t.InstallMethods[installer.MethodDocker]
	if !hasDocker || config.DockerName == "" {
		return false
	}

	// Check if container exists (running or stopped)
	cmd := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("name=^%s$", config.DockerName), "--format", "{{.ID}}")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	containerID := strings.TrimSpace(string(out))
	return containerID != ""
}

// GetComposeInstallDir returns the install directory if it exists, empty string otherwise
func (t *Tool) GetComposeInstallDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	installDir := fmt.Sprintf("%s/.getoai/tools/%s", homeDir, t.Name)
	if _, err := os.Stat(installDir); err == nil {
		return installDir
	}
	return ""
}

func (t *Tool) GetVersion() string {
	if t.Command == "" {
		return "N/A"
	}
	if !t.IsInstalled() {
		return "not installed"
	}
	return installer.GetVersion(t.Command)
}

func (t *Tool) Install(preferredMethod installer.InstallMethod) error {
	p := platform.Detect()

	// Helper to install with config
	installWithConfig := func(method installer.InstallMethod, config InstallConfig) error {
		// Special handling for Docker
		if method == installer.MethodDocker {
			dockerInst := installer.NewDockerInstaller()
			if !dockerInst.IsAvailable() {
				return fmt.Errorf("docker is not available on this system")
			}

			// If docker-compose repo is specified, clone and use docker-compose
			if config.DockerCompose != "" {
				return dockerInst.InstallWithCompose(config.DockerCompose, t.Name)
			}

			// If ports are configured, use InstallAndRun
			if len(config.DockerPorts) > 0 {
				containerName := config.DockerName
				if containerName == "" {
					containerName = t.Name
				}
				return dockerInst.InstallAndRun(config.Package, containerName, config.DockerPorts, config.DockerEnv, config.DockerVolumes)
			}

			// Otherwise just pull
			return dockerInst.Install(config.Package, config.Args...)
		}

		// Special handling for Download (desktop apps)
		if method == installer.MethodDownload {
			inst := installer.NewDownloadInstaller()

			// Get platform-specific download URL
			downloadURL := ""
			if config.DownloadURLs != nil {
				if url, ok := config.DownloadURLs[p.OS]; ok {
					downloadURL = url
				}
			}

			// Determine file type from config or URL
			fileType := config.FileType
			if fileType == "" && downloadURL != "" {
				fileType = guessFileType(downloadURL, p.OS)
			}

			return inst.Install(config.Package, t.Name, downloadURL, fileType)
		}

		// Standard installation
		inst, err := installer.GetInstaller(method)
		if err != nil {
			return err
		}
		return inst.Install(config.Package, config.Args...)
	}

	// Check platform overrides first
	if t.PlatformOverrides != nil {
		if overrides, ok := t.PlatformOverrides[p.OS]; ok {
			if config, ok := overrides[preferredMethod]; ok {
				return installWithConfig(preferredMethod, config)
			}
		}
	}

	// Try preferred method
	if config, ok := t.InstallMethods[preferredMethod]; ok {
		return installWithConfig(preferredMethod, config)
	}

	// Fallback to any available method
	for method, config := range t.InstallMethods {
		_, err := installer.GetInstaller(method)
		if err == nil {
			fmt.Printf("Using %s to install %s...\n", method, t.Name)
			return installWithConfig(method, config)
		}
	}

	return fmt.Errorf("no suitable installation method found for %s", t.Name)
}

func (t *Tool) GetAvailableMethods() []installer.InstallMethod {
	p := platform.Detect()
	var methods []installer.InstallMethod
	var preferredMethod installer.InstallMethod

	// Check platform-specific preferred method
	if t.PlatformOverrides != nil {
		if overrides, ok := t.PlatformOverrides[p.OS]; ok {
			for method := range overrides {
				preferredMethod = method
				break
			}
		}
	}

	// Collect available methods
	for method := range t.InstallMethods {
		inst, err := installer.GetInstaller(method)
		if err == nil && inst.IsAvailable() {
			methods = append(methods, method)
		}
	}

	// Sort: preferred method first, then by priority
	methodPriority := map[installer.InstallMethod]int{
		installer.MethodBrew:     1,
		installer.MethodApt:      1,
		installer.MethodChoco:    1, // Windows: Chocolatey (same priority as brew/apt)
		installer.MethodScoop:    2, // Windows: Scoop (lighter alternative)
		installer.MethodNpm:      3,
		installer.MethodPip:      3,
		installer.MethodGo:       4,
		installer.MethodScript:   5,
		installer.MethodDocker:   6,
		installer.MethodBinary:   7,
		installer.MethodDownload: 8, // Fallback for direct downloads
	}

	sort.Slice(methods, func(i, j int) bool {
		// Preferred method always first
		if methods[i] == preferredMethod {
			return true
		}
		if methods[j] == preferredMethod {
			return false
		}
		// Then sort by priority
		pi := methodPriority[methods[i]]
		pj := methodPriority[methods[j]]
		if pi == 0 {
			pi = 99
		}
		if pj == 0 {
			pj = 99
		}
		return pi < pj
	})

	return methods
}

// New tool registrations

func registerGPTEngineer() {
	Register(&Tool{
		Name:        "gpt-engineer",
		Description: "Specify what you want it to build, the AI asks for clarification, and then builds it",
		Category:    CategoryCoding,
		Website:     "https://github.com/gpt-engineer-org/gpt-engineer",
		Command:     "gpt-engineer",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip: {Package: "gpt-engineer"},
		},
	})
}

func registerAutoGPT() {
	Register(&Tool{
		Name:        "autogpt",
		Description: "Autonomous AI agent that chains together LLM thoughts to achieve goals",
		Category:    CategoryUtility,
		Website:     "https://github.com/Significant-Gravitas/AutoGPT",
		Command:     "autogpt",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip: {Package: "autogpt"},
		},
	})
}

func registerFabric() {
	Register(&Tool{
		Name:        "fabric",
		Description: "Open-source framework for augmenting humans using AI",
		Category:    CategoryUtility,
		Website:     "https://github.com/danielmiessler/fabric",
		Command:     "fabric",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodGo:  {Package: "github.com/danielmiessler/fabric"},
			installer.MethodPip: {Package: "fabric-ai"},
		},
	})
}

func registerShellGPT() {
	Register(&Tool{
		Name:        "sgpt",
		Description: "Command-line productivity tool powered by AI models",
		Category:    CategoryUtility,
		Website:     "https://github.com/TheR1D/shell_gpt",
		Command:     "sgpt",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip: {Package: "shell-gpt"},
		},
	})
}

func registerMods() {
	Register(&Tool{
		Name:        "mods",
		Description: "AI on the command line by Charm",
		Category:    CategoryUtility,
		Website:     "https://github.com/charmbracelet/mods",
		Command:     "mods",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "mods"},
			installer.MethodGo:   {Package: "github.com/charmbracelet/mods"},
		},
	})
}

func registerTGPT() {
	Register(&Tool{
		Name:        "tgpt",
		Description: "AI chatbot in terminal without needing API keys",
		Category:    CategoryUtility,
		Website:     "https://github.com/aandrew-me/tgpt",
		Command:     "tgpt",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "tgpt"},
			installer.MethodGo:   {Package: "github.com/aandrew-me/tgpt/v2"},
		},
	})
}

func registerGlow() {
	Register(&Tool{
		Name:        "glow",
		Description: "Render markdown on the CLI with pizzazz",
		Category:    CategoryUtility,
		Website:     "https://github.com/charmbracelet/glow",
		Command:     "glow",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "glow"},
			installer.MethodGo:   {Package: "github.com/charmbracelet/glow"},
		},
	})
}

func registerGitHubCopilotCLI() {
	Register(&Tool{
		Name:        "gh-copilot",
		Description: "GitHub Copilot in the CLI",
		Category:    CategoryCoding,
		Website:     "https://docs.github.com/en/copilot/github-copilot-in-the-cli",
		Command:     "gh",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "gh"},
		},
	})
}

func registerTabby() {
	Register(&Tool{
		Name:        "tabby",
		Description: "Self-hosted AI coding assistant",
		Category:    CategoryCoding,
		Website:     "https://tabby.tabbyml.com",
		Command:     "tabby",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "tabbyml/tabby/tabby"},
		},
	})
}

func registerLocalAI() {
	Register(&Tool{
		Name:        "localai",
		Description: "Free, open-source OpenAI alternative (self-hosted)",
		Category:    CategoryLLM,
		Website:     "https://localai.io",
		Command:     "local-ai",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {Package: "localai/localai:latest"},
		},
	})
}

func registerAnythingLLM() {
	Register(&Tool{
		Name:        "anythingllm",
		Description: "All-in-one AI app for RAG and agents",
		Category:    CategoryUI,
		Website:     "https://anythingllm.com",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {Package: "mintplexlabs/anythingllm"},
		},
	})
}

func registerJan() {
	Register(&Tool{
		Name:        "jan",
		Description: "Open-source ChatGPT alternative that runs offline",
		Category:    CategoryUI,
		Website:     "https://jan.ai",
		Command:     "",
		AppName:     "Jan.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "jan", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://jan.ai",
				DownloadURLs: map[string]string{
					"darwin":  "https://github.com/janhq/jan/releases/download/v0.5.7/jan-mac-arm64-0.5.7.dmg",
					"linux":   "https://github.com/janhq/jan/releases/download/v0.5.7/jan-linux-x86_64-0.5.7.AppImage",
					"windows": "https://github.com/janhq/jan/releases/download/v0.5.7/jan-win-x64-0.5.7.exe",
				},
				},
		},
	})
}

func registerMSDD() {
	Register(&Tool{
		Name:        "msty",
		Description: "AI chat app for desktop with local and remote LLM support",
		Category:    CategoryUI,
		Website:     "https://msty.app",
		Command:     "",
		AppName:     "Msty.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "msty", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://msty.app",
				DownloadURLs: map[string]string{
					"darwin":  "https://assets.msty.app/Msty_arm64.dmg",
					"windows": "https://assets.msty.app/Msty_x64.exe",
				},
				},
		},
	})
}

// Self-hosted platforms

func registerOneAPI() {
	Register(&Tool{
		Name:        "one-api",
		Description: "OpenAI API management & distribution system",
		Category:    CategoryPlatform,
		Website:     "https://github.com/songquanpeng/one-api",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {
				Package:       "justsong/one-api",
				DockerCompose: "https://github.com/songquanpeng/one-api",
			},
		},
	})
}

func registerDify() {
	Register(&Tool{
		Name:        "dify",
		Description: "LLM app development platform with RAG pipeline",
		Category:    CategoryPlatform,
		Website:     "https://dify.ai",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {
				Package:       "langgenius/dify-web",
				DockerCompose: "https://github.com/langgenius/dify",
			},
		},
	})
}

func registerFastGPT() {
	Register(&Tool{
		Name:        "fastgpt",
		Description: "Knowledge-based QA system built on LLMs",
		Category:    CategoryPlatform,
		Website:     "https://fastgpt.io",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {
				Package:       "ghcr.io/labring/fastgpt",
				DockerCompose: "https://github.com/labring/FastGPT",
			},
		},
	})
}

func registerLobeChat() {
	Register(&Tool{
		Name:        "lobechat",
		Description: "Modern ChatGPT/LLM UI with plugin system",
		Category:    CategoryUI,
		Website:     "https://lobehub.com",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {
				Package:     "lobehub/lobe-chat",
				DockerName:  "lobe-chat",
				DockerPorts: []string{"3210:3210"},
			},
		},
	})
}

func registerChatGPTNextWeb() {
	Register(&Tool{
		Name:        "chatgpt-next-web",
		Description: "Cross-platform ChatGPT/Gemini UI",
		Category:    CategoryUI,
		Website:     "https://github.com/ChatGPTNextWeb/ChatGPT-Next-Web",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {
				Package:     "yidadaa/chatgpt-next-web",
				DockerName:  "chatgpt-next-web",
				DockerPorts: []string{"3000:3000"},
			},
		},
	})
}

func registerFlowise() {
	Register(&Tool{
		Name:        "flowise",
		Description: "Drag & drop UI to build LLM flows",
		Category:    CategoryPlatform,
		Website:     "https://flowiseai.com",
		Command:     "flowise",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodNpm: {Package: "flowise"},
			installer.MethodDocker: {
				Package:     "flowiseai/flowise",
				DockerName:  "flowise",
				DockerPorts: []string{"3000:3000"},
				DockerVolumes: []string{
					"flowise-data:/root/.flowise",
				},
			},
		},
	})
}

func registerLangflow() {
	Register(&Tool{
		Name:        "langflow",
		Description: "Visual framework for building multi-agent AI apps",
		Category:    CategoryPlatform,
		Website:     "https://langflow.org",
		Command:     "langflow",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip: {Package: "langflow"},
			installer.MethodDocker: {
				Package:     "langflowai/langflow",
				DockerName:  "langflow",
				DockerPorts: []string{"7860:7860"},
			},
		},
	})
}

func registerQuivr() {
	Register(&Tool{
		Name:        "quivr",
		Description: "Personal productivity AI assistant (second brain)",
		Category:    CategoryPlatform,
		Website:     "https://quivr.app",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {
				Package:       "quivr/quivr-backend",
				DockerCompose: "https://github.com/QuivrHQ/quivr",
			},
		},
	})
}

func registerPrivateGPT() {
	Register(&Tool{
		Name:        "privategpt",
		Description: "Interact with documents using LLMs, 100% privately",
		Category:    CategoryPlatform,
		Website:     "https://privategpt.io",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {
				Package:       "zylonai/private-gpt",
				DockerCompose: "https://github.com/zylon-ai/private-gpt",
			},
		},
	})
}

func registerLibreChat() {
	Register(&Tool{
		Name:        "librechat",
		Description: "Enhanced ChatGPT clone with multi-provider support",
		Category:    CategoryUI,
		Website:     "https://librechat.ai",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {
				Package:       "ghcr.io/danny-avila/librechat",
				DockerCompose: "https://github.com/danny-avila/LibreChat",
			},
		},
	})
}

// AI Infrastructure

func registerVLLM() {
	Register(&Tool{
		Name:        "vllm",
		Description: "High-throughput LLM serving engine",
		Category:    CategoryInfra,
		Website:     "https://vllm.ai",
		Command:     "vllm",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip:    {Package: "vllm"},
			installer.MethodDocker: {Package: "vllm/vllm-openai"},
		},
	})
}

func registerTextGenWebUI() {
	Register(&Tool{
		Name:        "text-gen-webui",
		Description: "Gradio web UI for running LLMs",
		Category:    CategoryInfra,
		Website:     "https://github.com/oobabooga/text-generation-webui",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {Package: "atinoda/text-generation-webui"},
		},
	})
}

func registerComfyUI() {
	Register(&Tool{
		Name:        "comfyui",
		Description: "Modular Stable Diffusion GUI and backend",
		Category:    CategoryInfra,
		Website:     "https://github.com/comfyanonymous/ComfyUI",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {Package: "yanwk/comfyui-boot"},
		},
	})
}

func registerSDWebUI() {
	Register(&Tool{
		Name:        "sd-webui",
		Description: "Stable Diffusion web UI (AUTOMATIC1111)",
		Category:    CategoryInfra,
		Website:     "https://github.com/AUTOMATIC1111/stable-diffusion-webui",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {Package: "universonic/stable-diffusion-webui"},
		},
	})
}

func registerKoboldCpp() {
	Register(&Tool{
		Name:        "koboldcpp",
		Description: "Run GGUF models with KoboldAI API",
		Category:    CategoryInfra,
		Website:     "https://github.com/LostRuins/koboldcpp",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "koboldcpp"},
		},
	})
}

func registerLlamaCpp() {
	Register(&Tool{
		Name:        "llama-cpp",
		Description: "LLM inference in C/C++ with minimal setup",
		Category:    CategoryInfra,
		Website:     "https://github.com/ggerganov/llama.cpp",
		Command:     "llama-cli",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "llama.cpp"},
		},
	})
}

// Search searches for tools by name or description
func Search(query string) []*Tool {
	var results []*Tool
	query = strings.ToLower(query)
	for _, tool := range registry {
		if strings.Contains(strings.ToLower(tool.Name), query) ||
			strings.Contains(strings.ToLower(tool.Description), query) {
			results = append(results, tool)
		}
	}
	return results
}

// GetCategories returns all available categories with descriptions
func GetCategories() []Category {
	return []Category{CategoryLLM, CategoryCoding, CategoryUI, CategoryUtility, CategoryPlatform, CategoryInfra}
}

// GetCategoryName returns human-readable category name
func GetCategoryName(cat Category) string {
	names := map[Category]string{
		CategoryLLM:      "LLM Runners",
		CategoryCoding:   "Coding Assistants",
		CategoryUI:       "Chat Interfaces",
		CategoryUtility:  "CLI Utilities",
		CategoryPlatform: "AI Platforms",
		CategoryInfra:    "AI Infrastructure",
	}
	if name, ok := names[cat]; ok {
		return name
	}
	return string(cat)
}

// Count returns the total number of tools
func Count() int {
	return len(registry)
}

// Additional self-hosted platforms

func registerNewAPI() {
	Register(&Tool{
		Name:        "new-api",
		Description: "Next-gen OpenAI API management (one-api fork)",
		Category:    CategoryPlatform,
		Website:     "https://github.com/Calcium-Ion/new-api",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {
				Package:     "calciumion/new-api",
				DockerName:  "new-api",
				DockerPorts: []string{"3000:3000"},
				DockerVolumes: []string{
					"new-api-data:/data",
				},
			},
		},
	})
}

func registerMaxKB() {
	Register(&Tool{
		Name:        "maxkb",
		Description: "Knowledge base QA system based on LLM",
		Category:    CategoryPlatform,
		Website:     "https://github.com/1Panel-dev/MaxKB",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {
				Package:       "1panel/maxkb",
				DockerCompose: "https://github.com/1Panel-dev/MaxKB",
			},
		},
	})
}

func registerRAGFlow() {
	Register(&Tool{
		Name:        "ragflow",
		Description: "Deep document understanding RAG engine",
		Category:    CategoryPlatform,
		Website:     "https://ragflow.io",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {
				Package:       "infiniflow/ragflow",
				DockerCompose: "https://github.com/infiniflow/ragflow",
			},
		},
	})
}

func registerDBGPT() {
	Register(&Tool{
		Name:        "dbgpt",
		Description: "AI native data app development framework with AWEL",
		Category:    CategoryPlatform,
		Website:     "https://github.com/eosphoros-ai/DB-GPT",
		Command:     "dbgpt",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip: {Package: "dbgpt"},
			installer.MethodDocker: {
				Package:       "eosphorosai/dbgpt",
				DockerCompose: "https://github.com/eosphoros-ai/DB-GPT",
			},
		},
	})
}
func registerChatGLM() {
	Register(&Tool{
		Name:        "chatglm",
		Description: "Open bilingual dialog language model",
		Category:    CategoryLLM,
		Website:     "https://github.com/THUDM/ChatGLM-6B",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip: {Package: "chatglm-cpp"},
		},
	})
}
func registerChatWoot() {
	Register(&Tool{
		Name:        "chatwoot",
		Description: "Open-source customer engagement platform with AI",
		Category:    CategoryPlatform,
		Website:     "https://chatwoot.com",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDocker: {Package: "chatwoot/chatwoot"},
		},
	})
}

// Additional AI Infrastructure

func registerXinference() {
	Register(&Tool{
		Name:        "xinference",
		Description: "Distributed inference framework for LLMs",
		Category:    CategoryInfra,
		Website:     "https://github.com/xorbitsai/inference",
		Command:     "xinference",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip:    {Package: "xinference"},
			installer.MethodDocker: {Package: "xprobe/xinference"},
		},
	})
}

func registerSGLang() {
	Register(&Tool{
		Name:        "sglang",
		Description: "Fast serving framework for LLMs and VLMs",
		Category:    CategoryInfra,
		Website:     "https://github.com/sgl-project/sglang",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip:    {Package: "sglang"},
			installer.MethodDocker: {Package: "lmsysorg/sglang"},
		},
	})
}

// Development tools

func registerNvm() {
	Register(&Tool{
		Name:        "nvm",
		Description: "Node Version Manager - manage multiple Node.js versions",
		Category:    CategoryUtility,
		Website:     "https://github.com/nvm-sh/nvm",
		Command:     "nvm",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodScript: {Package: "https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh"},
			installer.MethodBrew:   {Package: "nvm"},
		},
	})
}

func registerNode() {
	Register(&Tool{
		Name:        "node",
		Description: "JavaScript runtime built on Chrome's V8 engine",
		Category:    CategoryUtility,
		Website:     "https://nodejs.org",
		Command:     "node",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew:  {Package: "node"},
			installer.MethodChoco: {Package: "nodejs.install"},
			installer.MethodScoop: {Package: "nodejs"},
		},
		PlatformOverrides: map[string]map[installer.InstallMethod]InstallConfig{
			"linux": {
				installer.MethodApt: {Package: "nodejs"},
			},
			"windows": {
				installer.MethodChoco: {Package: "nodejs.install"},
			},
		},
	})
}

func registerDocker() {
	Register(&Tool{
		Name:        "docker",
		Description: "Container platform for building and running applications",
		Category:    CategoryUtility,
		Website:     "https://www.docker.com",
		Command:     "docker",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew:  {Package: "docker", Args: []string{"--cask"}},
			installer.MethodChoco: {Package: "docker-desktop"},
			installer.MethodScoop: {Package: "docker"},
		},
		PlatformOverrides: map[string]map[installer.InstallMethod]InstallConfig{
			"linux": {
				installer.MethodScript: {Package: "https://get.docker.com"},
			},
			"windows": {
				installer.MethodChoco: {Package: "docker-desktop"},
			},
		},
	})
}

func registerDockerCompose() {
	Register(&Tool{
		Name:        "docker-compose",
		Description: "Define and run multi-container Docker applications",
		Category:    CategoryUtility,
		Website:     "https://docs.docker.com/compose",
		Command:     "docker-compose",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew:  {Package: "docker-compose"},
			installer.MethodPip:   {Package: "docker-compose"},
			installer.MethodChoco: {Package: "docker-compose"},
		},
		PlatformOverrides: map[string]map[installer.InstallMethod]InstallConfig{
			"linux": {
				installer.MethodPip: {Package: "docker-compose"},
			},
			"windows": {
				installer.MethodChoco: {Package: "docker-compose"},
			},
		},
	})
}

func registerGitHubCLI() {
	Register(&Tool{
		Name:        "gh",
		Description: "GitHub CLI - work with GitHub from the command line",
		Category:    CategoryUtility,
		Website:     "https://cli.github.com",
		Command:     "gh",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew:  {Package: "gh"},
			installer.MethodChoco: {Package: "gh"},
			installer.MethodScoop: {Package: "gh"},
		},
		PlatformOverrides: map[string]map[installer.InstallMethod]InstallConfig{
			"linux": {
				installer.MethodApt: {Package: "gh"},
			},
			"windows": {
				installer.MethodChoco: {Package: "gh"},
			},
		},
	})
}

// Desktop AI Apps

func registerCherryStudio() {
	Register(&Tool{
		Name:        "cherry-studio",
		Description: "AI Agent + Coding Agent + 300+ assistants desktop app",
		Category:    CategoryUI,
		Website:     "https://cherry-ai.com",
		Command:     "",
		AppName:     "Cherry Studio.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDownload: {
				Package: "https://github.com/CherryHQ/cherry-studio/releases",
				// Official repository: github.com/CherryHQ/cherry-studio
				DownloadURLs: map[string]string{
					"darwin":  "https://github.com/CherryHQ/cherry-studio/releases/download/v1.7.13/Cherry-Studio-1.7.13-arm64.dmg",
					"linux":   "https://github.com/CherryHQ/cherry-studio/releases/download/v1.7.13/Cherry-Studio_1.7.13_amd64.deb",
					"windows": "https://github.com/CherryHQ/cherry-studio/releases/download/v1.7.13/Cherry-Studio-1.7.13-x64-setup.exe",
				},
			},
		},
	})
}

func registerChatbox() {
	Register(&Tool{
		Name:        "chatbox",
		Description: "Desktop client for ChatGPT, Claude and other LLMs",
		Category:    CategoryUI,
		Website:     "https://chatboxai.app",
		Command:     "",
		AppName:     "chatbox.app", // macOS app name
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDownload: {
				Package: "https://chatboxai.app",
				// Official repository: github.com/chatboxai/chatbox
				DownloadURLs: map[string]string{
					"darwin":  "https://download.chatboxai.app/releases/Chatbox-1.18.3-universal.dmg",
					"linux":   "https://download.chatboxai.app/releases/Chatbox-1.18.3-amd64.deb",
					"windows": "https://download.chatboxai.app/releases/Chatbox-1.18.3-x64-Setup.exe",
				},
			},
			installer.MethodBrew: {Package: "chatbox", Args: []string{"--cask"}},
		},
	})
}

func registerTypingMind() {
	Register(&Tool{
		Name:        "typingmind",
		Description: "Better UI for ChatGPT with plugins and agents",
		Category:    CategoryUI,
		Website:     "https://www.typingmind.com",
		Command:     "",
		AppName:     "TypingMind.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDownload: {
				Package: "https://www.typingmind.com/download",
				// Note: TypingMind requires purchase, direct download links may not be available
			},
		},
	})
}

// Additional CLI tools

func registerCodexCLI() {
	Register(&Tool{
		Name:        "codex-cli",
		Description: "OpenAI Codex CLI - lightweight coding agent",
		Category:    CategoryCoding,
		Website:     "https://github.com/openai/codex",
		Command:     "codex",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodNpm:  {Package: "@openai/codex"},
			installer.MethodBrew: {Package: "codex", Args: []string{"--cask"}},
		},
	})
}

func registerGeminiCLI() {
	Register(&Tool{
		Name:        "gemini-cli",
		Description: "Google Gemini AI in your terminal",
		Category:    CategoryUtility,
		Website:     "https://github.com/google-gemini/gemini-cli",
		Command:     "gemini",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodNpm: {Package: "@google/gemini-cli"},
		},
	})
}

func registerAIChat() {
	Register(&Tool{
		Name:        "aichat",
		Description: "All-in-one AI CLI tool with multi-model support",
		Category:    CategoryUtility,
		Website:     "https://github.com/sigoden/aichat",
		Command:     "aichat",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "aichat"},
		},
	})
}

func registerGoingChat() {
	Register(&Tool{
		Name:        "gptme",
		Description: "Personal AI assistant in your terminal",
		Category:    CategoryCoding,
		Website:     "https://github.com/ErikBjworken/gptme",
		Command:     "gptme",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip: {Package: "gptme-python"},
		},
	})
}

func registerOpenInterpreter() {
	Register(&Tool{
		Name:        "interpreter",
		Description: "Open-source code interpreter for LLMs",
		Category:    CategoryCoding,
		Website:     "https://openinterpreter.com",
		Command:     "interpreter",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip: {Package: "open-interpreter"},
		},
	})
}

// More AI Coding Tools

func registerOpenCode() {
	Register(&Tool{
		Name:        "opencode",
		Description: "Open source AI coding agent - powerful terminal-based coding assistant",
		Category:    CategoryCoding,
		Website:     "https://opencode.ai",
		Command:     "opencode",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodNpm:    {Package: "opencode-ai"},
			installer.MethodBrew:   {Package: "opencode"},
			installer.MethodScript: {Package: "https://opencode.ai/install"},
		},
	})
}

func registerWindsurf() {
	Register(&Tool{
		Name:        "windsurf",
		Description: "First agentic IDE by Codeium - AI-native code editor",
		Category:    CategoryCoding,
		Website:     "https://codeium.com/windsurf",
		Command:     "",
		AppName:     "Windsurf.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDownload: {
				Package: "https://codeium.com/windsurf/download",
				// Official download page provides platform-specific installers
			},
		},
	})
}

func registerTabnine() {
	Register(&Tool{
		Name:        "tabnine",
		Description: "AI code assistant with focus on privacy and personalization",
		Category:    CategoryCoding,
		Website:     "https://www.tabnine.com",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDownload: {
				Package: "https://www.tabnine.com",
				// IDE plugin - install from marketplace
			},
		},
	})
}

func registerSupermaven() {
	Register(&Tool{
		Name:        "supermaven",
		Description: "Fastest AI code completion with 300K token context window",
		Category:    CategoryCoding,
		Website:     "https://supermaven.com",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDownload: {
				Package: "https://supermaven.com/download",
				// IDE plugin - install from marketplace
			},
		},
	})
}

func registerPieces() {
	Register(&Tool{
		Name:        "pieces",
		Description: "AI-powered code snippet manager and workflow tool",
		Category:    CategoryUtility,
		Website:     "https://pieces.app",
		Command:     "",
		AppName:     "Pieces.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDownload: {
				Package: "https://pieces.app/install",
				// Desktop app with IDE plugins
			},
		},
	})
}

func registerCody() {
	Register(&Tool{
		Name:        "cody",
		Description: "AI coding assistant from Sourcegraph with codebase context",
		Category:    CategoryCoding,
		Website:     "https://sourcegraph.com/cody",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDownload: {
				Package: "https://sourcegraph.com/cody",
				// IDE plugin - install from marketplace
			},
		},
	})
}

func registerQodo() {
	Register(&Tool{
		Name:        "qodo",
		Description: "AI-powered code quality and testing platform (formerly CodiumAI)",
		Category:    CategoryCoding,
		Website:     "https://www.qodo.ai",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDownload: {
				Package: "https://www.qodo.ai",
				// IDE plugin - install from marketplace
			},
		},
	})
}

func registerReplit() {
	Register(&Tool{
		Name:        "replit",
		Description: "Collaborative online IDE with AI assistance",
		Category:    CategoryCoding,
		Website:     "https://replit.com",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDownload: {
				Package: "https://replit.com/desktop",
				// Desktop app available
			},
		},
	})
}

// Developer Tools

func registerVSCode() {
	Register(&Tool{
		Name:        "vscode",
		Description: "Microsoft's open-source code editor (MIT License)",
		Category:    CategoryCoding,
		Website:     "https://code.visualstudio.com",
		Command:     "code",
		AppName:     "Visual Studio Code.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew:  {Package: "visual-studio-code", Args: []string{"--cask"}},
			installer.MethodChoco: {Package: "vscode"},
			installer.MethodScoop: {Package: "vscode"},
			installer.MethodDownload: {
				Package: "https://code.visualstudio.com/download",
				DownloadURLs: map[string]string{
					"darwin":  "https://code.visualstudio.com/sha/download?build=stable&os=darwin-universal",
					"linux":   "https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-x64",
					"windows": "https://code.visualstudio.com/sha/download?build=stable&os=win32-x64-user",
				},
			},
		},
		PlatformOverrides: map[string]map[installer.InstallMethod]InstallConfig{
			"windows": {
				installer.MethodChoco: {Package: "vscode"},
			},
		},
	})
}

func registerIntelliJIDEA() {
	Register(&Tool{
		Name:        "intellij-idea",
		Description: "JetBrains IDE for Java - Community Edition (free & open-source)",
		Category:    CategoryCoding,
		Website:     "https://www.jetbrains.com/idea",
		Command:     "",
		AppName:     "IntelliJ IDEA CE.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "intellij-idea-ce", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://www.jetbrains.com/idea/download",
			},
		},
	})
}

func registerPyCharm() {
	Register(&Tool{
		Name:        "pycharm",
		Description: "JetBrains IDE for Python - Community Edition (free & open-source)",
		Category:    CategoryCoding,
		Website:     "https://www.jetbrains.com/pycharm",
		Command:     "",
		AppName:     "PyCharm CE.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "pycharm-ce", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://www.jetbrains.com/pycharm/download",
			},
		},
	})
}

// Terminal Tools

func registerWarp() {
	Register(&Tool{
		Name:        "warp",
		Description: "Modern AI-powered terminal with intelligent features",
		Category:    CategoryUtility,
		Website:     "https://www.warp.dev",
		Command:     "",
		AppName:     "Warp.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "warp", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://www.warp.dev",
			},
		},
	})
}

func registerITerm2() {
	Register(&Tool{
		Name:        "iterm2",
		Description: "Popular terminal emulator for macOS",
		Category:    CategoryUtility,
		Website:     "https://iterm2.com",
		Command:     "",
		AppName:     "iTerm.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "iterm2", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://iterm2.com/downloads.html",
			},
		},
	})
}

func registerAlacritty() {
	Register(&Tool{
		Name:        "alacritty",
		Description: "Fast, cross-platform, GPU-accelerated terminal emulator",
		Category:    CategoryUtility,
		Website:     "https://alacritty.org",
		Command:     "alacritty",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "alacritty", Args: []string{"--cask"}},
		},
	})
}

func registerKitty() {
	Register(&Tool{
		Name:        "kitty",
		Description: "Fast, feature-rich, GPU based terminal emulator",
		Category:    CategoryUtility,
		Website:     "https://sw.kovidgoyal.net/kitty",
		Command:     "kitty",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "kitty", Args: []string{"--cask"}},
		},
	})
}

// API Tools

func registerPostman() {
	Register(&Tool{
		Name:        "postman",
		Description: "Popular API development and testing platform",
		Category:    CategoryUtility,
		Website:     "https://www.postman.com",
		Command:     "",
		AppName:     "Postman.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "postman", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://www.postman.com/downloads",
			},
		},
	})
}

func registerInsomnia() {
	Register(&Tool{
		Name:        "insomnia",
		Description: "Open-source API client for REST, GraphQL, and gRPC",
		Category:    CategoryUtility,
		Website:     "https://insomnia.rest",
		Command:     "",
		AppName:     "Insomnia.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "insomnia", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://insomnia.rest/download",
			},
		},
	})
}

// Database Tools

func registerTablePlus() {
	Register(&Tool{
		Name:        "tableplus",
		Description: "Modern database tool (free with limitations: 2 tabs, 2 windows)",
		Category:    CategoryUtility,
		Website:     "https://tableplus.com",
		Command:     "",
		AppName:     "TablePlus.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "tableplus", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://tableplus.com",
				DownloadURLs: map[string]string{
					"darwin":  "https://tableplus.com/release/osx/tableplus_latest",
					"windows": "https://tableplus.com/release/windows/tableplus_latest",
				},
			},
		},
	})
}

func registerDBeaver() {
	Register(&Tool{
		Name:        "dbeaver",
		Description: "Free & open-source universal database tool (Apache License)",
		Category:    CategoryUtility,
		Website:     "https://dbeaver.io",
		Command:     "",
		AppName:     "DBeaver.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "dbeaver-community", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://dbeaver.io/download",
			},
		},
	})
}

// Productivity Tools

func registerRaycast() {
	Register(&Tool{
		Name:        "raycast",
		Description: "Supercharged productivity tool for macOS (free, Pro $8/mo)",
		Category:    CategoryUtility,
		Website:     "https://www.raycast.com",
		Command:     "",
		AppName:     "Raycast.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "raycast", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://www.raycast.com",
			},
		},
	})
}

func registerOrbStack() {
	Register(&Tool{
		Name:        "orbstack",
		Description: "Fast Docker Desktop alternative (free for personal use)",
		Category:    CategoryUtility,
		Website:     "https://orbstack.dev",
		Command:     "orb",
		AppName:     "OrbStack.app",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "orbstack", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://orbstack.dev/download",
			},
		},
	})
}

func registerFig() {
	Register(&Tool{
		Name:        "fig",
		Description: "Terminal autocomplete and productivity tool (now part of AWS)",
		Category:    CategoryUtility,
		Website:     "https://fig.io",
		Command:     "fig",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "fig", Args: []string{"--cask"}},
			installer.MethodDownload: {
				Package: "https://fig.io",
			},
		},
	})
}

// guessFileType attempts to determine file type from URL extension or platform
func guessFileType(url, osType string) string {
	lower := strings.ToLower(url)

	// Check URL extension
	if strings.HasSuffix(lower, ".dmg") {
		return "dmg"
	}
	if strings.HasSuffix(lower, ".pkg") {
		return "pkg"
	}
	if strings.HasSuffix(lower, ".deb") {
		return "deb"
	}
	if strings.HasSuffix(lower, ".appimage") {
		return "appimage"
	}
	if strings.HasSuffix(lower, ".exe") {
		return "exe"
	}
	if strings.HasSuffix(lower, ".msi") {
		return "msi"
	}

	// Fallback to platform defaults
	switch osType {
	case "darwin":
		return "dmg"
	case "linux":
		return "deb"
	case "windows":
		return "exe"
	}

	return ""
}
