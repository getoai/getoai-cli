package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/getoai/getoai-cli/internal/platform"
)

type InstallMethod string

const (
	MethodScript   InstallMethod = "script"
	MethodBrew     InstallMethod = "brew"
	MethodNpm      InstallMethod = "npm"
	MethodPip      InstallMethod = "pip"
	MethodGo       InstallMethod = "go"
	MethodDocker   InstallMethod = "docker"
	MethodBinary   InstallMethod = "binary"
	MethodApt      InstallMethod = "apt"
	MethodChoco    InstallMethod = "choco"
	MethodScoop    InstallMethod = "scoop"
	MethodDownload InstallMethod = "download" // Manual download from website
)

// GetMethodDescription returns a human-friendly description for each install method
func GetMethodDescription(method InstallMethod) string {
	descriptions := map[InstallMethod]string{
		MethodBrew:     "Homebrew package manager (macOS/Linux)",
		MethodApt:      "APT package manager (Debian/Ubuntu)",
		MethodNpm:      "NPM package manager (Node.js)",
		MethodPip:      "Pip package manager (Python)",
		MethodGo:       "Go install (requires Go toolchain)",
		MethodScript:   "Installation script (curl/wget)",
		MethodDocker:   "Docker container",
		MethodBinary:   "Pre-built binary",
		MethodDownload: "Manual download and install",
		MethodChoco:    "Chocolatey package manager (Windows)",
		MethodScoop:    "Scoop package manager (Windows)",
	}

	if desc, ok := descriptions[method]; ok {
		return desc
	}
	return string(method)
}

type Installer interface {
	Install(name string, args ...string) error
	Uninstall(name string, args ...string) error
	IsAvailable() bool
	Name() string
}

type BaseInstaller struct {
	platform *platform.Platform
}

func NewBaseInstaller() *BaseInstaller {
	return &BaseInstaller{
		platform: platform.Detect(),
	}
}

func (b *BaseInstaller) RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (b *BaseInstaller) RunCommandSilent(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// Script Installer - runs shell scripts
type ScriptInstaller struct {
	*BaseInstaller
}

func NewScriptInstaller() *ScriptInstaller {
	return &ScriptInstaller{BaseInstaller: NewBaseInstaller()}
}

func (s *ScriptInstaller) Name() string { return "script" }

func (s *ScriptInstaller) IsAvailable() bool {
	return s.platform.HasCurl || s.platform.HasWget
}

func (s *ScriptInstaller) Install(url string, args ...string) error {
	if s.platform.HasCurl {
		shellArgs := []string{"-fsSL", url}
		cmd := exec.Command("curl", shellArgs...)
		pipe, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		if err := cmd.Start(); err != nil {
			return err
		}

		shell := exec.Command("sh", "-s", "--")
		shell.Stdin = pipe
		shell.Stdout = os.Stdout
		shell.Stderr = os.Stderr
		if err := shell.Run(); err != nil {
			return err
		}
		return cmd.Wait()
	}
	return fmt.Errorf("curl not available")
}

func (s *ScriptInstaller) Uninstall(name string, args ...string) error {
	return fmt.Errorf("script installer does not support uninstall")
}

// Brew Installer
type BrewInstaller struct {
	*BaseInstaller
}

func NewBrewInstaller() *BrewInstaller {
	return &BrewInstaller{BaseInstaller: NewBaseInstaller()}
}

func (b *BrewInstaller) Name() string { return "brew" }

func (b *BrewInstaller) IsAvailable() bool {
	return b.platform.HasBrew
}

func (b *BrewInstaller) Install(name string, args ...string) error {
	allArgs := append([]string{"install", name}, args...)
	return b.RunCommand("brew", allArgs...)
}

func (b *BrewInstaller) Uninstall(name string, args ...string) error {
	return b.RunCommand("brew", "uninstall", name)
}

// Apt Installer (for Linux)
type AptInstaller struct {
	*BaseInstaller
}

func NewAptInstaller() *AptInstaller {
	return &AptInstaller{BaseInstaller: NewBaseInstaller()}
}

func (a *AptInstaller) Name() string { return "apt" }

func (a *AptInstaller) IsAvailable() bool {
	return a.platform.HasApt
}

func (a *AptInstaller) Install(name string, args ...string) error {
	// Update package list first
	if err := a.RunCommand("sudo", "apt-get", "update"); err != nil {
		fmt.Println("Warning: failed to update package list")
	}
	allArgs := append([]string{"apt-get", "install", "-y", name}, args...)
	return a.RunCommand("sudo", allArgs...)
}

func (a *AptInstaller) Uninstall(name string, args ...string) error {
	return a.RunCommand("sudo", "apt-get", "remove", "-y", name)
}

// Npm Installer
type NpmInstaller struct {
	*BaseInstaller
}

func NewNpmInstaller() *NpmInstaller {
	return &NpmInstaller{BaseInstaller: NewBaseInstaller()}
}

func (n *NpmInstaller) Name() string { return "npm" }

func (n *NpmInstaller) IsAvailable() bool {
	return n.platform.HasNpm
}

func (n *NpmInstaller) Install(name string, args ...string) error {
	allArgs := append([]string{"install", "-g", name}, args...)
	return n.RunCommand("npm", allArgs...)
}

func (n *NpmInstaller) Uninstall(name string, args ...string) error {
	return n.RunCommand("npm", "uninstall", "-g", name)
}

// Pip Installer
type PipInstaller struct {
	*BaseInstaller
	usePip3 bool
}

func NewPipInstaller() *PipInstaller {
	base := NewBaseInstaller()
	return &PipInstaller{
		BaseInstaller: base,
		usePip3:       base.platform.HasPip3,
	}
}

func (p *PipInstaller) Name() string { return "pip" }

func (p *PipInstaller) IsAvailable() bool {
	return p.platform.HasPip || p.platform.HasPip3
}

func (p *PipInstaller) pipCmd() string {
	if p.usePip3 {
		return "pip3"
	}
	return "pip"
}

func (p *PipInstaller) Install(name string, args ...string) error {
	allArgs := append([]string{"install", name}, args...)
	return p.RunCommand(p.pipCmd(), allArgs...)
}

func (p *PipInstaller) Uninstall(name string, args ...string) error {
	return p.RunCommand(p.pipCmd(), "uninstall", "-y", name)
}

// Go Installer
type GoInstaller struct {
	*BaseInstaller
}

func NewGoInstaller() *GoInstaller {
	return &GoInstaller{BaseInstaller: NewBaseInstaller()}
}

func (g *GoInstaller) Name() string { return "go" }

func (g *GoInstaller) IsAvailable() bool {
	return g.platform.HasGo
}

func (g *GoInstaller) Install(name string, args ...string) error {
	allArgs := append([]string{"install", name + "@latest"}, args...)
	return g.RunCommand("go", allArgs...)
}

func (g *GoInstaller) Uninstall(name string, args ...string) error {
	return fmt.Errorf("go installer does not support uninstall, manually remove from $GOPATH/bin")
}

// Docker Installer
type DockerInstaller struct {
	*BaseInstaller
}

func NewDockerInstaller() *DockerInstaller {
	return &DockerInstaller{BaseInstaller: NewBaseInstaller()}
}

func (d *DockerInstaller) Name() string { return "docker" }

func (d *DockerInstaller) IsAvailable() bool {
	return d.platform.HasDocker
}

// CheckDockerAvailable checks if Docker is available and running
func CheckDockerAvailable() bool {
	// Check if docker command exists
	if _, err := exec.LookPath("docker"); err != nil {
		fmt.Println()
		fmt.Println("\033[33mDocker is not installed.\033[0m")
		fmt.Println()
		fmt.Println("Install Docker first:")
		fmt.Println("  getoai install docker")
		fmt.Println()
		fmt.Println("Or install manually from: https://www.docker.com")
		fmt.Println()
		return false
	}

	// Check if Docker daemon is running
	out, err := exec.Command("docker", "info").CombinedOutput()
	if err != nil {
		fmt.Println()
		fmt.Println("\033[33mDocker is installed but not running.\033[0m")
		fmt.Println()
		// Check for common error messages
		outStr := string(out)
		if strings.Contains(outStr, "Cannot connect to the Docker daemon") ||
			strings.Contains(outStr, "Is the docker daemon running") ||
			strings.Contains(outStr, "permission denied") {
			fmt.Println("Please start Docker Desktop or the Docker service:")
			fmt.Println()
			fmt.Println("  macOS/Windows: Start Docker Desktop application")
			fmt.Println("  Linux:         sudo systemctl start docker")
			fmt.Println()
		} else {
			fmt.Printf("Error: %s\n", strings.TrimSpace(outStr))
			fmt.Println()
		}
		return false
	}
	return true
}

// CheckDockerComposeAvailable checks if docker-compose is available
func CheckDockerComposeAvailable() bool {
	// Check docker compose v2
	out, _ := exec.Command("docker", "compose", "version").CombinedOutput()
	if strings.Contains(string(out), "Docker Compose") {
		return true
	}
	// Check docker-compose v1
	if _, err := exec.LookPath("docker-compose"); err == nil {
		return true
	}
	fmt.Println()
	fmt.Println("\033[33mDocker Compose is not installed.\033[0m")
	fmt.Println()
	fmt.Println("Install Docker Compose first:")
	fmt.Println("  getoai install docker-compose")
	fmt.Println()
	return false
}

func (d *DockerInstaller) Install(image string, args ...string) error {
	// Check if Docker is running first
	if !CheckDockerAvailable() {
		return fmt.Errorf("docker is required but not running")
	}

	// Pull the image
	allArgs := append([]string{"pull", image}, args...)
	if err := d.RunCommand("docker", allArgs...); err != nil {
		showDockerMirrorHelp()
		return err
	}
	return nil
}

// InstallAndRun pulls the image and runs the container in the background
func (d *DockerInstaller) InstallAndRun(image string, containerName string, ports []string, env map[string]string, volumes []string) error {
	// Check dependencies
	if !CheckDockerAvailable() {
		return fmt.Errorf("docker is required but not running")
	}

	// First pull the image
	fmt.Printf("Pulling image %s...\n", image)
	if err := d.RunCommand("docker", "pull", image); err != nil {
		showDockerMirrorHelp()
		return fmt.Errorf("failed to pull image: %w", err)
	}

	// Check if container already exists
	out, _ := d.RunCommandSilent("docker", "ps", "-a", "--filter", fmt.Sprintf("name=^%s$", containerName), "--format", "{{.Names}}")
	if strings.TrimSpace(out) == containerName {
		fmt.Printf("Container '%s' already exists. Removing...\n", containerName)
		_, _ = d.RunCommandSilent("docker", "rm", "-f", containerName)
	}

	// Build run command
	runArgs := []string{"run", "-d", "--name", containerName, "--restart", "unless-stopped"}

	// Add port mappings
	for _, port := range ports {
		runArgs = append(runArgs, "-p", port)
	}

	// Add environment variables
	for k, v := range env {
		runArgs = append(runArgs, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	// Add volume mappings
	for _, vol := range volumes {
		runArgs = append(runArgs, "-v", vol)
	}

	runArgs = append(runArgs, image)

	fmt.Printf("Starting container '%s'...\n", containerName)
	if err := d.RunCommand("docker", runArgs...); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	// Show container status
	fmt.Println()
	fmt.Printf("\033[32m✓ Container '%s' started successfully!\033[0m\n", containerName)

	// Show access URL if ports are mapped
	if len(ports) > 0 {
		fmt.Println()
		fmt.Println("Access URLs:")
		for _, port := range ports {
			parts := strings.Split(port, ":")
			if len(parts) >= 1 {
				hostPort := parts[0]
				fmt.Printf("  http://localhost:%s\n", hostPort)
			}
		}
	}

	fmt.Println()
	fmt.Println("Useful commands:")
	fmt.Printf("  docker logs %s      View logs\n", containerName)
	fmt.Printf("  docker stop %s      Stop container\n", containerName)
	fmt.Printf("  docker start %s     Start container\n", containerName)
	fmt.Printf("  docker rm -f %s     Remove container\n", containerName)

	return nil
}

// InstallWithCompose clones the repo and starts with docker-compose
func (d *DockerInstaller) InstallWithCompose(repoURL string, appName string) error {
	// Check dependencies
	if !CheckDockerAvailable() {
		return fmt.Errorf("docker is required but not running")
	}
	if !CheckDockerComposeAvailable() {
		return fmt.Errorf("docker-compose is required but not available")
	}

	// Get install directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	installDir := fmt.Sprintf("%s/.getoai/tools/%s", homeDir, appName)

	// Check if already installed
	if _, err := os.Stat(installDir); err == nil {
		fmt.Printf("Directory %s already exists.\n", installDir)
		fmt.Println("Updating and restarting...")

		// Pull latest changes
		if err := d.RunCommand("git", "-C", installDir, "pull"); err != nil {
			fmt.Printf("Warning: failed to pull updates: %v\n", err)
		}
	} else {
		// Clone the repository
		fmt.Printf("Cloning %s...\n", repoURL)

		// Create parent directory
		parentDir := fmt.Sprintf("%s/.getoai/tools", homeDir)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		if err := d.RunCommand("git", "clone", repoURL, installDir); err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	}

	// Find docker-compose file
	composeFile := findComposeFile(installDir)
	if composeFile == "" {
		fmt.Println()
		fmt.Printf("\033[33mNo docker-compose file found in %s\033[0m\n", installDir)
		fmt.Println("Please check the repository documentation for deployment instructions.")
		fmt.Printf("Repository: %s\n", repoURL)
		return nil
	}

	// Get compose file directory
	composeDir := composeFile[:strings.LastIndex(composeFile, "/")]

	// Check for .env.example and copy to .env if .env doesn't exist
	envExample := composeDir + "/.env.example"
	envFile := composeDir + "/.env"
	if _, err := os.Stat(envExample); err == nil {
		if _, err := os.Stat(envFile); os.IsNotExist(err) {
			fmt.Println("Creating .env file from .env.example...")
			if err := copyFile(envExample, envFile); err != nil {
				fmt.Printf("Warning: failed to copy .env.example: %v\n", err)
			}
		}
	}

	// Start with docker-compose
	fmt.Printf("Starting %s with docker-compose...\n", appName)

	// Check for docker compose v2 vs v1
	var composeCmd string
	var composeArgs []string
	if _, err := exec.LookPath("docker"); err == nil {
		// Try docker compose (v2) first
		out, _ := exec.Command("docker", "compose", "version").CombinedOutput()
		if strings.Contains(string(out), "Docker Compose") {
			composeCmd = "docker"
			composeArgs = []string{"compose", "-f", composeFile, "up", "-d"}
		}
	}
	if composeCmd == "" {
		// Fallback to docker-compose (v1)
		if _, err := exec.LookPath("docker-compose"); err == nil {
			composeCmd = "docker-compose"
			composeArgs = []string{"-f", composeFile, "up", "-d"}
		}
	}

	if composeCmd == "" {
		return fmt.Errorf("docker-compose is not installed. Please install it first")
	}

	// Change to compose directory and run
	cmd := exec.Command(composeCmd, composeArgs...)
	cmd.Dir = composeDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		// Check if it's a network/timeout error and provide helpful message
		showDockerMirrorHelp()
		return fmt.Errorf("failed to start containers: %w", err)
	}

	// Show success message
	fmt.Println()
	fmt.Printf("\033[32m✓ %s started successfully!\033[0m\n", appName)
	fmt.Println()
	fmt.Printf("Install location: %s\n", installDir)
	fmt.Println()
	fmt.Println("Useful commands:")
	fmt.Printf("  cd %s && docker compose ps      View containers\n", composeDir)
	fmt.Printf("  cd %s && docker compose logs    View logs\n", composeDir)
	fmt.Printf("  cd %s && docker compose down    Stop services\n", composeDir)
	fmt.Printf("  cd %s && docker compose up -d   Start services\n", composeDir)

	return nil
}

// findComposeFile looks for docker-compose file in common locations
func findComposeFile(baseDir string) string {
	// Common locations for docker-compose files
	locations := []string{
		"docker/docker-compose.yaml",
		"docker/docker-compose.yml",
		"docker-compose.yaml",
		"docker-compose.yml",
		"compose.yaml",
		"compose.yml",
		"deploy/docker-compose.yaml",
		"deploy/docker-compose.yml",
	}

	for _, loc := range locations {
		path := fmt.Sprintf("%s/%s", baseDir, loc)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

// showDockerMirrorHelp displays instructions for configuring Docker registry mirrors
func showDockerMirrorHelp() {
	fmt.Println()
	fmt.Println("\033[33m╭─────────────────────────────────────────────────────────────────╮\033[0m")
	fmt.Println("\033[33m│ Docker 拉取镜像超时？请配置镜像加速器                              │\033[0m")
	fmt.Println("\033[33m╰─────────────────────────────────────────────────────────────────╯\033[0m")
	fmt.Println()
	fmt.Println("编辑 Docker 配置文件:")
	fmt.Println()
	fmt.Println("  \033[36m# Linux/macOS\033[0m")
	fmt.Println("  sudo mkdir -p /etc/docker")
	fmt.Println("  sudo tee /etc/docker/daemon.json <<EOF")
	fmt.Println("  {")
	fmt.Println("    \"registry-mirrors\": [")
	fmt.Println("      \"https://docker.1ms.run\",")
	fmt.Println("      \"https://docker.xuanyuan.me\"")
	fmt.Println("    ]")
	fmt.Println("  }")
	fmt.Println("  EOF")
	fmt.Println("  sudo systemctl restart docker  # Linux")
	fmt.Println("  # macOS: 重启 Docker Desktop")
	fmt.Println()
	fmt.Println("或者使用 Docker Desktop 设置:")
	fmt.Println("  Settings -> Docker Engine -> 添加 registry-mirrors")
	fmt.Println()
	fmt.Println("常用镜像加速器:")
	fmt.Println("  • https://docker.1ms.run")
	fmt.Println("  • https://docker.xuanyuan.me")
	fmt.Println("  • https://dockerhub.icu")
	fmt.Println("  • https://hub.rat.dev")
	fmt.Println()
	fmt.Println("配置完成后重新运行安装命令。")
	fmt.Println()
}

func (d *DockerInstaller) Uninstall(image string, args ...string) error {
	return d.RunCommand("docker", "rmi", image)
}

// StopContainer stops and removes a container by name
func (d *DockerInstaller) StopContainer(containerName string) error {
	_, _ = d.RunCommandSilent("docker", "stop", containerName)
	return d.RunCommand("docker", "rm", containerName)
}

// UninstallCompose stops containers but keeps the install directory
func (d *DockerInstaller) UninstallCompose(installDir string) error {
	// Find docker-compose file
	composeFile := findComposeFile(installDir)
	if composeFile != "" {
		composeDir := composeFile[:strings.LastIndex(composeFile, "/")]

		// Check for docker compose v2 vs v1
		var composeCmd string
		var composeArgs []string
		if _, err := exec.LookPath("docker"); err == nil {
			out, _ := exec.Command("docker", "compose", "version").CombinedOutput()
			if strings.Contains(string(out), "Docker Compose") {
				composeCmd = "docker"
				composeArgs = []string{"compose", "-f", composeFile, "down"}
			}
		}
		if composeCmd == "" {
			if _, err := exec.LookPath("docker-compose"); err == nil {
				composeCmd = "docker-compose"
				composeArgs = []string{"-f", composeFile, "down"}
			}
		}

		if composeCmd != "" {
			fmt.Println("Stopping containers...")
			cmd := exec.Command(composeCmd, composeArgs...)
			cmd.Dir = composeDir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to stop containers: %w", err)
			}
		}
	}

	// Show manual cleanup instructions
	fmt.Println()
	fmt.Println("Containers stopped.")
	fmt.Printf("Data directory: %s\n", installDir)
	fmt.Println()
	fmt.Println("To completely remove (including all data):")
	fmt.Println("  macOS:  Move to Trash manually or use Finder")
	fmt.Println("  Linux:  trash-put or move to ~/.local/share/Trash/")
	fmt.Println()

	return nil
}

// DownloadInstaller - shows download instructions for desktop apps
type DownloadInstaller struct {
	*BaseInstaller
}

func NewDownloadInstaller() *DownloadInstaller {
	return &DownloadInstaller{BaseInstaller: NewBaseInstaller()}
}

func (d *DownloadInstaller) Name() string { return "download" }

func (d *DownloadInstaller) IsAvailable() bool {
	return true // Always available
}

func (d *DownloadInstaller) Install(url string, args ...string) error {
	appName := ""
	downloadURL := ""
	fileType := ""

	// Parse args: [appName, downloadURL, fileType]
	if len(args) > 0 {
		appName = args[0]
	}
	if len(args) > 1 {
		downloadURL = args[1]
	}
	if len(args) > 2 {
		fileType = args[2]
	}

	fmt.Println()
	fmt.Printf("\033[36m%s is a desktop application.\033[0m\n", appName)
	fmt.Println()

	// If no specific download URL provided, fallback to browser
	if downloadURL == "" {
		fmt.Println("Please download and install from:")
		fmt.Printf("  \033[4m%s\033[0m\n", url)
		fmt.Println()
		return d.openBrowser(url)
	}

	// Download and install
	return d.downloadAndInstall(appName, downloadURL, fileType)
}

func (d *DownloadInstaller) openBrowser(url string) error {
	var cmd *exec.Cmd
	switch d.platform.OS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	}

	if cmd != nil {
		fmt.Print("Opening download page in browser... ")
		if err := cmd.Run(); err == nil {
			fmt.Println("done")
		} else {
			fmt.Println("failed (please open the URL manually)")
		}
	}

	fmt.Println()
	return nil
}

func (d *DownloadInstaller) downloadAndInstall(appName, downloadURL, fileType string) error {
	// Create temp directory
	tmpDir := os.TempDir()
	fileName := getFileNameFromURL(downloadURL)
	filePath := fmt.Sprintf("%s/%s", tmpDir, fileName)

	// Download file
	fmt.Printf("Downloading %s...\n", fileName)
	fmt.Println()

	if err := d.downloadFile(downloadURL, filePath); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	fmt.Println()

	// Install based on file type
	return d.installFile(filePath, fileType, appName)
}

func (d *DownloadInstaller) downloadFile(url, filepath string) error {
	// Use curl with progress bar and follow redirects
	cmd := exec.Command("curl", "-L", "--progress-bar", "-o", filepath, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// Verify file was downloaded
	stat, err := os.Stat(filepath)
	if err != nil {
		return fmt.Errorf("downloaded file not found: %w", err)
	}

	// Check if file is too small (likely an error page)
	if stat.Size() < 1024 {
		content, _ := os.ReadFile(filepath)
		return fmt.Errorf("downloaded file too small (%d bytes), possible error: %s", stat.Size(), string(content))
	}

	fmt.Printf("Downloaded %s (%.2f MB)\n", filepath, float64(stat.Size())/1024/1024)

	return nil
}

func (d *DownloadInstaller) installFile(filePath, fileType, appName string) error {
	switch fileType {
	case "dmg":
		return d.installDMG(filePath, appName)
	case "pkg":
		return d.installPKG(filePath)
	case "deb":
		return d.installDEB(filePath)
	case "appimage":
		return d.installAppImage(filePath, appName)
	case "exe":
		return d.installEXE(filePath)
	case "msi":
		return d.installMSI(filePath)
	default:
		return fmt.Errorf("unsupported file type: %s", fileType)
	}
}

func (d *DownloadInstaller) installDMG(dmgPath, appName string) error {
	fmt.Println("Installing DMG file...")

	// Mount the DMG
	fmt.Println("Mounting DMG...")
	out, err := exec.Command("hdiutil", "attach", dmgPath, "-nobrowse", "-mountrandom", "/tmp").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to mount DMG: %w\n%s", err, string(out))
	}

	outStr := string(out)
	fmt.Printf("Mount output: %s\n", outStr)

	// Parse mount point from output - get the last mount point in the output
	mountPoint := parseMountPointFromHdiutil(outStr)
	if mountPoint == "" {
		return fmt.Errorf("failed to find mount point in output: %s", outStr)
	}

	fmt.Printf("Mount point: %s\n", mountPoint)

	defer func() {
		fmt.Println("Unmounting DMG...")
		_ = exec.Command("hdiutil", "detach", mountPoint, "-force").Run()
	}()

	// Find .app in mount point
	entries, err := os.ReadDir(mountPoint)
	if err != nil {
		return fmt.Errorf("failed to read mount point %s: %w", mountPoint, err)
	}

	var appPath string
	for _, entry := range entries {
		fmt.Printf("Found entry: %s (isDir: %v)\n", entry.Name(), entry.IsDir())
		if strings.HasSuffix(entry.Name(), ".app") {
			appPath = fmt.Sprintf("%s/%s", mountPoint, entry.Name())
			break
		}
	}

	if appPath == "" {
		return fmt.Errorf("no .app file found in DMG, entries: %v", entries)
	}

	// Copy to Applications
	destPath := fmt.Sprintf("/Applications/%s", strings.TrimPrefix(appPath, mountPoint+"/"))
	fmt.Printf("Installing to: %s\n", destPath)

	// Remove existing app if present
	if _, err := os.Stat(destPath); err == nil {
		fmt.Println("Removing existing installation...")
		if err := os.RemoveAll(destPath); err != nil {
			return fmt.Errorf("failed to remove existing app: %w", err)
		}
	}

	if err := exec.Command("cp", "-R", appPath, "/Applications/").Run(); err != nil {
		return fmt.Errorf("failed to copy app: %w", err)
	}

	fmt.Println()
	fmt.Printf("\033[32m✓ %s installed successfully to /Applications\033[0m\n", appName)
	fmt.Println()

	return nil
}

func (d *DownloadInstaller) installPKG(pkgPath string) error {
	fmt.Println("Installing PKG file...")
	fmt.Println("This requires administrator privileges.")
	fmt.Println()

	if err := exec.Command("sudo", "installer", "-pkg", pkgPath, "-target", "/").Run(); err != nil {
		return fmt.Errorf("failed to install PKG: %w", err)
	}

	fmt.Println()
	fmt.Println("\033[32m✓ Installation completed\033[0m")
	fmt.Println()

	return nil
}

func (d *DownloadInstaller) installDEB(debPath string) error {
	fmt.Println("Installing DEB package...")
	fmt.Println("This requires administrator privileges.")
	fmt.Println()

	if err := exec.Command("sudo", "dpkg", "-i", debPath).Run(); err != nil {
		// Try to fix dependencies
		fmt.Println("Fixing dependencies...")
		_ = exec.Command("sudo", "apt-get", "install", "-f", "-y").Run()
		return fmt.Errorf("failed to install DEB: %w", err)
	}

	fmt.Println()
	fmt.Println("\033[32m✓ Installation completed\033[0m")
	fmt.Println()

	return nil
}

func (d *DownloadInstaller) installAppImage(appImagePath, appName string) error {
	fmt.Println("Installing AppImage...")

	// Make executable
	if err := os.Chmod(appImagePath, 0755); err != nil {
		return fmt.Errorf("failed to make executable: %w", err)
	}

	// Move to user's local bin
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	destDir := fmt.Sprintf("%s/.local/bin", homeDir)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	destPath := fmt.Sprintf("%s/%s.appimage", destDir, appName)
	if err := os.Rename(appImagePath, destPath); err != nil {
		return fmt.Errorf("failed to move AppImage: %w", err)
	}

	fmt.Println()
	fmt.Printf("\033[32m✓ AppImage installed to: %s\033[0m\n", destPath)
	fmt.Printf("Make sure %s is in your PATH\n", destDir)
	fmt.Println()

	return nil
}

func (d *DownloadInstaller) installEXE(exePath string) error {
	fmt.Println("Launching installer...")
	fmt.Println()

	if err := exec.Command(exePath).Run(); err != nil {
		return fmt.Errorf("failed to launch installer: %w", err)
	}

	return nil
}

func (d *DownloadInstaller) installMSI(msiPath string) error {
	fmt.Println("Installing MSI package...")
	fmt.Println()

	if err := exec.Command("msiexec", "/i", msiPath).Run(); err != nil {
		return fmt.Errorf("failed to install MSI: %w", err)
	}

	return nil
}

func parseMountPointFromHdiutil(output string) string {
	// hdiutil output format:
	// /dev/disk4s2        	Apple_HFS                      	/private/tmp/dmg.XXXXXX
	// We want the last column which is the mount point
	lines := strings.Split(output, "\n")
	var mountPoint string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Look for lines with Apple_HFS or similar filesystem types
		if strings.Contains(line, "Apple_HFS") || strings.Contains(line, "Apple_APFS") {
			// Split by tabs or multiple spaces
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				// The last part should be the mount point
				mountPoint = parts[len(parts)-1]
			}
		}

		// Also check for /tmp or /Volumes paths
		if (strings.Contains(line, "/tmp/") || strings.Contains(line, "/Volumes/")) && !strings.Contains(line, "/dev/") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasPrefix(part, "/tmp/") || strings.HasPrefix(part, "/Volumes/") {
					mountPoint = part
				}
			}
		}
	}

	return mountPoint
}

func getFileNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "download"
}

func (d *DownloadInstaller) Uninstall(name string, args ...string) error {
	appName := ""
	if len(args) > 0 {
		appName = args[0]
	}

	// If no app name provided, show manual instructions
	if appName == "" {
		fmt.Println()
		fmt.Printf("Please uninstall %s manually:\n", name)
		fmt.Println("  macOS: Move to Trash from Applications folder")
		fmt.Println("  Linux: Use your package manager or remove manually")
		fmt.Println("  Windows: Use Add/Remove Programs")
		fmt.Println()
		return nil
	}

	// Attempt to uninstall based on platform
	switch d.platform.OS {
	case "darwin":
		return d.uninstallMacOS(appName)
	case "linux":
		return d.uninstallLinux(name)
	case "windows":
		return d.uninstallWindows(appName)
	}

	return fmt.Errorf("unsupported platform for automatic uninstall")
}

func (d *DownloadInstaller) uninstallMacOS(appName string) error {
	// Check /Applications first
	appPath := fmt.Sprintf("/Applications/%s", appName)
	if _, err := os.Stat(appPath); err == nil {
		fmt.Printf("Removing %s...\n", appPath)
		if err := os.RemoveAll(appPath); err != nil {
			return fmt.Errorf("failed to remove app: %w", err)
		}
		fmt.Println("\033[32m✓ Application removed successfully\033[0m")
		return nil
	}

	// Check ~/Applications
	homeDir, _ := os.UserHomeDir()
	userAppPath := fmt.Sprintf("%s/Applications/%s", homeDir, appName)
	if _, err := os.Stat(userAppPath); err == nil {
		fmt.Printf("Removing %s...\n", userAppPath)
		if err := os.RemoveAll(userAppPath); err != nil {
			return fmt.Errorf("failed to remove app: %w", err)
		}
		fmt.Println("\033[32m✓ Application removed successfully\033[0m")
		return nil
	}

	return fmt.Errorf("application not found in /Applications or ~/Applications")
}

func (d *DownloadInstaller) uninstallLinux(name string) error {
	// Try to uninstall via package manager first
	if _, err := exec.LookPath("dpkg"); err == nil {
		fmt.Println("Attempting to remove via dpkg...")
		if err := exec.Command("sudo", "dpkg", "-r", name).Run(); err == nil {
			fmt.Println("\033[32m✓ Package removed successfully\033[0m")
			return nil
		}
	}

	// Check for AppImage
	homeDir, _ := os.UserHomeDir()
	appImagePath := fmt.Sprintf("%s/.local/bin/%s.appimage", homeDir, name)
	if _, err := os.Stat(appImagePath); err == nil {
		fmt.Printf("Removing %s...\n", appImagePath)
		if err := os.Remove(appImagePath); err != nil {
			return fmt.Errorf("failed to remove AppImage: %w", err)
		}
		fmt.Println("\033[32m✓ AppImage removed successfully\033[0m")
		return nil
	}

	return fmt.Errorf("application not found or unable to uninstall")
}

func (d *DownloadInstaller) uninstallWindows(appName string) error {
	fmt.Println()
	fmt.Println("Please uninstall via Windows Settings:")
	fmt.Println("  Settings -> Apps -> Apps & features")
	fmt.Printf("  Find and uninstall: %s\n", appName)
	fmt.Println()
	return nil
}

// GetInstaller returns the appropriate installer for the given method
func GetInstaller(method InstallMethod) (Installer, error) {
	var inst Installer

	switch method {
	case MethodScript:
		inst = NewScriptInstaller()
	case MethodBrew:
		inst = NewBrewInstaller()
	case MethodApt:
		inst = NewAptInstaller()
	case MethodNpm:
		inst = NewNpmInstaller()
	case MethodPip:
		inst = NewPipInstaller()
	case MethodGo:
		inst = NewGoInstaller()
	case MethodDocker:
		inst = NewDockerInstaller()
	case MethodDownload:
		inst = NewDownloadInstaller()
	default:
		return nil, fmt.Errorf("unknown install method: %s", method)
	}

	if !inst.IsAvailable() {
		return nil, fmt.Errorf("%s is not available on this system", inst.Name())
	}

	return inst, nil
}

// CheckInstalled checks if a command is installed
func CheckInstalled(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// RunCommandSilent runs a command and returns output without printing
func RunCommandSilent(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// GetVersion tries to get the version of an installed command
func GetVersion(cmd string) string {
	for _, flag := range []string{"--version", "-v", "version"} {
		out, err := exec.Command(cmd, flag).CombinedOutput()
		if err == nil {
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			if len(lines) > 0 {
				return lines[0]
			}
		}
	}
	return "unknown"
}
