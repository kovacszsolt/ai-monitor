package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ConfigFilePath returns the desktop YAML path used for load/save (absolute when possible).
func ConfigFilePath() string {
	return ResolveConfigFilePath()
}

// LoadFileSnapshot reads YAML from disk without auto-resolving database_path or applying env overrides.
func LoadFileSnapshot() (filePath string, cfg Config, err error) {
	filePath = ConfigFilePath()
	cfg = Default()
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return filePath, cfg, nil
		}
		return "", Config{}, fmt.Errorf("read config %s: %w", filePath, err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "", Config{}, fmt.Errorf("parse config %s: %w", filePath, err)
	}
	return filePath, cfg, nil
}

// WriteDatabasePath updates storage.database_path in the desktop config file (creates file if missing).
func WriteDatabasePath(configPath, databasePath string) error {
	if strings.TrimSpace(configPath) == "" {
		configPath = "config.yaml"
	}
	cfg := Default()
	if data, err := os.ReadFile(configPath); err == nil {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("parse config %s: %w", configPath, err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("read config %s: %w", configPath, err)
	}
	cfg.Storage.DatabasePath = strings.TrimSpace(databasePath)
	dir := filepath.Dir(configPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create config dir: %w", err)
		}
	}
	out, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	header := "# AI Proxy Manager desktop config\n"
	if err := os.WriteFile(configPath, append([]byte(header), out...), 0o644); err != nil {
		return fmt.Errorf("write config %s: %w", configPath, err)
	}
	return nil
}
