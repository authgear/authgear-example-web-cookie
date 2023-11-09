export include .env

FLAGS ?=
TAG ?= git-$(shell git rev-parse --short=8 HEAD)
IMAGE = quay.io/theauthgear/authgear-example-web-cookie:$(TAG)

.PHONY: vendor
vendor:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.50.1
	go mod download

.PHONY: test
test:
	go test ./pkg/... -timeout 1m30s $(FLAGS)

.PHONY: lint
lint:
	golangci-lint run ./cmd/... ./pkg/... --timeout 7m

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: check-tidy
check-tidy:
	$(MAKE) fmt
	go mod tidy
	git status --porcelain | grep '.*'; test $$? -eq 1

.PHONY: start
start:
	go run ./cmd/serve/main.go

.PHONY: image-build
image-build:
	docker buildx build --platform linux/amd64 . -t '$(IMAGE)' $(ARGS)
