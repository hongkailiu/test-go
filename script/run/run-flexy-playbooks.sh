#!/bin/bash

set -e

echo "secret_vars_path: ${secret_vars_file}"
echo "ansible_repo_path: ${ansible_repo_path}"
echo "install_glusterfs: ${install_glusterfs}"

if [[ -z "${secret_vars_file}" ]]; then
  echo "${secret_vars_file} is required"
  exit 1
fi

if [[ -z "${ansible_repo_path}" ]]; then
  echo "${ansible_repo_path} is required"
  exit 1
fi

if [[ -z "${install_glusterfs}" ]]; then
  echo "${install_glusterfs} is required"
  exit 1
fi

source "${secret_vars_file}"

if [[ -z "${REG_AUTH_USER}" ]]; then
  echo "REG_AUTH_USER is not defined in ${secret_vars_file}"
  exit 1
fi

ansible-playbook -i build/output/flexy/inv/2.file ${ansible_repo_path}/playbooks/prerequisites.yml
ansible-playbook -i build/output/flexy/inv/2.file ${ansible_repo_path}/playbooks/deploy_cluster.yml
#ansible-playbook -i build/output/flexy/inv/2.file ${ansible_repo_path}/playbooks/olm/config.yml

if [[ $(echo "${install_glusterfs}" | awk '{print tolower($0)}') = "true" ]]; then
  ansible-playbook -i build/output/flexy/inv/ ${ansible_repo_path}/playbooks/openshift-glusterfs/config.yml
fi