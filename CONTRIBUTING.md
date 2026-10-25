# Contributing

Thanks for your interest in AI Proxy Manager.

## Development

1. Clone this repository.
2. Follow [apps/proxy-desktop/README.md](apps/proxy-desktop/README.md).
3. Run tests before opening a PR:

```bash
cd apps/proxy-desktop
make test
cd frontend && pnpm test && pnpm run build
```

## Pull requests

- Keep changes focused; match existing Go and React patterns in the app.
- Update docs when behavior or configuration changes.
- Do not commit `config.yaml`, `proxy.bootstrap.yaml`, or local database files.

## Issues

Use GitHub issues for bugs and feature requests. Include steps to reproduce and your OS version for desktop issues.
