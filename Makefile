

imports:
	which gci || go install -mod=mod github.com/daixiang0/gci@latest
	gci write --custom-order -s standard -s default -s "prefix(k8s.io)" -s "prefix(github.com/openshift)" -s localmodule --skip-vendor .
.PHONY: imports

generate-go:
	go fmt ./...
	go mod tidy
.PHONY: generate-go

unit: gotestsum
	gotestsum --packages="./..."
.PHONY: unit

test: unit
.PHONY: test

verify: test imports generate-go lint
	git diff --exit-code
.PHONY: verify

gotestsum:
	which gotestsum || go install go install gotest.tools/gotestsum@v1.13.0
.PHONY: gotestsum

lint:
	which golangci-lint || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.3
	golangci-lint run --default all --new --disable wsl
.PHONY: lint
