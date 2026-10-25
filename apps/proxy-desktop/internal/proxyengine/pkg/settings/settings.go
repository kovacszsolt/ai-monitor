package settings

import (
	"context"
	"database/sql"

	proxyconfig "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/config"
	isettings "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/settings"
)

func SeedFromConfig(ctx context.Context, db *sql.DB, fileCfg proxyconfig.Config) error {
	return isettings.SeedFromConfig(ctx, db, fileCfg)
}

func LoadProxyConfig(ctx context.Context, db *sql.DB, base proxyconfig.Config) (proxyconfig.Config, error) {
	return isettings.LoadProxyConfig(ctx, db, base)
}

type ProxyVariables = isettings.ProxyVariables

func ReadProxyVariables(ctx context.Context, db *sql.DB, base proxyconfig.Config) (ProxyVariables, error) {
	return isettings.ReadProxyVariables(ctx, db, base)
}

func SaveProxyVariables(ctx context.Context, db *sql.DB, base proxyconfig.Config, in ProxyVariables) (proxyconfig.Config, error) {
	return isettings.SaveProxyVariables(ctx, db, base, in)
}
