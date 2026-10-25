package app

import iapp "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/app"
import "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/config"

type App = iapp.App

func Open(cfg config.Config, migrationsDir string) (*App, error) {
	return iapp.Open(cfg, migrationsDir)
}

func DefaultMigrationsDir() string { return iapp.DefaultMigrationsDir() }

func SetupLogger(cfg config.LoggingConfig) { iapp.SetupLogger(cfg) }
