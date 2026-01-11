# getoai-cli

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/getoai/getoai-cli)](https://github.com/getoai/getoai-cli/releases)

**AI 工具一键安装器**

[English](README.md) | [中文](README_zh.md)

getoai-cli 是一个跨平台的命令行工具，用于安装和管理 AI 相关的工具和命令行程序。无需再搜索安装教程，只需 `getoai install <tool>` 即可。

## 特性

- **50+ AI 工具** - 大模型运行器、编程助手、聊天界面、AI 平台等
- **多种安装方式** - Homebrew、npm、pip、Go、Docker、脚本安装
- **跨平台支持** - macOS (Intel & Apple Silicon)、Linux (amd64 & arm64)、Windows
- **智能检测** - 自动选择最适合您系统的安装方式
- **轻松管理** - 安装、列表、搜索、查看工具信息

## 快速开始

### 安装

**macOS / Linux:**

```bash
# 使用 curl
curl -fsSL https://raw.githubusercontent.com/getoai/getoai-cli/master/install.sh | bash

# 或使用 Homebrew (即将支持)
brew install getoai/tap/getoai
```

**从源码安装:**

```bash
go install github.com/getoai/getoai-cli/cmd/getoai@latest
```

**从 Releases 下载:**

从 [Releases](https://github.com/getoai/getoai-cli/releases) 页面下载适合您平台的二进制文件。

### 使用方法

```bash
# 列出所有可用工具
getoai list

# 安装工具
getoai install ollama
getoai install claude-code
getoai install aider

# 搜索工具
getoai search "coding"

# 显示工具信息
getoai info ollama

# 列出已安装的工具
getoai installed
```

## 支持的工具

### 大模型运行器
| 工具 | 描述 |
|------|------|
| ollama | 本地运行大语言模型 |
| localai | 免费开源的 OpenAI 替代品 |
| chatglm | 开源双语对话语言模型 |

### 编程助手
| 工具 | 描述 |
|------|------|
| claude-code | Claude AI 编程助手 CLI |
| aider | 终端中的 AI 结对编程 |
| cursor | AI 优先的代码编辑器 |
| gh-copilot | GitHub Copilot 命令行版 |
| tabby | 自托管 AI 编程助手 |
| gpt-engineer | AI 根据提示构建项目 |
| codex-cli | OpenAI Codex CLI |
| interpreter | 开源代码解释器 |

### 聊天界面
| 工具 | 描述 |
|------|------|
| open-webui | 友好的 LLM Web 界面 |
| lobechat | 现代化 ChatGPT/LLM 界面 |
| chatgpt-next-web | 跨平台 ChatGPT 界面 |
| chatbox | LLM 桌面客户端 |
| jan | 开源 ChatGPT 替代品 |
| librechat | 增强版 ChatGPT 克隆 |

### 命令行工具
| 工具 | 描述 |
|------|------|
| llm | 命令行访问 LLM |
| sgpt | Shell GPT - 终端中的 AI |
| mods | Charm 出品的命令行 AI |
| tgpt | 无需 API Key 的 AI 聊天 |
| fabric | AI 增强人类的框架 |
| aichat | 多合一 AI CLI 工具 |
| gemini-cli | 终端中的 Google Gemini |

### AI 平台
| 工具 | 描述 |
|------|------|
| dify | LLM 应用开发平台 |
| fastgpt | 基于知识库的问答系统 |
| flowise | 拖拽式 LLM 流程构建器 |
| langflow | 可视化多智能体框架 |
| one-api | OpenAI API 管理系统 |
| ragflow | 文档理解 RAG 引擎 |

### AI 基础设施
| 工具 | 描述 |
|------|------|
| vllm | 高吞吐量 LLM 推理服务 |
| llama-cpp | C/C++ 实现的 LLM 推理 |
| comfyui | 模块化 Stable Diffusion GUI |
| sd-webui | Stable Diffusion Web UI |
| xinference | 分布式 LLM 推理框架 |

运行 `getoai list` 查看全部 50+ 可用工具。

## 配置

getoai 的配置文件存储在 `~/.getoai/config.yaml`:

```yaml
# 默认安装方式偏好
default_method: brew

# 自定义工具注册表 (即将支持)
registries:
  - https://example.com/tools.json
```

## 开发

### 环境要求

- Go 1.21 或更高版本
- Make (可选)

### 构建

```bash
# 构建当前平台
make build

# 构建所有平台
make build-all

# 运行测试
make test

# 运行 linter
make lint
```

### 项目结构

```
getoai-cli/
├── cmd/getoai/       # 主入口
├── internal/
│   ├── cli/          # CLI 命令 (cobra)
│   ├── config/       # 配置管理
│   ├── installer/    # 安装方法
│   ├── platform/     # 平台检测
│   ├── tools/        # 工具注册表
│   └── util/         # 工具函数 (spinner 等)
├── Makefile
└── README.md
```

## 贡献

欢迎贡献代码！请随时提交 Pull Request。

### 添加新工具

1. 编辑 `internal/tools/registry.go`
2. 添加新的注册函数:

```go
func registerMyTool() {
    Register(&Tool{
        Name:        "my-tool",
        Description: "工具描述",
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

3. 在 `init()` 中调用该函数
4. 提交 PR

## 许可证

MIT 许可证 - 详见 [LICENSE](LICENSE)

## 相关项目

- [Ollama](https://ollama.ai) - 本地运行 LLM
- [Claude Code](https://claude.ai) - Claude AI 编程助手
- [Aider](https://aider.chat) - AI 结对编程
- [Open WebUI](https://openwebui.com) - LLM Web 界面
