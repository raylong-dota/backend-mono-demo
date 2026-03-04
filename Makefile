ROOT  := $(shell pwd)
GOBIN := $(ROOT)/.go/bin
TOOLS := $(ROOT)/.tools
T     := PATH="$(GOBIN):$(TOOLS)/go/bin:$(TOOLS):$$PATH" GOMODCACHE="$(ROOT)/.go/pkg/mod"

.PHONY: install
# install development tools to .tools/ (run once after cloning)
install:
	@bash scripts/install_base.sh

.PHONY: new
# create a new service: make new svc=<name>  (e.g. make new svc=order)
new:
	@bash scripts/new.sh $(svc)

.PHONY: image
# build docker image: make image svc=<name> [tag=<tag>]  (default tag: latest)
image:
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "error: docker not found, please install Docker first"; exit 1; \
	elif [ -z "$(svc)" ]; then \
		echo "error: svc is required, e.g. make image svc=order"; exit 1; \
	elif [ ! -d "app/$(svc)/service" ]; then \
		echo "error: app/$(svc)/service not found"; exit 1; \
	else \
		docker build \
			--build-arg APP_SVC=$(svc) \
			-t orbit-$(svc)-svc:$(if $(tag),$(tag),latest) \
			.; \
	fi

.PHONY: clean
# remove build artifacts
clean:
	rm -rf $(ROOT)/bin
	find app -mindepth 2 -maxdepth 2 -type d -exec test -d '{}/bin' \; -exec rm -rf '{}/bin' \;

.PHONY: tidy
# run go mod tidy with the project-local Go
tidy:
	$(T) go mod tidy

.PHONY: get
# add or upgrade a dependency: make get pkg=github.com/foo/bar@v1.2.3
get:
	$(T) go get $(pkg)

.PHONY: lint
# run golangci-lint and auto-fix issues (local development)
lint:
	$(T) golangci-lint run --timeout 10m --fix --path-mode abs --config configs/golangci.yaml ./...

.PHONY: lint-check
# run golangci-lint in check-only mode, no auto-fix (used by CI)
lint-check:
	$(T) golangci-lint run --timeout 10m --path-mode abs --config configs/golangci.yaml ./...

.PHONY: generate
# generate code (api+wire+proto): make generate svc=<name>|all
generate:
	@if [ -z "$(svc)" ]; then \
		echo "error: svc is required, e.g. make generate svc=order  or  make generate svc=all"; exit 1; \
	elif [ "$(svc)" = "all" ]; then \
		find app -mindepth 2 -maxdepth 2 -type d -print | xargs -L 1 bash -c 'cd "$$0" && echo "→ $$0" && $(MAKE) api proto wire'; \
	elif [ ! -d "app/$(svc)/service" ]; then \
		echo "error: app/$(svc)/service not found"; exit 1; \
	else \
		cd app/$(svc)/service && $(MAKE) api proto wire; \
	fi

.PHONY: run
# run a service locally: make run svc=<name>
run:
	@if [ -z "$(svc)" ]; then \
		echo "error: svc is required, e.g. make run svc=order"; exit 1; \
	elif [ ! -d "app/$(svc)/service" ]; then \
		echo "error: app/$(svc)/service not found"; exit 1; \
	else \
		cd app/$(svc)/service && $(MAKE) run; \
	fi

.PHONY: build
# build binary into bin/<svc>: make build svc=<name>|all
build:
	@mkdir -p $(ROOT)/bin
	@if [ -z "$(svc)" ]; then \
		echo "error: svc is required, e.g. make build svc=order  or  make build svc=all"; exit 1; \
	elif [ "$(svc)" = "all" ]; then \
		find app -mindepth 2 -maxdepth 2 -type d -print | xargs -L 1 bash -c \
			'svcname=$$(basename $$(dirname "$$0")) && echo "→ orbit-$$svcname-svc" && cd "$$0" && $(MAKE) build BINARY_DIR="$(ROOT)/bin" BINARY_NAME="orbit-$$svcname-svc"'; \
	elif [ ! -d "app/$(svc)/service" ]; then \
		echo "error: app/$(svc)/service not found"; exit 1; \
	else \
		cd app/$(svc)/service && $(MAKE) build BINARY_DIR="$(ROOT)/bin" BINARY_NAME="orbit-$(svc)-svc"; \
	fi