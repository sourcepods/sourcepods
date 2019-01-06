GO := CGO_ENABLED=0 go

.PHONY: generate
generate: apiv1

.PHONY: apiv1
apiv1: internal/api/v1/models internal/api/v1/restapi

SWAGGER ?= docker run --rm \
	--user=$(shell id -u $(USER)):$(shell id -g $(USER)) \
	-v $(shell pwd):/go/src/github.com/gitpods/gitpods \
	-w /go/src/github.com/gitpods/gitpods quay.io/goswagger/swagger:v0.18.0

internal/api/v1/models internal/api/v1/restapi: swagger.yaml
	-rm -r internal/api/v1/{models,restapi}
	$(SWAGGER) generate server -f swagger.yaml --exclude-main -A gitpods --target internal/api/v1

.PHONY: test
test:
	go test -coverprofile coverage.out -race -v ./... # -race needs Cgo

.PHONY: build
build: dev/api dev/gitpods dev/storage dev/ui

dev/api: cmd/api internal
	$(GO) build -v -ldflags '-w -extldflags '-static'' -o ./dev/api ./cmd/api

dev/gitpods: cmd/gitpods internal
	$(GO) build -v -ldflags '-w -extldflags '-static'' -o ./dev/gitpods ./cmd/gitpods

dev/storage: cmd/storage internal
	$(GO) build -v -ldflags '-w -extldflags '-static'' -o ./dev/storage ./cmd/storage

dev/ui: cmd/ui internal dev/packr
	./dev/packr
	$(GO) build -v -ldflags '-w -extldflags '-static'' -o ./dev/ui ./cmd/ui

dev/packr:
	mkdir -p dev/
	curl -s -L https://github.com/gobuffalo/packr/releases/download/v1.10.4/packr_1.10.4_linux_amd64.tar.gz | tar -xz -C dev packr
