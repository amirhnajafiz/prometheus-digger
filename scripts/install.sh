#!/usr/bin/env bash
set -euo pipefail

APP_NAME="promdigger"
INSTALL_DIR="/usr/local/bin"
REQUIRED_GO_MAJOR=1
REQUIRED_GO_MINOR=25

echo "🔍 Checking Go installation..."

if ! command -v go >/dev/null 2>&1; then
  echo "❌ Go is not installed."
  echo "Please install Go >= ${REQUIRED_GO_MAJOR}.${REQUIRED_GO_MINOR}"
  exit 1
fi

GO_VERSION_RAW="$(go version)"
# Example: go version go1.25.0 linux/amd64
GO_VERSION="$(echo "$GO_VERSION_RAW" | awk '{print $3}' | sed 's/^go//')"

GO_MAJOR="$(echo "$GO_VERSION" | cut -d. -f1)"
GO_MINOR="$(echo "$GO_VERSION" | cut -d. -f2)"

echo "✅ Found Go version: $GO_VERSION"

if [[ "$GO_MAJOR" -lt "$REQUIRED_GO_MAJOR" ]] || \
   [[ "$GO_MAJOR" -eq "$REQUIRED_GO_MAJOR" && "$GO_MINOR" -lt "$REQUIRED_GO_MINOR" ]]; then
  echo "❌ Go version ${GO_VERSION} is too old."
  echo "Please upgrade to Go >= ${REQUIRED_GO_MAJOR}.${REQUIRED_GO_MINOR}"
  exit 1
fi

echo "🚀 Go version is sufficient. Building..."

make

if [[ ! -f "$APP_NAME" ]]; then
  echo "❌ Build succeeded but binary '$APP_NAME' not found"
  exit 1
fi

echo "📦 Installing $APP_NAME to $INSTALL_DIR"

if [[ ! -w "$INSTALL_DIR" ]]; then
  sudo install -m 0755 "$APP_NAME" "$INSTALL_DIR/$APP_NAME"
else
  install -m 0755 "$APP_NAME" "$INSTALL_DIR/$APP_NAME"
fi

rm -rf ~/.promdigger
mkdir ~/.promdigger

cp example.config.json ~/.promdigger/config.json
chmod 644 ~/.promdigger/config.json

echo "✅ Installation complete!"
echo "Config: ~/.promdigger/config.json"
echo "👉 Run: $APP_NAME --help"
