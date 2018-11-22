# ocptf

## build

```bash
$ make build-ocptf

```

## run

```bash
$ export terraform_tf_state_file=$(readlink -f ./test_files/ocpft/unit.test.files/terraform.tfstate.json) 
$ ansible-playbook -i build/ocptf -i test_files/ocpft/inv/2.file ./test_files/ocpft/playbook/test.yaml
$ ansible-playbook -i build/ocptf -i test_files/ocpft/inv/2.file ./test_files/ocpft/playbook/test.yaml --list-hosts

```

## Images

Build images:

```bash
$ buildah bud --format=docker -f test_files/docker/Dockerfile.ocptf.txt -t quay.io/hongkailiu/test-go:ocptf-0.0.1 .
$ buildah push --creds=hongkailiu d58cbf2a06aa docker://quay.io/hongkailiu/test-go:ocptf-0.0.1

```

Use an existing image:

```bash
# podman run --rm -t -i quay.io/hongkailiu/test-go:ocptf-0.0.2 bash 
[root@c54a44b23f3f /]# git clone https://github.com/hongkailiu/svt-case-doc.git
[root@c54a44b23f3f /]# cd svt-case-doc/files/terraform/4_node_cluster/
### update the secret
[root@c54a44b23f3f 4_node_cluster]# vi secret.tfvars 
# terraform init -var-file="secret.tfvars"
# terraform apply -var-file="secret.tfvars" -auto-approve
# export terraform_tf_state_file=$(readlink -f ./terraform.tfstate)
### 
# cd /
[root@c54a44b23f3f /]# export ANSIBLE_HOST_KEY_CHECKING=False
###
# vi /perf.key
# chmod 0600 /perf.key
### (optional) test run for authentication and dynamic inventory
# ansible-playbook -i /bin/ocptf -i /inv/2.file /playbook/test.yaml -e "ansible_ssh_private_key_file=/perf.key"
###
# git clone https://github.com/openshift/openshift-ansible.git
# cd openshift-ansible/
# git checkout release-3.11

### export secrets for the playbook

# cd /
[root@77c9bf718b93 /]# install_ocp_gluster=false ansible-playbook -i /bin/ocptf -i /inv/2.file /openshift-ansible/playbooks/prerequisites.yml -e "ansible_ssh_private_key_file=/perf.key"
# install_ocp_gluster=false ansible-playbook -i /bin/ocptf -i /inv/2.file /openshift-ansible/playbooks/deploy_cluster.yml -e "ansible_ssh_private_key_file=/perf.key"
### optional: install glusterfs
# ansible-playbook -i /bin/ocptf -i /inv/ /openshift-ansible/playbooks/openshift-glusterfs/config.yml -e "ansible_ssh_private_key_file=/perf.key"

### clean up
# cd /
[root@5f82b77cab38 /]# cd svt-case-doc/files/terraform/4_node_cluster/
# terraform destroy -var-file="secret.tfvars" -auto-approve


```
