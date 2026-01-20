# Windows Testing Guide

This guide provides comprehensive testing instructions for verifying getoai-cli functionality on Windows 10/11.

## Prerequisites

### Required
- Windows 10/11 (64-bit)
- PowerShell (Administrator privileges required for initial setup)
- At least one Windows package manager installed

### Package Managers

#### Option 1: Install Chocolatey (Recommended)

Chocolatey is the most popular Windows package manager and provides the widest tool coverage.

```powershell
# Run in PowerShell (Admin)
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

Verify installation:
```powershell
choco --version
```

#### Option 2: Install Scoop (Lightweight Alternative)

Scoop is a lightweight command-line installer that doesn't require admin privileges.

```powershell
# Run in PowerShell (non-Admin is fine)
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex
```

Verify installation:
```powershell
scoop --version
```

## Test Cases

### Test 1: Platform Detection

Verify that getoai correctly detects Windows package managers.

```bash
getoai list
```

**Expected Behavior:**
- Tools should show `choco` and/or `scoop` in their available installation methods
- Windows-specific tools should appear in the list

### Test 2: Install Ollama via Chocolatey

Test installing a priority tool using Chocolatey.

```bash
getoai install ollama
```

**Expected Behavior:**
- If both Chocolatey and Scoop are installed, shows interactive menu
- Installs using Chocolatey if selected
- Shows installation progress with spinner
- Verifies installation: `ollama --version`

**Verification:**
```bash
ollama --version
```

### Test 3: Install Node.js

Test Node.js installation with proper package naming.

```bash
getoai install node
```

**Expected Behavior:**
- Chocolatey installs `nodejs.install` package
- Scoop installs `nodejs` package
- Successful installation message shows version
- npm should also be available

**Verification:**
```bash
node --version
npm --version
```

### Test 4: Install GitHub CLI

Test a straightforward tool with same package name in both managers.

```bash
getoai install gh
```

**Expected Behavior:**
- Package name `gh` used for both Chocolatey and Scoop
- Command line tool available after install

**Verification:**
```bash
gh --version
```

### Test 5: Install Docker Desktop

Test installation of a desktop application.

```bash
getoai install docker
```

**Expected Behavior:**
- Chocolatey installs `docker-desktop` package (GUI application)
- May require system restart or Docker service start
- Docker Desktop app should appear in Start Menu

**Verification:**
```bash
docker --version
docker ps
```

**Note:** Docker Desktop may require manual startup after installation.

### Test 6: Install VS Code

Test IDE installation with multiple methods available.

```bash
getoai install vscode
```

**Expected Behavior:**
- Multiple methods available: choco, scoop, download
- Chocolatey preferred by default
- VS Code command line tool `code` available

**Verification:**
```bash
code --version
```

### Test 7: Install Docker Compose

Test Chocolatey-only tool (no Scoop support).

```bash
getoai install docker-compose
```

**Expected Behavior:**
- Only Chocolatey method shown (if Scoop is sole package manager, should show clear error)
- Installs docker-compose v2.x

**Verification:**
```bash
docker-compose --version
```

### Test 8: Interactive Method Selection

Test method selection menu when multiple package managers are available.

**Prerequisites:** Both Chocolatey and Scoop installed

```bash
getoai install ollama
```

**Expected Behavior:**
1. Displays menu: "Multiple installation methods available for ollama"
2. Shows numbered options:
   - 1) choco - Chocolatey package manager for Windows
   - 2) scoop - Scoop package manager for Windows
   - 3) download - Direct download
   - 4) script - Script-based installation
3. Prompts: "Enter your choice (1-4):"
4. After selection, shows: "Selected installation method: [method]"
5. Proceeds with installation

### Test 9: Method Selection via Flag

Test --method flag to bypass interactive menu.

```bash
getoai install ollama --method choco
```

**Expected Behavior:**
- Skips interactive menu
- Directly installs using specified method
- Shows error if method not available

### Test 10: Missing Package Manager Handling

Test behavior when required package manager is missing.

**Prerequisites:** Uninstall both Chocolatey and Scoop, OR test with a tool requiring a missing dependency

```bash
getoai install ollama
```

**Expected Behavior:**
- Shows clear error message: "No installation method available for ollama on this system"
- Suggests installing Chocolatey or Scoop
- Provides link to tool website for manual installation

### Test 11: List All Tools

Verify Windows methods appear in tool listings.

```bash
getoai list
```

**Expected Behavior:**
- All 81 tools displayed
- Priority tools show choco/scoop methods
- Python/Node.js tools show pip/npm methods
- Desktop apps show download method
- Categories properly displayed

### Test 12: Uninstall Tool

Test tool removal using Windows package managers.

```bash
getoai uninstall node
```

**Expected Behavior:**
- Prompts for confirmation: "Are you sure you want to uninstall node? [y/N]"
- Uninstalls via Chocolatey: `choco uninstall nodejs.install -y`
- Shows success message
- Verifies removal

**Verification:**
```bash
node --version  # Should show "command not found"
```

### Test 13: Force Uninstall

Test --force flag to skip confirmation.

```bash
getoai uninstall gh --force
```

**Expected Behavior:**
- Skips confirmation prompt
- Immediately proceeds with uninstallation
- Shows success message

### Test 14: Python CLI Tools

Test pip-based tools work on Windows (requires Python installed).

```bash
# Install Python first if needed
getoai install python

# Test pip-based tool
getoai install aider
```

**Expected Behavior:**
- Python tools install via pip
- Works on Windows without special handling

**Verification:**
```bash
aider --version
```

### Test 15: Check Tool Status

Test status checking for installed and non-installed tools.

```bash
getoai status ollama
getoai status nonexistent-tool
```

**Expected Behavior:**
- Installed tools show version
- Non-installed tools show "not installed"
- Unknown tools show error with suggestions

### Test 16: Automatic Installation for Desktop Apps

Test automatic installation of desktop applications via direct download.

**Prerequisites:** No package managers needed for this test

```bash
getoai install cursor
# or
getoai install windsurf
```

**Expected Behavior:**
1. Downloads .exe installer to temp directory
2. Shows download progress with curl
3. Attempts silent installation:
   - First tries `/S` flag (NSIS)
   - Then tries `/VERYSILENT` (Inno Setup)
4. If silent fails, launches interactive installer with clear message:
   - "Silent installation not supported, launching interactive installer..."
   - "Please follow the on-screen instructions to complete installation."
5. Shows success message: "✓ Installation completed"
6. Reminds user: "You may need to restart your shell for the changes to take effect."

**For MSI packages:**
```bash
# If any tool uses MSI installer
getoai install <tool-with-msi>
```

**Expected Behavior:**
1. Uses `/passive` mode (shows progress bar, no interaction needed)
2. Falls back to `/qn` (fully silent) if passive fails
3. Falls back to interactive if silent modes fail
4. Prevents automatic system restart (`/norestart`)

**Verification:**
- Desktop app appears in Start Menu
- Command line tool available (if applicable)
- Installation completes without errors

## Expected Overall Behavior

### Method Priority on Windows

When multiple methods are available, priority order:
1. **choco** (priority 1) - Preferred for Windows
2. **scoop** (priority 2) - Lightweight alternative
3. **npm** (priority 3) - For Node.js tools
4. **pip** (priority 3) - For Python tools
5. **go** (priority 4) - For Go tools
6. **script** (priority 5) - Shell scripts (may not work on Windows)
7. **docker** (priority 6) - Docker containers
8. **download** (priority 8) - Direct download fallback

### Graceful Degradation

- If Chocolatey not installed, falls back to Scoop
- If neither package manager available, uses download method (if available)
- Clear error messages guide users to install package managers

### Platform-Specific Behavior

- **PlatformOverrides** ensure Windows uses appropriate methods
- Linux/macOS-specific methods (apt, brew) not shown on Windows
- Windows-specific packages (docker-desktop vs docker) handled correctly

### Automatic Installation for Desktop Apps

For desktop applications installed via the download method:

**EXE Installers:**
- Automatically attempts silent installation using common parameters:
  - `/S` (NSIS installers - most common)
  - `/VERYSILENT` (Inno Setup installers)
- Falls back to interactive installation if silent mode not supported
- Provides clear feedback during installation process

**MSI Installers:**
- Uses `/passive` mode by default (unattended with progress bar)
- Falls back to `/qn` (fully silent) if passive mode fails
- Falls back to interactive mode if silent installation fails
- Automatically prevents system restart during installation (`/norestart`)

**Installation Flow:**
1. Download the installer to temp directory
2. Attempt automatic/silent installation
3. Display progress and status messages
4. Clean up installer file after completion
5. Notify user if shell restart needed

This means desktop apps can be installed with a single command, without manual download and installation steps.

## Tools Working on Windows

After implementation, the following tools should work on Windows:

### Via Chocolatey/Scoop (6 tools)
- ollama
- node (nodejs.install / nodejs)
- gh
- docker (docker-desktop / docker)
- vscode
- docker-compose (Chocolatey only)

### Via pip (20+ Python tools)
- aider
- llm
- openai-cli
- gpt-engineer
- sgpt (shell-gpt)
- fabric
- And other Python-based CLI tools

### Via npm (15+ Node.js tools)
- claude-code
- codex-cli
- And other npm packages

### Via go install (5+ Go tools)
- mods
- glow
- And other Go-based tools

### Via Direct Download (8+ desktop apps)
- cursor
- lmstudio
- jan
- msty
- chatbox
- cherry-studio
- windsurf
- vscode

**Total: 40+ tools available on Windows**

## Known Limitations

### Docker on Windows
- Docker Desktop requires Windows 10/11 Pro, Enterprise, or Education (for Hyper-V)
- Windows Home users need WSL2 backend
- May require system restart after installation
- Docker daemon must be manually started

### Script-based Installations
- Shell scripts (MethodScript) may not work on Windows
- getoai will skip script methods on Windows unless they're bash scripts running under Git Bash or WSL

### Path Issues
- Some tools may require shell restart to appear in PATH
- Chocolatey modifies PATH system-wide (requires new shell)
- Scoop modifies PATH for current user (requires shell restart)

### Admin Privileges
- Chocolatey requires Administrator privileges for installation
- Scoop does not require admin (per-user installation)
- Some tools may require admin for installation

## Troubleshooting

### Issue: "choco: command not found"

**Solution:** Install Chocolatey or use Scoop
```powershell
# Install Chocolatey (see Prerequisites section)
# OR install Scoop
```

### Issue: "No installation method available"

**Solution:** Install at least one Windows package manager (Chocolatey or Scoop)

### Issue: Tool installed but command not found

**Solution:** Restart PowerShell or Command Prompt to refresh PATH
```powershell
# Close and reopen PowerShell
# OR manually refresh PATH
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
```

### Issue: Chocolatey installation requires admin

**Solution:** Run PowerShell as Administrator
```
Right-click PowerShell → "Run as Administrator"
```

### Issue: Docker Desktop won't start

**Solution:** Enable Hyper-V or WSL2
```powershell
# For Hyper-V (Windows Pro/Enterprise)
Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All

# For WSL2 (Windows Home)
wsl --install
```

## Reporting Issues

If you encounter issues during testing:

1. **Check Prerequisites:** Ensure Chocolatey or Scoop is properly installed
2. **Verify Package Names:** Check official Chocolatey/Scoop repositories for correct package names
3. **Check Logs:** Look for error messages in PowerShell output
4. **Report:** Open an issue at https://github.com/getoai/getoai-cli/issues with:
   - Windows version (run `winver`)
   - Package manager version (`choco --version` or `scoop --version`)
   - Tool name and command used
   - Full error output

## Success Criteria

Testing is successful when:

- ✅ All 6 priority tools install successfully via Chocolatey
- ✅ All 6 priority tools install successfully via Scoop (except docker-compose)
- ✅ Interactive method selection works correctly
- ✅ --method flag properly overrides method selection
- ✅ Error messages are clear and helpful
- ✅ Uninstall functionality works for Windows package managers
- ✅ Python/Node.js/Go tools work via pip/npm/go install
- ✅ Download method works as fallback
- ✅ EXE installers attempt silent installation automatically
- ✅ MSI installers use passive/silent mode by default
- ✅ Interactive installation launched if silent mode fails
- ✅ Desktop apps install successfully via download method

## Additional Testing Notes

### Cross-Platform Consistency
- Verify that Windows installation experience matches macOS (Homebrew) and Linux (APT)
- Error messages should be consistent across platforms
- Interactive prompts should behave identically

### Package Name Verification
All package names have been verified from official sources (January 2026):
- **Chocolatey:** https://community.chocolatey.org/packages
- **Scoop:** https://github.com/ScoopInstaller/Main

### Performance
- Installation should complete within reasonable time (varies by tool)
- Chocolatey typically slower than Scoop (more comprehensive)
- Download method fastest for small tools

## Next Steps After Testing

1. Document any issues found during testing
2. Verify package names are still current (packages may be updated)
3. Test on both Windows 10 and Windows 11
4. Test with different PowerShell versions (5.1 and 7.x)
5. Consider automating tests with GitHub Actions Windows runners
