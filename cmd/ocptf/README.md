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

```bash
$ buildah bud --format=docker -f test_files/docker/Dockerfile.ocptf.txt -t quay.io/hongkailiu/test-go:ocptf-0.0.1 .
$ buildah push --creds=hongkailiu d58cbf2a06aa docker://quay.io/hongkailiu/test-go:ocptf-0.0.1

```