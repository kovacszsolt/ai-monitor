.PHONY: dev build test frontend-install frontend-build

BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X 'github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/buildinfo.BuildTime=$(BUILD_TIME)'

dev:
	wails dev -ldflags "$(LDFLAGS)"

build: frontend-build
	wails build -ldflags "$(LDFLAGS)"

frontend-install:
	cd frontend && pnpm install

frontend-build:
	cd frontend && pnpm run build

test:
	go test ./...
