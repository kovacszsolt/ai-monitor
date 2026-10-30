package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Storage     StorageConfig     `yaml:"storage"`
	Logging     LoggingConfig     `yaml:"logging"`
	Capture     CaptureConfig     `yaml:"capture"`
	Persistence PersistenceConfig `yaml:"persistence"`
}

type ServerConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	MaxRequestBytes int64  `yaml:"max_request_bytes"`
}

type StorageConfig struct {
	DatabasePath  string `yaml:"database_path"`
	JournalMode   string `yaml:"journal_mode"`
	BusyTimeoutMS int    `yaml:"busy_timeout_ms"`
	Synchronous   string `yaml:"synchronous"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"` // console (colored) or json
}

type FileDumpConfig struct {
	Enabled      bool   `yaml:"enabled"`
	Directory    string `yaml:"directory"`
	ErrorsOnly   bool   `yaml:"errors_only"`
	MaxBodyBytes int64  `yaml:"max_body_bytes"`
}

type CaptureConfig struct {
	Mode                     string         `yaml:"mode"`
	StoreContent             bool           `yaml:"store_content"`
	StoreUpstreamBodyOnError bool           `yaml:"store_upstream_body_on_error"`
	FileDump                 FileDumpConfig `yaml:"file_dump"`
	MaxPromptBytes           int64          `yaml:"max_prompt_bytes"`
	MaxResponseBytes         int64          `yaml:"max_response_bytes"`
}

type PersistenceConfig struct {
	BatchSize     int `yaml:"batch_size"`
	BatchDelayMS  int `yaml:"batch_delay_ms"`
	QueueCapacity int `yaml:"queue_capacity"`
}

func Default() Config {
	return Config{
		Server: ServerConfig{
			Host:            "127.0.0.1",
			Port:            8080,
			MaxRequestBytes: 64 * 1024 * 1024,
		},
		Storage: StorageConfig{
			DatabasePath:  "./data/proxy.db",
			JournalMode:   "WAL",
			BusyTimeoutMS: 5000,
			Synchronous:   "NORMAL",
		},
		Logging: LoggingConfig{Level: "info", Format: "console"},
		Capture: CaptureConfig{
			Mode:                     "metadata",
			StoreContent:             false,
			StoreUpstreamBodyOnError: true,
			MaxPromptBytes:           256 * 1024,
			MaxResponseBytes:         1024 * 1024,
		},
		Persistence: PersistenceConfig{
			BatchSize:     100,
			BatchDelayMS:  250,
			QueueCapacity: 10000,
		},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return Config{}, fmt.Errorf("read config %s: %w", path, err)
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return Config{}, fmt.Errorf("parse config: %w", err)
		}
	}
	applyEnv(&cfg)
	cfg.normalize()
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c *Config) normalize() {
	if c.Capture.FileDump.Enabled && c.Capture.FileDump.MaxBodyBytes == 0 {
		c.Capture.FileDump.MaxBodyBytes = c.Capture.MaxResponseBytes
	}
}

// ApplyEnvOverrides applies AI_PROXY_* environment variables to cfg.
func ApplyEnvOverrides(cfg *Config) {
	applyEnv(cfg)
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("AI_PROXY_DATABASE"); v != "" {
		cfg.Storage.DatabasePath = v
	}
	if v := os.Getenv("AI_PROXY_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("AI_PROXY_PORT"); v != "" {
		if p, err := parseInt(v); err == nil {
			cfg.Server.Port = p
		}
	}
	if v := os.Getenv("AI_PROXY_LOG_FORMAT"); v != "" {
		cfg.Logging.Format = v
	}
	if v := os.Getenv("AI_PROXY_HTTP_DUMP_DIR"); v != "" {
		cfg.Capture.FileDump.Enabled = true
		cfg.Capture.FileDump.Directory = v
	}
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}

func (c Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be 1-65535")
	}
	if c.Storage.DatabasePath == "" {
		return fmt.Errorf("storage.database_path is required")
	}
	if c.Server.MaxRequestBytes < 1 {
		return fmt.Errorf("server.max_request_bytes must be >= 1")
	}
	if c.Capture.Mode != "metadata" {
		return fmt.Errorf("capture.mode must be metadata")
	}
	if c.Capture.MaxPromptBytes < 1 {
		return fmt.Errorf("capture.max_prompt_bytes must be >= 1")
	}
	if c.Capture.MaxResponseBytes < 1 {
		return fmt.Errorf("capture.max_response_bytes must be >= 1")
	}
	if c.Capture.MaxPromptBytes > c.Server.MaxRequestBytes {
		return fmt.Errorf("capture.max_prompt_bytes cannot exceed server.max_request_bytes")
	}
	if c.Capture.MaxResponseBytes > c.Server.MaxRequestBytes {
		return fmt.Errorf("capture.max_response_bytes cannot exceed server.max_request_bytes")
	}
	if c.Capture.FileDump.Enabled {
		if strings.TrimSpace(c.Capture.FileDump.Directory) == "" {
			return fmt.Errorf("capture.file_dump.directory is required when file_dump.enabled is true")
		}
		if c.Capture.FileDump.MaxBodyBytes < 1 {
			return fmt.Errorf("capture.file_dump.max_body_bytes must be >= 1 when file_dump is enabled")
		}
		if c.Capture.FileDump.MaxBodyBytes > c.Server.MaxRequestBytes {
			return fmt.Errorf("capture.file_dump.max_body_bytes cannot exceed server.max_request_bytes")
		}
	}
	if c.Persistence.BatchSize < 1 {
		return fmt.Errorf("persistence.batch_size must be >= 1")
	}
	if c.Persistence.QueueCapacity < 1 {
		return fmt.Errorf("persistence.queue_capacity must be >= 1")
	}
	switch c.Logging.Format {
	case "", "console", "json":
	default:
		return fmt.Errorf("logging.format must be console or json")
	}
	return nil
}

func (c Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
