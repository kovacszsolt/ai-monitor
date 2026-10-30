package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Storage StorageConfig `yaml:"storage"`
}

type StorageConfig struct {
	DatabasePath  string `yaml:"database_path"`
	JournalMode   string `yaml:"journal_mode"`
	BusyTimeoutMS int    `yaml:"busy_timeout_ms"`
	Synchronous   string `yaml:"synchronous"`
}

func Default() Config {
	return Config{
		Storage: StorageConfig{
			DatabasePath:  "",
			JournalMode:   "WAL",
			BusyTimeoutMS: 5000,
			Synchronous:   "NORMAL",
		},
	}
}

func Load() (Config, error) {
	cfg := Default()
	applyEnv(&cfg)

	path := ResolveConfigFilePath()
	readPath := path
	if path != "" {
		if _, err := os.Stat(readPath); err != nil {
			example := filepath.Join(filepath.Dir(path), "config.example.yaml")
			if _, err2 := os.Stat(example); err2 == nil {
				readPath = example
			}
		}
	}
	if readPath != "" {
		data, err := os.ReadFile(readPath)
		if err != nil {
			return Config{}, fmt.Errorf("read config %s: %w", path, err)
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return Config{}, fmt.Errorf("parse config: %w", err)
		}
		applyEnv(&cfg)
	}

	if cfg.Storage.DatabasePath == "" {
		resolved, err := ResolveDefaultDatabasePath()
		if err != nil {
			return Config{}, err
		}
		cfg.Storage.DatabasePath = resolved
	}
	if v := os.Getenv("AI_PROXY_DESKTOP_DATABASE_PATH"); v != "" {
		cfg.Storage.DatabasePath = v
	}
	return cfg, nil
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("AI_PROXY_DESKTOP_DATABASE_PATH"); v != "" {
		cfg.Storage.DatabasePath = v
	}
}

// ResolveDefaultDatabasePath returns the standard per-user database location (~/.config/ai-proxy-manager/proxy.db).
// Honors XDG_CONFIG_HOME when set.
func ResolveDefaultDatabasePath() (string, error) {
	base, err := userConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve default database path: %w", err)
	}
	return filepath.Join(base, "ai-proxy-manager", "proxy.db"), nil
}

func userConfigDir() (string, error) {
	if v := strings.TrimSpace(os.Getenv("XDG_CONFIG_HOME")); v != "" {
		return v, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config"), nil
}

func IsUniqueConstraint(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "constraint")
}
