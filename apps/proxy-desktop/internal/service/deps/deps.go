package deps

import (
	"context"
	"database/sql"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/config"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/auth"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/projects"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/providers"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/routing"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyruntime"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/repository"
)

// Desktop holds shared dependencies for Wails-facing domain services.
type Desktop struct {
	Cfg     config.Config
	Bridge  config.Bridge
	DB      *sql.DB
	Runtime *proxyruntime.Runtime

	Projects *projects.Store
	Providers *providers.Store
	Routes   *routing.Store
	Keys     *auth.KeyStore
	Requests *repository.AIRequests
	ChatPins *repository.ChatPinnedPrompts

	Ctx func() context.Context
}

func NewDesktop(cfg config.Config, bridge config.Bridge, db *sql.DB, runtime *proxyruntime.Runtime, ctxFn func() context.Context) *Desktop {
	d := &Desktop{
		Cfg:     cfg,
		Bridge:  bridge,
		DB:      db,
		Runtime: runtime,
		Ctx:     ctxFn,
	}
	d.RebindDB(db)
	return d
}

func (d *Desktop) RebindDB(db *sql.DB) {
	d.DB = db
	d.Projects = projects.NewStore(db)
	d.Providers = providers.NewStore(db)
	d.Routes = routing.NewStore(db)
	d.Keys = auth.NewKeyStore(db)
	d.Requests = repository.NewAIRequests(db)
	d.ChatPins = repository.NewChatPinnedPrompts(db)
}
