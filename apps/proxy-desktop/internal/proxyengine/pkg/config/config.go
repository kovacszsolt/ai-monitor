package config

import icfg "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"

type Config = icfg.Config
type StorageConfig = icfg.StorageConfig
type ServerConfig = icfg.ServerConfig
type LoggingConfig = icfg.LoggingConfig
type CaptureConfig = icfg.CaptureConfig
type PersistenceConfig = icfg.PersistenceConfig

func Default() Config { return icfg.Default() }
func Load(path string) (Config, error) { return icfg.Load(path) }
func ApplyEnvOverrides(cfg *Config) { icfg.ApplyEnvOverrides(cfg) }
