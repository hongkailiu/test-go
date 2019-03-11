.PHONY : build-k8s
build-k8s:
	./script/ci/build-k8s.sh

.PHONY : build-oc
build-oc:
	./script/ci/build-oc.sh

.PHONY : update-dep
update-dep:
	./script/ci/update-dep.sh


.PHONY : build-swagger
build-swagger:
	./script/ci/build-swagger.sh


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


.PHONY : build-flexy
build-flexy:
	./script/ci/build-flexy.sh

.PHONY : create-flexy-instances
create-flexy-instances:
	./script/run/create-flexy-instances.sh

###make run-flexy-playbooks secret_vars_file=/tmp/aaa.sh ansible_repo_path=/home/hongkliu/repo/openshift/openshift-ansible install_glusterfs=false
.PHONY : run-flexy-playbooks
run-flexy-playbooks:
	./script/run/run-flexy-playbooks.sh

.PHONY : release-flexy
release-flexy:
	./script/ci/release-flexy.sh

.PHONY : test-flexy
test-flexy:
	./script/ci/test-flexy.sh

.PHONY : code-gen-clean
code-gen-clean:
	./script/ci/code-gen-clean.sh

.PHONY : code-gen
code-gen:
	./script/ci/code-gen.sh

.PHONY : build-code-gen
build-code-gen:
	./script/ci/build-code-gen.sh

.PHONY : build-http
build-http:
	./script/ci/build-http.sh

.PHONY : test-pb
test-pb:
	./script/ci/test-pb.sh

.PHONY : test-lc
test-lc:
	./script/ci/test-lc.sh

.PHONY : build-others
build-others:
	go build -o ./build/hello ./pkg/hello/
	go build -o ./build/worker_pool ./pkg/channel/

.PHONY : test-others
test-others:
	go test -v ./pkg/hello/...
	go test -v ./pkg/doc/...

.PHONY : gen-coverage
gen-coverage:
	if [ ! -d "build" ]; then mkdir -v build; fi
	go test -v -coverprofile build/coverage.out ./...

.PHONY : coveralls
coveralls:
	./script/ci/coveralls.sh

.PHONY : gen-images
gen-images:
	docker build -f test_files/docker/Dockerfile.http.txt -t quay.io/hongkailiu/test-go:http-travis .
	docker build -f test_files/docker/Dockerfile.ocptf.txt -t quay.io/hongkailiu/test-go:ocptf-travis .
	docker images

.PHONY : build-ocptf
build-ocptf:
	go build -o ./build/ocptf ./cmd/ocptf/

.PHONY : build-ocpsanity
build-ocpsanity:
	go build -o ./build/ocpsanity ./cmd/ocpsanity/

.PHONY : bazel-build
bazel-build:
	bazel build --host_jvm_args=-Xmx500m --host_jvm_args=-Xms500m //cmd/...
	### ignore those package: seems bazel needs test file and target file are in the same pkg
	### however, it is not the case for those 2 pkgs
	bazel test -- //... -//pkg/flexy/... -//pkg/ocptf/...

.PHONY : ci-install
ci-install:
	go get github.com/onsi/ginkgo/ginkgo
	cp test_files/flexy/unit.test.files/gce.json /tmp/
	@echo "install bazel ..."
	./script/ci/install-bazel.sh


.PHONY : ci-before-script
ci-before-script:
	echo "GOPATH: $${GOPATH}"
	go version
	ginkgo version
	docker version
	make --version
	java -version
	bazel version

CI_SCRIPT_DEPS := build-k8s
CI_SCRIPT_DEPS += build-oc
CI_SCRIPT_DEPS += build-swagger
CI_SCRIPT_DEPS += build-others
CI_SCRIPT_DEPS += build-flexy
CI_SCRIPT_DEPS += test-flexy
CI_SCRIPT_DEPS += code-gen
CI_SCRIPT_DEPS += build-code-gen
CI_SCRIPT_DEPS += build-http
CI_SCRIPT_DEPS += test-pb
CI_SCRIPT_DEPS += test-lc
CI_SCRIPT_DEPS += test-others
CI_SCRIPT_DEPS += gen-coverage
CI_SCRIPT_DEPS += gen-images
CI_SCRIPT_DEPS += build-ocptf
CI_SCRIPT_DEPS += build-ocpsanity
CI_SCRIPT_DEPS += bazel-build

.PHONY : ci-script
ci-script: $(CI_SCRIPT_DEPS)

.PHONY : ci-package
ci-package:
	./script/ci/package-ocptf.sh
	./script/ci/package-ocpsanity.sh
	ls -al ./build/*.tar.gz

.PHONY : ci-all
ci-all: ci-install ci-before-script ci-script ci-package coveralls
