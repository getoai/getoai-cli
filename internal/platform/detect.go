package platform

import (
	"os/exec"
	"runtime"
	"strings"
)

type Platform struct {
	OS        string
	Arch      string
	HasBrew   bool
	HasApt    bool
	HasYum    bool
	HasDnf    bool
	HasPacman bool
	HasChoco  bool
	HasScoop  bool
	HasNpm    bool
	HasPip    bool
	HasPip3   bool
	HasDocker bool
	HasGo     bool
	HasCurl   bool
	HasWget   bool
	HomeDir   string
	IsWSL     bool
}

var current *Platform

func Detect() *Platform {
	if current != nil {
		return current
	}

	p := &Platform{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	p.HasBrew = commandExists("brew")
	p.HasApt = commandExists("apt-get")
	p.HasYum = commandExists("yum")
	p.HasDnf = commandExists("dnf")
	p.HasPacman = commandExists("pacman")
	p.HasChoco = commandExists("choco")
	p.HasScoop = commandExists("scoop")
	p.HasNpm = commandExists("npm")
	p.HasPip = commandExists("pip")
	p.HasPip3 = commandExists("pip3")
	p.HasDocker = commandExists("docker")
	p.HasGo = commandExists("go")
	p.HasCurl = commandExists("curl")
	p.HasWget = commandExists("wget")
	p.IsWSL = detectWSL()

	current = p
	return p
}

// Refresh clears the cached platform detection and re-detects
func Refresh() *Platform {
	current = nil
	return Detect()
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func detectWSL() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(out)), "microsoft") ||
		strings.Contains(strings.ToLower(string(out)), "wsl")
}

func (p *Platform) GetPackageManager() string {
	switch p.OS {
	case "darwin":
		if p.HasBrew {
			return "brew"
		}
	case "linux":
		if p.HasApt {
			return "apt"
		}
		if p.HasDnf {
			return "dnf"
		}
		if p.HasYum {
			return "yum"
		}
		if p.HasPacman {
			return "pacman"
		}
	case "windows":
		if p.HasScoop {
			return "scoop"
		}
		if p.HasChoco {
			return "choco"
		}
	}
	return ""
}

func (p *Platform) IsDarwin() bool {
	return p.OS == "darwin"
}

func (p *Platform) IsLinux() bool {
	return p.OS == "linux"
}

func (p *Platform) IsWindows() bool {
	return p.OS == "windows"
}

func (p *Platform) String() string {
	return p.OS + "/" + p.Arch
}
