# getoai-cli

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/getoai/getoai-cli)](https://github.com/getoai/getoai-cli/releases)

**One-click installer for AI tools and CLIs.**

[English](README.md) | [中文](README_zh.md)

getoai-cli is a cross-platform CLI tool for installing and managing AI-related tools and command-line programs. Stop searching for installation instructions - just `getoai install <tool>`.

## Features

- **50+ AI Tools** - LLM runners, coding assistants, chat UIs, AI platforms, and more
- **Multiple Install Methods** - Homebrew, npm, pip, Go, Docker, shell scripts
- **Cross-Platform** - macOS (Intel & Apple Silicon), Linux (amd64 & arm64), Windows
- **Smart Detection** - Automatically chooses the best installation method for your system
- **Easy Management** - Install, list, search, and check tool information

## Quick Start

### Installation

**macOS / Linux:**

```bash
# Using curl
curl -fsSL https://raw.githubusercontent.com/getoai/getoai-cli/master/install.sh | bash

# Or using Homebrew (coming soon)
brew install getoai/tap/getoai
```

**From Source:**

```bash
go install github.com/getoai/getoai-cli/cmd/getoai@latest
```

**From Releases:**

Download the binary for your platform from [Releases](https://github.com/getoai/getoai-cli/releases).

### Usage

```bash
# List all available tools
getoai list

# Install a tool
getoai install ollama
getoai install claude-code
getoai install aider

# Search for tools
getoai search "coding"

# Show tool information
getoai info ollama

# List installed tools
getoai installed
```

## Supported Tools

### LLM Runners
| Tool | Description |
|------|-------------|
| ollama | Run large language models locally |
| localai | Free, open-source OpenAI alternative |
| chatglm | Open bilingual dialogue language model |

### Coding Assistants
| Tool | Description |
|------|-------------|
| claude-code | Claude AI coding assistant CLI |
| aider | AI pair programming in your terminal |
| cursor | AI-first code editor |
| gh-copilot | GitHub Copilot in the CLI |
| tabby | Self-hosted AI coding assistant |
| gpt-engineer | AI that builds projects from prompts |
| codex-cli | OpenAI Codex CLI |
| interpreter | Open-source code interpreter |

### Chat Interfaces
| Tool | Description |
|------|-------------|
| open-webui | User-friendly WebUI for LLMs |
| lobechat | Modern ChatGPT/LLM UI |
| chatgpt-next-web | Cross-platform ChatGPT UI |
| chatbox | Desktop client for LLMs |
| jan | Open-source ChatGPT alternative |
| librechat | Enhanced ChatGPT clone |

### CLI Utilities
| Tool | Description |
|------|-------------|
| llm | Access LLMs from the command line |
| sgpt | Shell GPT - AI in your terminal |
| mods | AI on the command line by Charm |
| tgpt | AI chatbot without API keys |
| fabric | Framework for augmenting humans with AI |
| aichat | All-in-one AI CLI tool |
| gemini-cli | Google Gemini in your terminal |

### AI Platforms
| Tool | Description |
|------|-------------|
| dify | LLM app development platform |
| fastgpt | Knowledge-based QA system |
| flowise | Drag & drop LLM flow builder |
| langflow | Visual multi-agent AI framework |
| one-api | OpenAI API management system |
| ragflow | Document understanding RAG engine |

### AI Infrastructure
| Tool | Description |
|------|-------------|
| vllm | High-throughput LLM serving |
| llama-cpp | LLM inference in C/C++ |
| comfyui | Modular Stable Diffusion GUI |
| sd-webui | Stable Diffusion web UI |
| xinference | Distributed LLM inference |

Run `getoai list` to see all 50+ available tools.

## Configuration

getoai stores its configuration in `~/.getoai/config.yaml`:

```yaml
# Default installation method preference
default_method: brew

# Custom tool registry (coming soon)
registries:
  - https://example.com/tools.json
```

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional)

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run linter
make lint
```

### Project Structure

```
getoai-cli/
├── cmd/getoai/       # Main entry point
├── internal/
│   ├── cli/          # CLI commands (cobra)
│   ├── config/       # Configuration management
│   ├── installer/    # Installation methods
│   ├── platform/     # Platform detection
│   ├── tools/        # Tool registry
│   └── util/         # Utilities (spinner, etc.)
├── Makefile
└── README.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Adding a New Tool

1. Edit `internal/tools/registry.go`
2. Add a new registration function:

```go
func registerMyTool() {
    Register(&Tool{
        Name:        "my-tool",
        Description: "Description of my tool",
        Category:    CategoryUtility,
        Website:     "https://mytool.com",
        Command:     "mytool",
        InstallMethods: map[installer.InstallMethod]InstallConfig{
            installer.MethodBrew: {Package: "my-tool"},
            installer.MethodPip:  {Package: "my-tool"},
        },
    })
}
```

3. Call the function in `init()`
4. Submit a PR

## License

MIT License - see [LICENSE](LICENSE) for details.

## Related Projects

- [Ollama](https://ollama.ai) - Run LLMs locally
- [Claude Code](https://claude.ai) - Claude AI coding assistant
- [Aider](https://aider.chat) - AI pair programming
- [Open WebUI](https://openwebui.com) - LLM web interface
