.PHONY : build-k8s
build-k8s:
	./script/ci/build-k8s.sh

.PHONY : build-oc
build-oc:
	./script/ci/build-oc.sh

.PHONY : update-dep
update-dep:
	dep ensure
	bazel run //:gazelle

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
	docker build -f test_files/docker/Dockerfile.testctl.txt -t quay.io/hongkailiu/test-go:testctl-travis .
	docker build -f test_files/docker/Dockerfile.ocptf.txt -t quay.io/hongkailiu/test-go:ocptf-travis .
	docker images

.PHONY : build-ocptf
build-ocptf:
	go build -o ./build/ocptf ./cmd/ocptf/

.PHONY : bazel-build
bazel-build:
	./script/ci/bazel-all.sh

build_version := $(shell git describe --tags --always --dirty)

.PHONY : build-testctl
build-testctl:
	sed -i -e "s|{buildVersion}|$(build_version)|g" ./pkg/testctl/cmd/config/version.go
	go build -o ./build/testctl ./cmd/testctl/
	cp -rv pkg/http/static build/
	cp -rv pkg/http/swagger build/
	git checkout ./pkg/testctl/cmd/config/version.go

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
CI_SCRIPT_DEPS += test-pb
CI_SCRIPT_DEPS += test-lc
CI_SCRIPT_DEPS += test-others
CI_SCRIPT_DEPS += gen-coverage
CI_SCRIPT_DEPS += build-testctl
CI_SCRIPT_DEPS += gen-images
CI_SCRIPT_DEPS += build-ocptf
CI_SCRIPT_DEPS += bazel-build

.PHONY : ci-script
ci-script: $(CI_SCRIPT_DEPS)

.PHONY : ci-package
ci-package:
	./script/ci/package-ocptf.sh
	ls -al ./build/*.tar.gz

.PHONY : ci-all
ci-all: ci-install ci-before-script ci-script ci-package coveralls

current_oc_context := $(shell oc config current-context)
oc_project := $(shell echo $(current_oc_context) | cut -d "/" -f1)
oc_server := $(shell echo $(current_oc_context) | cut -d "/" -f2)
oc_user := $(shell echo $(current_oc_context) | cut -d "/" -f3)

expected_oc_server := api-hongkliu1-qe-devcluster-openshift-com:6443
expected_oc_user := kube:admin
web_secret_file := /home/hongkliu/repo/me/svt-secret/test_go/web_secret.yaml
grafana_secret_file := /home/hongkliu/repo/me/svt-secret/test_go/grafana_secret.yaml
slack_api_secret_value := $(shell head -n 1 /home/hongkliu/repo/me/svt-secret/test_go/slack_api_secret.txt)

.PHONY : oc-deploy-testctl
oc-deploy-testctl:
	@echo "deploy testctl on openshift starter ... with $(current_oc_context)"
	@echo "oc_project: $(oc_project)"
	@echo "oc_server: $(oc_server)"
	@echo "oc_user: $(oc_user)"
ifeq ($(oc_server),$(expected_oc_server))
	@echo "server match!"
else
	@echo "server do NOT match: exiting ..."
	@echo "expected_oc_server: $(expected_oc_server)"
	false
endif
ifeq ($(oc_user),$(expected_oc_user))
	@echo "user match!"
else
	@echo "user do NOT match: exiting ..."
	@echo "expected_oc_user: $(expected_oc_user)"
	false
endif
	@echo "deploy component http web server ..."
	oc apply -f $(web_secret_file)
	oc apply -f ./deploy/testctl_http/web_deploy.yaml
	oc create configmap -n hongkliu-stage prometheus --from-file=./deploy/testctl_http/prometheus.yml --from-file=alert.rules.yml=./deploy/testctl_http/prometheus_alert.rules.yml --dry-run -o yaml | oc apply -f -
	oc apply -f ./deploy/testctl_http/prometheus_deploy.yaml
	oc apply -f ./deploy/testctl_http/status_deploy.yaml
	oc apply -f $(grafana_secret_file)
	oc create -n hongkliu-stage configmap grafana-config --from-file=./deploy/testctl_http/grafana.ini --dry-run -o yaml | oc apply -f -
	oc create -n hongkliu-stage configmap grafana-datasources --from-file=.datasources.yaml=./deploy/testctl_http/grafana_datasources.yaml --dry-run -o yaml | oc apply -f -
	oc create -n hongkliu-stage configmap grafana-dashboards --from-file=dashboards.yaml=./deploy/testctl_http/grafana_dashboards.yaml --dry-run -o yaml | oc apply -f -
	oc create -n hongkliu-stage configmap grafana-dashboard-test-go --from-file=test-go.json=./deploy/testctl_http/test_go_dashboard.json --dry-run -o yaml | oc apply -f -
	oc apply -f ./deploy/testctl_http/grafana_deploy.yaml
	sed -e "s|{slack_api_secret}|$(slack_api_secret_value)|g" ./deploy/testctl_http/alertmanager.yml > /tmp/alertmanager_decoded.yml
	oc create -n hongkliu-stage configmap alert-manager-config --from-file=alertmanager.yml=/tmp/alertmanager_decoded.yml --from-file=msg.tmpl=./deploy/testctl_http/alert_manager.msg.tmpl --dry-run -o yaml | oc apply -f -
	rm -vf /tmp/alertmanager_decoded.yml
	oc apply -f ./deploy/testctl_http/alert_manager_deploy.yaml
	#https://github.com/kubernetes/kubernetes/issues/13488#issuecomment-481023838
	#kubectl rollout restart #this will be available soon
	@echo "deployed successfully!"
