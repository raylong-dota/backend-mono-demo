#!/usr/bin/env bash
# Install project development tools (project-local, git-ignored)
#   .go/bin/     — Go tool binaries (GOBIN)
#   .go/pkg/mod/ — Go module cache (GOMODCACHE)
#   .tools/      — non-Go tool binaries (protoc, etc.)
set -eo pipefail

# ── Versions (update here when upgrading) ────────────────────────────────────
PROTOC_VERSION="33.4"

# ── Paths ────────────────────────────────────────────────────────────────────
ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
GOBIN_DIR="$ROOT_DIR/.go/bin"
TOOLS_DIR="$ROOT_DIR/.tools"

mkdir -p "$GOBIN_DIR" "$TOOLS_DIR"

echo "→ Installing tools"
echo ""

# ── protoc → .tools/ ─────────────────────────────────────────────────────────
install_protoc() {
    if [ -f "$TOOLS_DIR/protoc" ]; then
        echo "✓ $("$TOOLS_DIR/protoc" --version) (already installed)"
        return
    fi

    local os arch zip_name
    os=$(uname -s)
    arch=$(uname -m)

    case "$os" in
        Darwin)
            case "$arch" in
                arm64)  zip_name="protoc-${PROTOC_VERSION}-osx-aarch_64.zip" ;;
                x86_64) zip_name="protoc-${PROTOC_VERSION}-osx-x86_64.zip" ;;
                *) echo "✗ unsupported macOS arch: $arch"; exit 1 ;;
            esac ;;
        Linux)
            case "$arch" in
                aarch64|arm64) zip_name="protoc-${PROTOC_VERSION}-linux-aarch_64.zip" ;;
                x86_64)        zip_name="protoc-${PROTOC_VERSION}-linux-x86_64.zip" ;;
                *) echo "✗ unsupported Linux arch: $arch"; exit 1 ;;
            esac ;;
        *) echo "✗ unsupported OS: $os"; exit 1 ;;
    esac

    local url="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${zip_name}"
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap 'rm -rf "$tmp_dir"' EXIT

    echo -n "  installing protoc v${PROTOC_VERSION} ... "
    curl -fsSL "$url" -o "$tmp_dir/protoc.zip"
    unzip -q "$tmp_dir/protoc.zip" bin/protoc -d "$tmp_dir"
    mv "$tmp_dir/bin/protoc" "$TOOLS_DIR/protoc"
    chmod +x "$TOOLS_DIR/protoc"
    trap - EXIT
    rm -rf "$tmp_dir"
    echo "✓"
}

install_protoc

# ── Go tool binaries → .go/bin/ ──────────────────────────────────────────────
install_go_tool() {
    local pkg=$1
    local bin=$2

    if [ -f "$GOBIN_DIR/$bin" ]; then
        echo "✓ $bin (already installed)"
        return
    fi

    echo -n "  installing $bin ... "
    GOBIN="$GOBIN_DIR" GOMODCACHE="$ROOT_DIR/.go/pkg/mod" go install "$pkg"
    echo "✓"
}

install_go_tool "google.golang.org/protobuf/cmd/protoc-gen-go@latest"                   "protoc-gen-go"
install_go_tool "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"                  "protoc-gen-go-grpc"
install_go_tool "github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest"          "protoc-gen-go-http"
install_go_tool "github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest"        "protoc-gen-go-errors"
install_go_tool "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest" "protoc-gen-openapiv2"
install_go_tool "github.com/google/wire/cmd/wire@latest"                                "wire"
install_go_tool "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"         "golangci-lint"

echo ""

# ── Go module dependencies → .go/pkg/mod/ ────────────────────────────────────
echo -n "  downloading go module dependencies ... "
GOMODCACHE="$ROOT_DIR/.go/pkg/mod" go mod download
echo "✓"

echo ""
echo "✓ All done:"
echo "  .tools/       → protoc v${PROTOC_VERSION}"
echo "  .go/bin/      → Go tool binaries"
echo "  .go/pkg/mod/  → Go module source (browse with IDE)"
echo ""
echo "  Run 'make api', 'make wire', or 'make lint' in any service directory."
