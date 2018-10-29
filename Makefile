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
	swagger validate ./swagger/swagger/swagger.yml


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
		--target=./swagger/swagger \
		--spec=./swagger/swagger/swagger.yml \
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
