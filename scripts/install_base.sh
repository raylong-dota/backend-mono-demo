#!/usr/bin/env bash
# Install project development tools to .tools/ (project-local, git-ignored)
set -eo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
TOOLS_DIR="$ROOT_DIR/.tools"

mkdir -p "$TOOLS_DIR"

echo "→ Installing tools to .tools/"
echo ""

# ── protoc (system-level, must be installed externally) ─────────────────────
if ! command -v protoc &>/dev/null; then
    echo "✗ protoc not found. Install it first:"
    echo "  macOS:  brew install protobuf"
    echo "  Ubuntu: apt install -y protobuf-compiler"
    exit 1
fi
echo "✓ protoc $(protoc --version)"

# ── Go-based protoc plugins & tools ─────────────────────────────────────────
install_tool() {
    local pkg=$1
    local bin=$2
    echo -n "  installing $bin ... "
    GOBIN="$TOOLS_DIR" go install "$pkg"
    echo "✓"
}

install_tool "google.golang.org/protobuf/cmd/protoc-gen-go@latest"                  "protoc-gen-go"
install_tool "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"                 "protoc-gen-go-grpc"
install_tool "github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest"         "protoc-gen-go-http"
install_tool "github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest"       "protoc-gen-go-errors"
install_tool "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest" "protoc-gen-openapiv2"
install_tool "github.com/google/wire/cmd/wire@latest"                               "wire"

echo ""

# ── Download go module dependencies to project-local cache ──────────────────
echo -n "  downloading go module dependencies ... "
GOMODCACHE="$ROOT_DIR/.go/pkg/mod" go mod download
echo "✓"

echo ""
echo "✓ All done:"
echo "  .tools/       → protoc plugins & wire"
echo "  .go/pkg/mod/  → go module source (browse with IDE)"
echo ""
echo "  Run 'make api' or 'make wire' in any service directory."
