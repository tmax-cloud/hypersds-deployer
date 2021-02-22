#!/bin/bash

# working path is $GITHUB_WORKSPACE (e.g. /home/runner/hypersds-provisioner/hypsersds-provisioner/)
minikube start --kubernetes-version v1.19.4 --mount --mount-string="./tests/e2e/:/e2e/"
eval $(minikube docker-env)
make build
make container
VAGRANT_VAGRANTFILE=./tests/ci/Vagrantfile vagrant up
ginkgo ./tests/e2e
