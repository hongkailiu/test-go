

imports:
	which gci || go install -mod=mod github.com/daixiang0/gci@latest
	gci write --custom-order -s standard -s default -s "prefix(k8s.io)" -s "prefix(github.com/openshift)" -s localmodule --skip-vendor .
.PHONY: imports

verify-go:
	go fmt ./...
	go mod tidy
	git diff --exit-code
.PHONY: verify-go

verify: imports verify-go
.PHONY: verify

run:
	go run ./cmd/cincinnati
.PHONY: run
