# Flexy

# Run with pre-build binary

```bash
### install ginkgo
$ go get github.com/onsi/ginkgo/ginkgo
### download the binary, see more versions: https://github.com/cduser/svt-release/branches
$ curl -LO https://github.com/cduser/svt-release/raw/travis_flexy_13/flexy-afd8da2-Linux-x86_64.tar.gz
$ tar -xvf flexy-afd8da2-Linux-x86_64.tar.gz
$ cd flexy/
### make sure aws credentials is working
$ cat ~/.aws/credentials
$ make create-flexy-instances
### make sure aaa.sh has the right info
### $cat /tmp/aaa.sh
### export REG_AUTH_USER="aos-qe-pull36"
### export REG_AUTH_PASSWORD="aaa"
### export AWS_ACCESS_KEY_ID="bbb"
### export AWS_SECRET_ACCESS_KEY="ccc"
$ make run-flexy-playbooks secret_vars_file=/tmp/aaa.sh ansible_repo_path=/root/openshift-ansible install_glusterfs=false

```
 