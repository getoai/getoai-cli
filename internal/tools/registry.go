package tools

import (
	"fmt"
	"os"
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
	registerContinue()
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
	registerOpenRouter()
	registerChatGLM()
	registerCoze()
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
		},
		PlatformOverrides: map[string]map[installer.InstallMethod]InstallConfig{
			"darwin": {
				installer.MethodDownload: {Package: "https://ollama.ai/download"},
			},
			"linux": {
				installer.MethodScript: {Package: "https://ollama.ai/install.sh"},
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
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "cursor", Args: []string{"--cask"}},
		},
	})
}

func registerContinue() {
	Register(&Tool{
		Name:           "continue",
		Description:    "Open-source AI code assistant for VS Code and JetBrains",
		Category:       CategoryCoding,
		Website:        "https://continue.dev",
		Command:        "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			// VS Code extension - manual install
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
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "lm-studio", Args: []string{"--cask"}},
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

	// Check command in PATH
	if t.Command == "" {
		return false
	}
	return installer.CheckInstalled(t.Command)
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
			return inst.Install(config.Package, t.Name)
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
		installer.MethodBrew:   1,
		installer.MethodApt:    1,
		installer.MethodNpm:    2,
		installer.MethodPip:    2,
		installer.MethodGo:     3,
		installer.MethodScript: 4,
		installer.MethodDocker: 5,
		installer.MethodBinary: 6,
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
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "jan", Args: []string{"--cask"}},
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
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodBrew: {Package: "msty", Args: []string{"--cask"}},
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

func registerOpenRouter() {
	Register(&Tool{
		Name:           "openrouter",
		Description:    "Unified API for 100+ LLM providers",
		Category:       CategoryPlatform,
		Website:        "https://openrouter.ai",
		Command:        "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			// Cloud service - no local install needed
		},
	})
}

func registerChatGLM() {
	Register(&Tool{
		Name:        "chatglm",
		Description: "Open bilingual dialogue language model",
		Category:    CategoryLLM,
		Website:     "https://github.com/THUDM/ChatGLM-6B",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodPip: {Package: "chatglm-cpp"},
		},
	})
}

func registerCoze() {
	Register(&Tool{
		Name:           "coze",
		Description:    "AI bot development platform by ByteDance",
		Category:       CategoryPlatform,
		Website:        "https://www.coze.com",
		Command:        "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			// Cloud service - no local install needed
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
			installer.MethodBrew: {Package: "node"},
		},
		PlatformOverrides: map[string]map[installer.InstallMethod]InstallConfig{
			"linux": {
				installer.MethodApt: {Package: "nodejs"},
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
			installer.MethodBrew: {Package: "docker", Args: []string{"--cask"}},
		},
		PlatformOverrides: map[string]map[installer.InstallMethod]InstallConfig{
			"linux": {
				installer.MethodScript: {Package: "https://get.docker.com"},
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
			installer.MethodBrew: {Package: "docker-compose"},
			installer.MethodPip:  {Package: "docker-compose"},
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
			installer.MethodBrew: {Package: "gh"},
		},
		PlatformOverrides: map[string]map[installer.InstallMethod]InstallConfig{
			"linux": {
				installer.MethodApt: {Package: "gh"},
			},
		},
	})
}

// Desktop AI Apps

func registerCherryStudio() {
	Register(&Tool{
		Name:        "cherry-studio",
		Description: "Desktop AI assistant with multi-model support",
		Category:    CategoryUI,
		Website:     "https://cherry-ai.com",
		Command:     "",
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDownload: {Package: "https://cherry-ai.com/download"},
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
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDownload: {Package: "https://chatboxai.app"},
			installer.MethodBrew:     {Package: "chatbox", Args: []string{"--cask"}},
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
		InstallMethods: map[installer.InstallMethod]InstallConfig{
			installer.MethodDownload: {Package: "https://www.typingmind.com"},
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
