GOFLAGS := -mod=vendor
GO := GOFLAGS=$(GOFLAGS) GO111MODULE=on CGO_ENABLED=0 go
GOTEST := GOFLAGS=$(GOFLAGS) GO111MODULE=on CGO_ENABLED=1 go # -race needs cgo
GO_PKG_FILES := $(shell find ./pkg/ -name "*.go" -type f ! -name "*_test.go")
DARTFILES := $(shell find ./ui/lib/ -name "*.dart" -type f)

.PHONY: generate
generate: apiv1

.PHONY: apiv1
apiv1: pkg/api/v1/models pkg/api/v1/restapi ui/lib/src/api

GOSWAGGER ?= docker run --rm \
	--user=$(shell id -u $(USER)):$(shell id -g $(USER)) \
	-v $(shell pwd):/go/src/github.com/sourcepods/sourcepods \
	-w /go/src/github.com/sourcepods/sourcepods quay.io/goswagger/swagger:v0.18.0

pkg/api/v1/models pkg/api/v1/restapi: swagger.yaml
	-rm -r pkg/api/v1/{models,restapi}
	$(GOSWAGGER) generate server -f swagger.yaml --exclude-main -A sourcepods --target pkg/api/v1

SWAGGER ?= docker run --rm \
		--user=$(shell id -u $(USER)):$(shell id -g $(USER)) \
		-v $(shell pwd):/local \
		swaggerapi/swagger-codegen-cli:2.4.0

ui/lib/src/api: swagger.yaml
	-rm -rf ui/lib/src/api
	$(SWAGGER) generate -i /local/swagger.yaml -l dart -o /local/tmp/dart
	mv tmp/dart/lib ui/lib/src/api
	-rm -rf tmp/

.PHONY: lint
lint:
	golint $(shell $(GO) list ./pkg/sourcepods/...)

.PHONY: check-vendor
check-vendor:
	$(GO) mod tidy
	$(GO) mod vendor
	git update-index --refresh
	git diff-index --quiet HEAD


.PHONY: test
test:
	$(GOTEST) test -coverprofile coverage.out -race -v ./...

.PHONY: build
build: dev/api dev/sourcepods-dev dev/storage dev/ui

dev/api: cmd/api $(GO_PKG_FILES)
	$(GO) build -v -ldflags '-w -extldflags '-static'' -o ./dev/api ./cmd/api

dev/sourcepods-dev: cmd/sourcepods-dev $(GO_PKG_FILES)
	$(GO) build -v -ldflags '-w -extldflags '-static'' -o ./dev/sourcepods-dev ./cmd/sourcepods-dev

dev/storage: cmd/storage $(GO_PKG_FILES)
	$(GO) build -v -ldflags '-w -extldflags '-static'' -o ./dev/storage ./cmd/storage

dev/ui: cmd/ui $(GO_PKG_FILES) ui/build dev/packr
	./dev/packr
	$(GO) build -v -ldflags '-w -extldflags '-static'' -o ./dev/ui ./cmd/ui

ui/build: $(DARTFILES)
	cd ui && webdev build

dev/packr:
	mkdir -p dev/
	curl -s -L https://github.com/gobuffalo/packr/releases/download/v1.10.4/packr_1.10.4_linux_amd64.tar.gz | tar -xz -C dev packr
