GOBIN := $(shell pwd)/.go/bin
TOOLS := $(shell pwd)/.tools
T     := PATH="$(GOBIN):$(TOOLS):$$PATH" GOMODCACHE="$(shell pwd)/.go/pkg/mod"

.PHONY: install
# install development tools to .tools/ (run once after cloning)
install:
	@bash scripts/install_base.sh

.PHONY: new
# create a new service: make new svcn=<name>  (e.g. make new svcn=order)
new:
	@bash scripts/new.sh $(svcn)

.PHONY: lint
# run golangci-lint and auto-fix issues (local development)
lint:
	$(T) golangci-lint run --timeout 10m --fix --path-mode abs --config configs/golangci.yaml ./...

.PHONY: lint-check
# run golangci-lint in check-only mode, no auto-fix (used by CI)
lint-check:
	$(T) golangci-lint run --timeout 10m --path-mode abs --config configs/golangci.yaml ./...

.PHONY: api
# generate api
api:
	find app -mindepth 2 -maxdepth 2 -type d -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) api'

.PHONY: wire
# generate wire
wire:
	find app -mindepth 2 -maxdepth 2 -type d -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) wire'

.PHONY: proto
# generate proto
proto:
	find app -mindepth 2 -maxdepth 2 -type d -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) proto'