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