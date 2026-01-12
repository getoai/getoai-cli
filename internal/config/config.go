package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	// Proxy settings
	HttpProxy  string `json:"http_proxy,omitempty"`
	HttpsProxy string `json:"https_proxy,omitempty"`

	// Mirror settings
	NpmRegistry string `json:"npm_registry,omitempty"`
	PypiMirror  string `json:"pypi_mirror,omitempty"`
	GoProxy     string `json:"go_proxy,omitempty"`

	// Installation preferences
	PreferredMethod map[string]string `json:"preferred_method,omitempty"`

	// Install paths
	BinPath string `json:"bin_path,omitempty"`
}

var (
	configDir  string
	configFile string
	current    *Config
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	configDir = filepath.Join(home, ".config", "getoai")
	configFile = filepath.Join(configDir, "config.json")
}

func Load() (*Config, error) {
	if current != nil {
		return current, nil
	}

	current = &Config{
		PreferredMethod: make(map[string]string),
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return current, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, current); err != nil {
		return nil, err
	}

	return current, nil
}

func Save(cfg *Config) error {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}

func Get() *Config {
	if current == nil {
		_, _ = Load()
	}
	return current
}

func GetConfigPath() string {
	return configFile
}

func (c *Config) SetProxy(httpProxy, httpsProxy string) {
	c.HttpProxy = httpProxy
	c.HttpsProxy = httpsProxy
}

func (c *Config) SetNpmRegistry(registry string) {
	c.NpmRegistry = registry
}

func (c *Config) SetPypiMirror(mirror string) {
	c.PypiMirror = mirror
}

func (c *Config) SetGoProxy(proxy string) {
	c.GoProxy = proxy
}

func (c *Config) GetPreferredMethod(tool string) string {
	if c.PreferredMethod == nil {
		return ""
	}
	return c.PreferredMethod[tool]
}

func (c *Config) SetPreferredMethod(tool, method string) {
	if c.PreferredMethod == nil {
		c.PreferredMethod = make(map[string]string)
	}
	c.PreferredMethod[tool] = method
}

// ApplyEnv applies proxy and mirror settings to environment
func (c *Config) ApplyEnv() {
	if c.HttpProxy != "" {
		os.Setenv("HTTP_PROXY", c.HttpProxy)
		os.Setenv("http_proxy", c.HttpProxy)
	}
	if c.HttpsProxy != "" {
		os.Setenv("HTTPS_PROXY", c.HttpsProxy)
		os.Setenv("https_proxy", c.HttpsProxy)
	}
	if c.GoProxy != "" {
		os.Setenv("GOPROXY", c.GoProxy)
	}
}
