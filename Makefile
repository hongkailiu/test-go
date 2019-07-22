.PHONY : build-k8s
build-k8s:
	go build -o build/k8s ./pkg/k8s/

.PHONY : build-oc
build-oc:
	go build -o build/oc ./pkg/oc/

.PHONY : update-dep
update-dep:
	dep ensure
	bazel run //:gazelle

.PHONY : build-swagger
build-swagger:
	go build -o build/hello-swagger ./pkg/swagger/


# The `validate-swagger` target checks for errors and inconsistencies in
# our specification of an API. This target can check if we're
# referencing inexistent definitions and gives us hints to where
# to fix problems with our API in a static manner.
validate-swagger:
	swagger validate ./pkg/swagger/swagger/swagger.yml


# The `gen-swagger` target depends on the `validate` target as
# it will only successfully generate the code if the specification
# is valid.
#
# Here we're specifying some flags:
# --target              the base directory for generating the files;
# --spec                path to the swagger specification;
# --exclude-main        generates only the library code and not a
#                       sample CLI application;
# --name                the name of the application.
gen-swagger: validate-swagger
	swagger generate server \
		--target=./pkg/swagger/swagger \
		--spec=./pkg/swagger/swagger/swagger.yml \
		--exclude-main \
		--name=hello

.PHONY : code-gen-clean
code-gen-clean:
	rm -rfv pkg/codegen/pkg/client
	rm -fv pkg/codegen/pkg/apis/app.example.com/v1alpha1/zz_generated.deepcopy.go

.PHONY : code-gen
code-gen:
	./pkg/codegen/hack/update-codegen.sh

.PHONY : build-code-gen
build-code-gen:
	go build -o build/example ./pkg/codegen/cmd/example/

.PHONY : test-pb
test-pb:
	go test -v ./pkg/probuf/unittest/...

.PHONY : test-lc
test-lc:
	go test -v ./pkg/lc/...

.PHONY : build-others
build-others:
	go build -o ./build/hello ./pkg/hello/
	go build -o ./build/worker_pool ./pkg/channel/

.PHONY : test-others
test-others:
	go test -v ./pkg/hello/...
	go test -v ./pkg/doc/...
	go test -v ./pkg/json/...
	go test -v ./pkg/http/...

.PHONY : gen-coverage
gen-coverage:
	if [ ! -d "build" ]; then mkdir -v build; fi
	go test -v -coverprofile build/coverage.out ./...

.PHONY : coveralls
coveralls:
ifneq ($(CI), true)
	echo "please run this on ci system, like travis ci or circle ci"
	false
endif
	go get -u github.com/mattn/goveralls
ifeq ($(TRAVIS), true)
	"${GOPATH}/bin/goveralls" -coverprofile=build/coverage.out -service=travis-ci
endif
ifeq ($(CIRCLECI), true)
	###https://github.com/lemurheavy/coveralls-public/issues/632
	"${GOPATH}/bin/goveralls" -coverprofile=build/coverage.out -service=circle-ci -repotoken="${COVERALLS_TOKEN}"
endif

.PHONY : gen-images
gen-images:
	docker build -f test_files/docker/Dockerfile.testctl.txt -t quay.io/hongkailiu/test-go:testctl-travis .
	docker build -f test_files/docker/Dockerfile.ocptf.txt -t quay.io/hongkailiu/test-go:ocptf-travis .
	docker build -f test_files/docker/Dockerfile.circleci.txt -t quay.io/hongkailiu/test-go:circleci-travis .
	docker images

.PHONY : build-ocptf
build-ocptf:
	go build -o ./build/ocptf ./cmd/ocptf/

.PHONY : bazel-all
bazel-all:
ifneq ($(CIRCLECI), true)
	bazel build --jobs=1 --jvmopt='-Xmx:2048m' --jvmopt='-Xms:2048m' //cmd/...
else
	bazel build //cmd/...
endif
	bazel test -- //... -//pkg/ocptf/...

build_version := $(shell git describe --tags --always --dirty)

.PHONY : build-testctl
build-testctl:
	sed -i -e "s|{buildVersion}|$(build_version)|g" ./pkg/testctl/cmd/config/version.go
	go build -o ./build/testctl ./cmd/testctl/
	cp -rv pkg/http/static build/
	cp -rv pkg/http/swagger build/
	git checkout ./pkg/testctl/cmd/config/version.go

BAZELISK_VERSION := v0.0.8

.PHONY : ci-install
ci-install:
ifneq ($(CI), true)
	echo "not supported CI environment ... failing"
	false
endif
ifeq ($(TRAVIS), true)
	echo "deb [arch=amd64] http://storage.googleapis.com/bazel-apt stable jdk1.8" | sudo tee /etc/apt/sources.list.d/bazel.list
	curl https://bazel.build/bazel-release.pub.gpg | sudo apt-key add -
	sudo apt-get update
	sudo apt-get install bazel
endif
	curl -OL https://github.com/bazelbuild/bazelisk/releases/download/${BAZELISK_VERSION}/bazelisk-linux-amd64
	sudo mv ./bazelisk-linux-amd64 /usr/bin/bazel
	sudo chmod +x /usr/bin/bazel

.PHONY : ci-before-script
ci-before-script:
	echo "GOPATH: $${GOPATH}"
	go version
	docker version
	make --version
	java -version
	bazel version

CI_SCRIPT_DEPS := build-k8s
CI_SCRIPT_DEPS += build-oc
CI_SCRIPT_DEPS += build-swagger
CI_SCRIPT_DEPS += build-others
CI_SCRIPT_DEPS += code-gen
CI_SCRIPT_DEPS += build-code-gen
CI_SCRIPT_DEPS += test-pb
CI_SCRIPT_DEPS += test-lc
CI_SCRIPT_DEPS += test-others
CI_SCRIPT_DEPS += gen-coverage
CI_SCRIPT_DEPS += build-testctl
CI_SCRIPT_DEPS += gen-images
CI_SCRIPT_DEPS += build-ocptf
CI_SCRIPT_DEPS += bazel-all

.PHONY : ci-script
ci-script: $(CI_SCRIPT_DEPS)

.PHONY : ci-package
ci-package:
	./script/ci/package-ocptf.sh
	ls -al ./build/*.tar.gz

.PHONY : ci-all
ci-all: ci-install ci-before-script ci-script ci-package coveralls
