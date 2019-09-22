IMAGE_NAME = quay.io/hongkailiu/test-go:s2i-1.0.1

.PHONY: build
build:
	buildah bud --label "version=$(git describe --tags --always --dirty)" --label "url=https://github.com/hongkailiu/test-go"  --format=docker -f ./Dockerfile -t $(IMAGE_NAME) .

.PHONY: test
test:
	docker build -t $(IMAGE_NAME)-candidate .
	IMAGE_NAME=$(IMAGE_NAME)-candidate test/run
