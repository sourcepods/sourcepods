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
