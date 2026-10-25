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

// ResolveDefaultDatabasePath walks up from cwd and the executable dir looking for data/proxy.db.
func ResolveDefaultDatabasePath() (string, error) {
	if cwd, err := os.Getwd(); err == nil {
		if p, ok := findProxyDBFrom(cwd); ok {
			return p, nil
		}
	}
	exe, err := os.Executable()
	if err == nil {
		if p, ok := findProxyDBFrom(filepath.Dir(exe)); ok {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find data/proxy.db; set storage.database_path or AI_PROXY_DESKTOP_DATABASE_PATH")
}

func findProxyDBFrom(start string) (string, bool) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", false
	}
	for i := 0; i < 12; i++ {
		candidate := filepath.Join(dir, "data", "proxy.db")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}

func IsUniqueConstraint(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "constraint")
}
