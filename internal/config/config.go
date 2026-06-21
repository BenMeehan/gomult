package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig          `yaml:"server"`
	Sandbox   SandboxConfig         `yaml:"sandbox"`
	Languages map[string]LangConfig `yaml:"languages"`
}

type ServerConfig struct {
	Port        int    `yaml:"port"`
	ReadTimeout string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
	MaxCodeSize int64  `yaml:"max_code_size"`
}

type SandboxConfig struct {
	Mode          string `yaml:"mode"`
	NsjailPath    string `yaml:"nsjail_path"`
	TimeLimit     int    `yaml:"time_limit"`
	MaxMemory     int64  `yaml:"max_memory"`
	MaxFileSize   int64  `yaml:"max_file_size"`
	MaxProcesses  int    `yaml:"max_processes"`
	RuntimeUser   int    `yaml:"runtime_user"`
	RuntimeGroup  int    `yaml:"runtime_group"`
}

type LangConfig struct {
	Name       string   `yaml:"name"`
	Extension  string   `yaml:"extension"`
	CompileCmd []string `yaml:"compile_cmd"`
	RunCmd     []string `yaml:"run_cmd"`
	SourceFile string   `yaml:"source_file"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.ReadTimeout == "" {
		cfg.Server.ReadTimeout = "10s"
	}
	if cfg.Server.WriteTimeout == "" {
		cfg.Server.WriteTimeout = "15s"
	}
	if cfg.Server.MaxCodeSize == 0 {
		cfg.Server.MaxCodeSize = 1 << 20
	}
	if cfg.Sandbox.Mode == "" {
		cfg.Sandbox.Mode = "auto"
	}
	if cfg.Sandbox.NsjailPath == "" {
		cfg.Sandbox.NsjailPath = "nsjail"
	}
	if cfg.Sandbox.TimeLimit == 0 {
		cfg.Sandbox.TimeLimit = 5
	}
	if cfg.Sandbox.MaxMemory == 0 {
		cfg.Sandbox.MaxMemory = 256 << 20
	}
	if cfg.Sandbox.MaxFileSize == 0 {
		cfg.Sandbox.MaxFileSize = 10 << 20
	}
	if cfg.Sandbox.MaxProcesses == 0 {
		cfg.Sandbox.MaxProcesses = 32
	}

	return cfg, nil
}

func (c *ServerConfig) ReadTimeoutDuration() time.Duration {
	d, _ := time.ParseDuration(c.ReadTimeout)
	return d
}

func (c *ServerConfig) WriteTimeoutDuration() time.Duration {
	d, _ := time.ParseDuration(c.WriteTimeout)
	return d
}

func (c *SandboxConfig) TimeLimitDuration() time.Duration {
	return time.Duration(c.TimeLimit) * time.Second
}
