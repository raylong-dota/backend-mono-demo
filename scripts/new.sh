#!/usr/bin/env bash
# Usage: make new <service-name>
# Example: make new order
set -eo pipefail

SVC="${1}"

# ── Validate ────────────────────────────────────────────────────────────────
if [ -z "$SVC" ]; then
    echo "Error: service name required"
    echo "Usage: make new <service-name>"
    echo "Example: make new order"
    exit 1
fi

if [ "$SVC" = "all" ]; then
    echo "Error: 'all' is a reserved keyword and cannot be used as a service name"
    exit 1
fi

if ! echo "$SVC" | grep -qE '^[a-z][a-z0-9]*$'; then
    echo "Error: service name must be lowercase letters and numbers only (e.g. order, quote, position)"
    exit 1
fi

# ── Name variants ───────────────────────────────────────────────────────────
# order → Order  (Go type prefix)
SVC_TITLE="$(echo "${SVC:0:1}" | tr '[:lower:]' '[:upper:]')${SVC:1}"
# order → ORDER  (enum values)
SVC_UPPER="$(echo "$SVC" | tr '[:lower:]' '[:upper:]')"

# ── Paths ───────────────────────────────────────────────────────────────────
ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
TPL_APP="$ROOT_DIR/app/helloworld/service"
TPL_API="$ROOT_DIR/api/helloworld/service/v1"
DST_APP="$ROOT_DIR/app/$SVC/service"
DST_API="$ROOT_DIR/api/$SVC/service/v1"

if [ -d "$DST_APP" ]; then
    echo "Error: app/$SVC/service already exists"
    exit 1
fi

# ── Portable sed -i ─────────────────────────────────────────────────────────
sedi() {
    if [ "$(uname)" = "Darwin" ]; then
        sed -i '' "$@"
    else
        sed -i "$@"
    fi
}

echo "→ Creating service: $SVC"

# ── 1. Directory tree ───────────────────────────────────────────────────────
mkdir -p \
    "$DST_API" \
    "$DST_APP/cmd/server" \
    "$DST_APP/configs" \
    "$DST_APP/internal/biz" \
    "$DST_APP/internal/conf" \
    "$DST_APP/internal/data" \
    "$DST_APP/internal/server" \
    "$DST_APP/internal/service"

# ── 2. Copy source files (skip generated: *.pb.go, wire_gen.go, *.swagger.json) ──
# API: proto definitions only
cp "$TPL_API/greeter.proto"       "$DST_API/${SVC}.proto"
cp "$TPL_API/error_reason.proto"  "$DST_API/error_reason.proto"

# App entry point (not wire_gen.go)
cp "$TPL_APP/cmd/server/main.go"          "$DST_APP/cmd/server/main.go"
cp "$TPL_APP/cmd/server/wire.go"          "$DST_APP/cmd/server/wire.go"

# Business logic layer
cp "$TPL_APP/internal/biz/biz.go"         "$DST_APP/internal/biz/biz.go"
cp "$TPL_APP/internal/biz/greeter.go"     "$DST_APP/internal/biz/${SVC}.go"

# Service (transport handler) layer
cp "$TPL_APP/internal/service/service.go" "$DST_APP/internal/service/service.go"
cp "$TPL_APP/internal/service/greeter.go" "$DST_APP/internal/service/${SVC}.go"

# Data access layer
cp "$TPL_APP/internal/data/data.go"       "$DST_APP/internal/data/data.go"
cp "$TPL_APP/internal/data/greeter.go"    "$DST_APP/internal/data/${SVC}.go"

# Server setup
cp "$TPL_APP/internal/server/server.go"   "$DST_APP/internal/server/server.go"
cp "$TPL_APP/internal/server/grpc.go"     "$DST_APP/internal/server/grpc.go"
cp "$TPL_APP/internal/server/http.go"     "$DST_APP/internal/server/http.go"

# Config proto (not conf.pb.go)
cp "$TPL_APP/internal/conf/conf.proto"    "$DST_APP/internal/conf/conf.proto"

# Runtime config and generate directive
cp "$TPL_APP/configs/config.yaml"         "$DST_APP/configs/config.yaml"
[ -f "$TPL_APP/generate.go" ] && cp "$TPL_APP/generate.go" "$DST_APP/generate.go"

# Makefile from root template
cp "$ROOT_DIR/app_makefile" "$DST_APP/Makefile"

# ── 3. Apply substitutions to all copied files ──────────────────────────────
find "$DST_API" "$DST_APP" -type f | while read -r f; do
    # Import paths first (most specific patterns, avoids partial double-replace)
    sedi "s|app/helloworld/service|app/${SVC}/service|g"       "$f"
    sedi "s|api/helloworld/service/v1|api/${SVC}/service/v1|g" "$f"

    # Remaining helloworld occurrences (proto package declarations, go_package options)
    sedi "s/helloworld/${SVC}/g" "$f"

    # Entity type names (Title case) — order matters: longer compound names first
    sedi "s/GreeterService/${SVC_TITLE}Service/g"   "$f"
    sedi "s/GreeterUsecase/${SVC_TITLE}Usecase/g"   "$f"
    sedi "s/GreeterRepo/${SVC_TITLE}Repo/g"         "$f"
    sedi "s/Greeter/${SVC_TITLE}/g"                 "$f"

    # Entity variable/function names (lowercase)
    sedi "s/greeterService/${SVC}Service/g"         "$f"
    sedi "s/greeterUsecase/${SVC}Usecase/g"         "$f"
    sedi "s/greeterRepo/${SVC}Repo/g"               "$f"
    sedi "s/greeter/${SVC}/g"                       "$f"

    # Enum values (UPPER case)
    sedi "s/GREETER/${SVC_UPPER}/g"                 "$f"

    # Proto message names
    sedi "s/HelloRequest/${SVC_TITLE}Request/g"     "$f"
    sedi "s/HelloReply/${SVC_TITLE}Reply/g"         "$f"
done

echo ""
echo "✓ Service created:"
echo "  api/$SVC/service/v1/"
echo "  app/$SVC/service/"
echo ""
echo "Next steps:"
echo "  make generate svc=$SVC   # generate protobuf + wire code"
echo "  make build    svc=$SVC   # build bin/orbit-$SVC-svc"
echo "  make run      svc=$SVC   # run the service"
