###
### go get -u k8s.io/client-go@kubernetes-1.15.0
### go get -u k8s.io/client-go@kubernetes-1.15.0
### go get -u github.com/openshift/client-go@master
### go get -u github.com/openshift/api@master
.PHONY : update-bazel
update-bazel:
	#dep ensure
	bazel run //:gazelle
	# gazelle:resolve did not work out
	# https://github.com/bazelbuild/bazel-gazelle/issues/432#issuecomment-457789836
	sed -i -e "s|//vendor/google.golang.org/grpc/naming:go_default_library|@org_golang_google_grpc//naming:go_default_library|g" ./vendor/google.golang.org/api/internal/BUILD.bazel

validate-modules:
	go mod tidy
	go mod vendor
	git status -s ./vendor/ go.mod go.sum
	test -z "$$(git status --porcelain go.mod go.sum)"
.PHONY: validate-modules

###deprecated
download-vendor:
	go mod vendor
.PHONY: download-vendor

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

.PHONY : test-lc
test-lc:
	go test -v ./pkg/lc/...

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
	#"${GOPATH}/bin/goveralls" -coverprofile=build/coverage.out -service=circle-ci -repotoken="${COVERALLS_TOKEN}"
	echo "skipping coveralls on circleci ..."
endif

go_version := $(shell go version)

.PHONY : gen-images
gen-images:
	docker build --label "version=$$(git describe --tags --always --dirty)" --label "url=https://github.com/hongkailiu/test-go" -f test_files/docker/Dockerfile.testctl.txt -t quay.io/hongkailiu/test-go:testctl-travis .
	docker build -f test_files/docker/Dockerfile.ocptf.txt -t quay.io/hongkailiu/test-go:ocptf-travis .
#	docker build --label "version=$(git describe --tags --always --dirty)" --label "url=https://github.com/hongkailiu/test-go" --label "build_time=$(date --utc +%FT%TZ)" -f test_files/docker/Dockerfile.circleci.txt -t quay.io/hongkailiu/test-go:circleci-travis .
ifeq ($(TRAVIS)$(findstring go1.12,$(go_version)), truego1.12)
	docker tag quay.io/hongkailiu/test-go:testctl-travis "quay.io/hongkailiu/ci-staging:testctl-$(USER)-${TRAVIS_JOB_NUMBER}"
	docker tag quay.io/hongkailiu/test-go:circleci-travis "quay.io/hongkailiu/ci-staging:circleci-$(USER)-${TRAVIS_JOB_NUMBER}"
	echo "$(quay_cli_password)" | docker login -u hongkailiu quay.io --password-stdin
	docker push "quay.io/hongkailiu/ci-staging:testctl-$(USER)-${TRAVIS_JOB_NUMBER}"
#	docker push "quay.io/hongkailiu/ci-staging:circleci-$(USER)-${TRAVIS_JOB_NUMBER}"
endif
	docker images

.PHONY : build-ocptf
build-ocptf:
	go build -o ./build/ocptf ./cmd/ocptf/

.PHONY : bazel-all
bazel-all: download-vendor update-bazel
ifeq ($(CIRCLECI), true)
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

BAZELISK_VERSION := v1.0

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
	# install bazelisk
	curl -OL https://github.com/bazelbuild/bazelisk/releases/download/${BAZELISK_VERSION}/bazelisk-linux-amd64
	sudo mv ./bazelisk-linux-amd64 /usr/bin/bazel
	sudo chmod +x /usr/bin/bazel
endif

.PHONY : ci-before-script
ci-before-script:
	echo "env. var. GOPATH: $${GOPATH}"
	echo "env. var. GO111MODULE: $${GO111MODULE}"
	echo "env. var. USE_BAZEL_VERSION: $${USE_BAZEL_VERSION}"
	echo "env. var. GOPROXY: $${GOPROXY}"
	go version
	go env
	docker version
	make --version
	java -version
	bazel version

CI_SCRIPT_DEPS += validate-modules
CI_SCRIPT_DEPS += build-swagger
CI_SCRIPT_DEPS += build-code-gen
CI_SCRIPT_DEPS += test-lc
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
