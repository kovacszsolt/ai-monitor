package settings

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"
)

const (
	KeyServerHost            = "server.host"
	KeyServerPort            = "server.port"
	KeyServerMaxRequestBytes = "server.max_request_bytes"
	KeyLoggingLevel          = "logging.level"
	KeyLoggingFormat         = "logging.format"
	KeyCaptureMode           = "capture.mode"
	KeyCaptureStoreContent   = "capture.store_content"
	KeyCaptureMaxPromptBytes = "capture.max_prompt_bytes"
	KeyCaptureMaxResponseBytes = "capture.max_response_bytes"
)

// ProxyVariables is the editable proxy runtime configuration stored in settings.
type ProxyVariables struct {
	ServerHost          string
	ServerPort          int
	ServerMaxRequestBytes int64
	LoggingLevel        string
	LoggingFormat       string
	CaptureMode         string
	CaptureStoreContent bool
	CaptureMaxPromptBytes int64
	CaptureMaxResponseBytes int64
}

func getOptional(ctx context.Context, db *sql.DB, key string) (string, bool, error) {
	v, err := Get(ctx, db, key)
	if err == nil {
		return v, true, nil
	}
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, ErrEmptyValue) {
		return "", false, nil
	}
	return "", false, err
}

func setIfMissing(ctx context.Context, db *sql.DB, key, value string) error {
	_, ok, err := getOptional(ctx, db, key)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return Set(ctx, db, key, value)
}

// SeedFromConfig writes file/bootstrap values into settings for keys that are not yet set.
func SeedFromConfig(ctx context.Context, db *sql.DB, fileCfg config.Config) error {
	if err := migrateLegacyKeyPepper(ctx, db); err != nil {
		return err
	}
	entries := []struct{ key, value string }{
		{KeyServerHost, fileCfg.Server.Host},
		{KeyServerMaxRequestBytes, strconv.FormatInt(fileCfg.Server.MaxRequestBytes, 10)},
		{KeyLoggingLevel, fileCfg.Logging.Level},
		{KeyLoggingFormat, fileCfg.Logging.Format},
		{KeyCaptureMode, fileCfg.Capture.Mode},
		{KeyCaptureStoreContent, strconv.FormatBool(fileCfg.Capture.StoreContent)},
		{KeyCaptureMaxPromptBytes, strconv.FormatInt(fileCfg.Capture.MaxPromptBytes, 10)},
		{KeyCaptureMaxResponseBytes, strconv.FormatInt(fileCfg.Capture.MaxResponseBytes, 10)},
	}
	if fileCfg.Server.Port >= 1 && fileCfg.Server.Port <= 65535 {
		entries = append(entries, struct{ key, value string }{KeyServerPort, strconv.Itoa(fileCfg.Server.Port)})
	}
	for _, e := range entries {
		if strings.TrimSpace(e.value) == "" {
			continue
		}
		if err := setIfMissing(ctx, db, e.key, e.value); err != nil {
			return err
		}
	}
	return nil
}

// LoadProxyConfig overlays settings from db onto base (typically YAML or defaults).
func LoadProxyConfig(ctx context.Context, db *sql.DB, base config.Config) (config.Config, error) {
	if err := migrateLegacyKeyPepper(ctx, db); err != nil {
		return config.Config{}, err
	}
	cfg := base
	if v, ok, err := getOptional(ctx, db, KeyServerHost); err != nil {
		return config.Config{}, err
	} else if ok {
		cfg.Server.Host = v
	}
	if v, ok, err := getOptional(ctx, db, KeyServerPort); err != nil {
		return config.Config{}, err
	} else if ok {
		n, err := strconv.Atoi(v)
		if err != nil {
			return config.Config{}, fmt.Errorf("settings %s: %w", KeyServerPort, err)
		}
		if n >= 1 && n <= 65535 {
			cfg.Server.Port = n
		}
	}
	if v, ok, err := getOptional(ctx, db, KeyServerMaxRequestBytes); err != nil {
		return config.Config{}, err
	} else if ok {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return config.Config{}, fmt.Errorf("settings %s: %w", KeyServerMaxRequestBytes, err)
		}
		cfg.Server.MaxRequestBytes = n
	}
	if v, ok, err := getOptional(ctx, db, KeyLoggingLevel); err != nil {
		return config.Config{}, err
	} else if ok {
		cfg.Logging.Level = v
	}
	if v, ok, err := getOptional(ctx, db, KeyLoggingFormat); err != nil {
		return config.Config{}, err
	} else if ok {
		cfg.Logging.Format = v
	}
	if v, ok, err := getOptional(ctx, db, KeyCaptureMode); err != nil {
		return config.Config{}, err
	} else if ok {
		cfg.Capture.Mode = v
	}
	if v, ok, err := getOptional(ctx, db, KeyCaptureStoreContent); err != nil {
		return config.Config{}, err
	} else if ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return config.Config{}, fmt.Errorf("settings %s: %w", KeyCaptureStoreContent, err)
		}
		cfg.Capture.StoreContent = b
	}
	if v, ok, err := getOptional(ctx, db, KeyCaptureMaxPromptBytes); err != nil {
		return config.Config{}, err
	} else if ok {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return config.Config{}, fmt.Errorf("settings %s: %w", KeyCaptureMaxPromptBytes, err)
		}
		cfg.Capture.MaxPromptBytes = n
	}
	if v, ok, err := getOptional(ctx, db, KeyCaptureMaxResponseBytes); err != nil {
		return config.Config{}, err
	} else if ok {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return config.Config{}, fmt.Errorf("settings %s: %w", KeyCaptureMaxResponseBytes, err)
		}
		cfg.Capture.MaxResponseBytes = n
	}
	return cfg, nil
}

// ReadProxyVariables returns current runtime values (DB overlay on base).
func ReadProxyVariables(ctx context.Context, db *sql.DB, base config.Config) (ProxyVariables, error) {
	cfg, err := LoadProxyConfig(ctx, db, base)
	if err != nil {
		return ProxyVariables{}, err
	}
	return ProxyVariables{
		ServerHost:              cfg.Server.Host,
		ServerPort:              cfg.Server.Port,
		ServerMaxRequestBytes:   cfg.Server.MaxRequestBytes,
		LoggingLevel:            cfg.Logging.Level,
		LoggingFormat:           cfg.Logging.Format,
		CaptureMode:             cfg.Capture.Mode,
		CaptureStoreContent:     cfg.Capture.StoreContent,
		CaptureMaxPromptBytes:   cfg.Capture.MaxPromptBytes,
		CaptureMaxResponseBytes: cfg.Capture.MaxResponseBytes,
	}, nil
}

// SaveProxyVariables validates and persists variables to settings.
func SaveProxyVariables(ctx context.Context, db *sql.DB, base config.Config, in ProxyVariables) (config.Config, error) {
	cfg, err := LoadProxyConfig(ctx, db, base)
	if err != nil {
		return config.Config{}, err
	}
	cfg.Server.Host = strings.TrimSpace(in.ServerHost)
	cfg.Server.Port = in.ServerPort
	cfg.Server.MaxRequestBytes = in.ServerMaxRequestBytes
	cfg.Logging.Level = strings.TrimSpace(in.LoggingLevel)
	cfg.Logging.Format = strings.TrimSpace(in.LoggingFormat)
	cfg.Capture.Mode = strings.TrimSpace(in.CaptureMode)
	if cfg.Capture.Mode == "" {
		cfg.Capture.Mode = "metadata"
	}
	cfg.Capture.StoreContent = in.CaptureStoreContent
	cfg.Capture.MaxPromptBytes = in.CaptureMaxPromptBytes
	cfg.Capture.MaxResponseBytes = in.CaptureMaxResponseBytes
	if err := cfg.Validate(); err != nil {
		return config.Config{}, err
	}
	writes := []struct{ key, value string }{
		{KeyServerHost, cfg.Server.Host},
		{KeyServerPort, strconv.Itoa(cfg.Server.Port)},
		{KeyServerMaxRequestBytes, strconv.FormatInt(cfg.Server.MaxRequestBytes, 10)},
		{KeyLoggingLevel, cfg.Logging.Level},
		{KeyLoggingFormat, cfg.Logging.Format},
		{KeyCaptureMode, cfg.Capture.Mode},
		{KeyCaptureStoreContent, strconv.FormatBool(cfg.Capture.StoreContent)},
		{KeyCaptureMaxPromptBytes, strconv.FormatInt(cfg.Capture.MaxPromptBytes, 10)},
		{KeyCaptureMaxResponseBytes, strconv.FormatInt(cfg.Capture.MaxResponseBytes, 10)},
	}
	for _, w := range writes {
		if err := Set(ctx, db, w.key, w.value); err != nil {
			return config.Config{}, err
		}
	}
	return cfg, nil
}
