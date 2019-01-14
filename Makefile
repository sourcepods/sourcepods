GO := CGO_ENABLED=0 go

.PHONY: generate
generate: apiv1

.PHONY: apiv1
apiv1: pkg/api/v1/models pkg/api/v1/restapi ui/lib/src/api

GOSWAGGER ?= docker run --rm \
	--user=$(shell id -u $(USER)):$(shell id -g $(USER)) \
	-v $(shell pwd):/go/src/github.com/gitpods/gitpods \
	-w /go/src/github.com/gitpods/gitpods quay.io/goswagger/swagger:v0.18.0

pkg/api/v1/models pkg/api/v1/restapi: swagger.yaml
	-rm -r pkg/api/v1/{models,restapi}
	$(GOSWAGGER) generate server -f swagger.yaml --exclude-main -A gitpods --target pkg/api/v1

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
	golint $(shell go list ./pkg/gitpods/...)

.PHONY: test
test:
	go test -coverprofile coverage.out -race -v ./... # -race needs Cgo

.PHONY: build
build: dev/api dev/gitpods dev/storage dev/ui

dev/api: cmd/api pkg
	$(GO) build -v -ldflags '-w -extldflags '-static'' -o ./dev/api ./cmd/api

dev/gitpods: cmd/gitpods pkg
	$(GO) build -v -ldflags '-w -extldflags '-static'' -o ./dev/gitpods ./cmd/gitpods

dev/storage: cmd/storage pkg
	$(GO) build -v -ldflags '-w -extldflags '-static'' -o ./dev/storage ./cmd/storage

dev/ui: cmd/ui pkg dev/packr
	./dev/packr
	$(GO) build -v -ldflags '-w -extldflags '-static'' -o ./dev/ui ./cmd/ui

dev/packr:
	mkdir -p dev/
	curl -s -L https://github.com/gobuffalo/packr/releases/download/v1.10.4/packr_1.10.4_linux_amd64.tar.gz | tar -xz -C dev packr
