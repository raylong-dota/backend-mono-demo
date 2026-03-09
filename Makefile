ROOT  := $(shell pwd)
GOBIN := $(ROOT)/.go/bin
TOOLS := $(ROOT)/.tools
T     := PATH="$(GOBIN):$(TOOLS)/go/bin:$(TOOLS):$$PATH" GOMODCACHE="$(ROOT)/.go/pkg/mod"

# Positional service argument: make <target> <svc>  (e.g. make build order)
SVC := $(word 2, $(MAKECMDGOALS))

# Absorb extra positional words so make doesn't error on unknown targets
%:
	@:

.PHONY: install
# install development tools to .tools/ (run once after cloning)
install:
	@bash scripts/install_base.sh

.PHONY: new
# create a new service: make new <svc>
new:
	@bash scripts/new.sh $(SVC)

.PHONY: generate
# generate code (api+wire+proto): make generate <svc>|all
generate:
	@if [ -z "$(SVC)" ]; then \
		echo "error: service name required, e.g. make generate order  or  make generate all"; exit 1; \
	elif [ "$(SVC)" = "all" ]; then \
		find app -mindepth 2 -maxdepth 2 -type d -print | xargs -L 1 bash -c 'cd "$$0" && echo "→ $$0" && $(MAKE) api proto wire'; \
	elif [ ! -d "app/$(SVC)/service" ]; then \
		echo "error: app/$(SVC)/service not found"; exit 1; \
	else \
		cd app/$(SVC)/service && $(MAKE) api proto wire; \
	fi

.PHONY: build
# build binary into bin/orbit-<svc>-svc: make build <svc>|all
build:
	@mkdir -p $(ROOT)/bin
	@if [ -z "$(SVC)" ]; then \
		echo "error: service name required, e.g. make build order  or  make build all"; exit 1; \
	elif [ "$(SVC)" = "all" ]; then \
		find app -mindepth 2 -maxdepth 2 -type d -print | xargs -L 1 bash -c \
			'svcname=$$(basename $$(dirname "$$0")) && echo "→ orbit-$$svcname-svc" && cd "$$0" && $(MAKE) build BINARY_DIR="$(ROOT)/bin" BINARY_NAME="orbit-$$svcname-svc"'; \
	elif [ ! -d "app/$(SVC)/service" ]; then \
		echo "error: app/$(SVC)/service not found"; exit 1; \
	else \
		cd app/$(SVC)/service && $(MAKE) build BINARY_DIR="$(ROOT)/bin" BINARY_NAME="orbit-$(SVC)-svc"; \
	fi

.PHONY: run
# run a service locally: make run <svc>
run:
	@if [ -z "$(SVC)" ]; then \
		echo "error: service name required, e.g. make run order"; exit 1; \
	elif [ ! -d "app/$(SVC)/service" ]; then \
		echo "error: app/$(SVC)/service not found"; exit 1; \
	else \
		cd app/$(SVC)/service && $(MAKE) run; \
	fi

.PHONY: image
# build docker image: make image <svc> [tag=<tag>]  (default tag: latest)
image:
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "error: docker not found, please install Docker first"; exit 1; \
	elif [ -z "$(SVC)" ]; then \
		echo "error: service name required, e.g. make image order"; exit 1; \
	elif [ ! -d "app/$(SVC)/service" ]; then \
		echo "error: app/$(SVC)/service not found"; exit 1; \
	else \
		if [ ! -d vendor ]; then \
			echo "→ vendor/ not found, running go mod vendor first ..."; \
			$(T) go mod vendor; \
		fi; \
		docker build \
			-f deploy/build/Dockerfile \
			--build-arg APP_SVC=$(SVC) \
			-t orbit-$(SVC)-svc:$(if $(tag),$(tag),latest) \
			.; \
	fi

.PHONY: deploy
# deploy a service to local k8s: make deploy <svc>
deploy:
	@if ! command -v kubectl >/dev/null 2>&1; then \
		echo "error: kubectl not found, please install kubectl first"; exit 1; \
	elif [ -z "$(SVC)" ]; then \
		echo "error: service name required, e.g. make deploy helloworld"; exit 1; \
	elif [ ! -d "deploy/k8s/$(SVC)" ]; then \
		echo "error: deploy/k8s/$(SVC) not found"; exit 1; \
	else \
		echo "→ deploying $(SVC) to $$(kubectl config current-context)"; \
		kubectl apply -f deploy/k8s/$(SVC)/; \
	fi

.PHONY: lint
# run golangci-lint and auto-fix issues (local development)
lint:
	$(T) golangci-lint run --timeout 10m --fix --path-mode abs --config configs/golangci.yaml ./...

.PHONY: lint-check
# run golangci-lint in check-only mode, no auto-fix (used by CI)
lint-check:
	$(T) golangci-lint run --timeout 10m --path-mode abs --config configs/golangci.yaml ./...

.PHONY: tidy
# run go mod tidy with the project-local Go
tidy:
	$(T) go mod tidy

.PHONY: vendor
# vendor dependencies for offline / Docker builds
vendor:
	$(T) go mod vendor

.PHONY: get
# add or upgrade a dependency: make get pkg=github.com/foo/bar@v1.2.3
get:
	$(T) go get $(pkg)

.PHONY: clean
# remove build artifacts
clean:
	rm -rf $(ROOT)/bin
	find app -mindepth 2 -maxdepth 2 -type d -exec test -d '{}/bin' \; -exec rm -rf '{}/bin' \;
