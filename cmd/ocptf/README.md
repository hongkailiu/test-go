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
```

### Dynamic inventory

```bash
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

### Known issue

* The above ansible-playbook for glusterfs installation does not work. It is up to the fact that
the quotes in `glusterfs_devices` cannot be handled correctly.

```bash
"ec2-54-202-76-127.us-west-2.compute.amazonaws.com": {
        "glusterfs_devices": "'[\"/dev/nvme2n1\"]'",
        "openshift_node_group_name": "node-config-compute",
        "openshift_public_hostname": "ec2-54-202-76-127.us-west-2.compute.amazonaws.com"
      }

``` 

It failed here:

```bash
TASK [openshift_storage_glusterfs : Load heketi topology] *********************************************************************
fatal: [ec2-34-209-244-212.us-west-2.compute.amazonaws.com]: FAILED! => {"changed": true, "cmd": ["oc", "--config=/tmp/openshift-glusterfs-ansible-r03NKl/admin.kubeconfig", "rsh", "--namespace=glusterfs", "deploy-heketi-storage-1-z6w6b", "heketi-cli", "-s", "http://localhost:8080", "--user", "admin", "--secret", "/qhUFBHAoFSAsOGlrDECDsTzUWLdLShtjsu1252qhpc=", "topology", "load", "--json=/tmp/openshift-glusterfs-ansible-r03NKl/topology.json", "2>&1"], "delta": "0:00:00.298894", "end": "2018-11-22 20:56:09.184860", "failed_when_result": true, "msg": "non-zero return code", "rc": 255, "start": "2018-11-22 20:56:08.885966", "stderr": "Error: Unable to parse config file\ncommand terminated with exit code 255", "stderr_lines": ["Error: Unable to parse config file", "command terminated with exit code 255"], "stdout": "", "stdout_lines": []}
	to retry, use: --limit @/openshift-ansible/playbooks/openshift-glusterfs/config.retry


```

The workaround: Manually copy/paste the `glusterfs group`.

```bash
# vi /inv/gfs.file
[glusterfs]
ec2-34-218-254-163.us-west-2.compute.amazonaws.com glusterfs_devices='["/dev/nvme2n1"]'
ec2-54-184-86-92.us-west-2.compute.amazonaws.com glusterfs_devices='["/dev/nvme2n1"]'
ec2-54-202-76-127.us-west-2.compute.amazonaws.com glusterfs_devices='["/dev/nvme2n1"]'

# install_ocp_gluster=false ansible-playbook -i /bin/ocptf -i /inv/ /openshift-ansible/playbooks/openshift-glusterfs/config.yml -e "ansible_ssh_private_key_file=/perf.key"

```

### Static inventory

```bash
### generate static inventory
# install_ocp_gluster=false /bin/ocptf --list --static > /inv/ocptf.file
# ansible-playbook  -i /inv/ /openshift-ansible/playbooks/prerequisites.yml -e "ansible_ssh_private_key_file=/perf.key"
# ansible-playbook  -i /inv/ /openshift-ansible/playbooks/deploy_cluster.yml -e "ansible_ssh_private_key_file=/perf.key"
### optional: install glusterfs
# /bin/ocptf --list --static > /inv/ocptf.file
# ansible-playbook  -i /inv/ /openshift-ansible/playbooks/openshift-glusterfs/config.yml -e "ansible_ssh_private_key_file=/perf.key"

```