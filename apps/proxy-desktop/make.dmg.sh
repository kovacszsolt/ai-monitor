#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="$ROOT/build/bin"
APP=ai-proxy-manager.app
VOL="${DMG_VOLUME_NAME:-AI Proxy Manager}"
DMG_BASENAME="${DMG_BASENAME:-ai-proxy-manager}"
if [[ -n "${VERSION:-}" && "${VERSION}" != "dev" ]]; then
  DMG="${DMG_BASENAME}-${VERSION}-mac.dmg"
else
  DMG="${DMG_BASENAME}-mac.dmg"
fi

if [[ ! -d "$BIN_DIR/$APP" ]]; then
  echo "error: $BIN_DIR/$APP not found — run make build first" >&2
  exit 1
fi

STAGE="$(mktemp -d)"
trap 'rm -rf "$STAGE"' EXIT

cp -R "$BIN_DIR/$APP" "$STAGE/"
ln -s /Applications "$STAGE/Applications"

hdiutil create -volname "$VOL" -srcfolder "$STAGE" -ov -format UDZO "$BIN_DIR/$DMG"
echo "Created $BIN_DIR/$DMG"
