#!/bin/bash

function minikube_start() {
	echo "--------- Launch k8s cluster by minikube ---------"
	# working path is $GITHUB_WORKSPACE (e.g. /home/runner_account/hypersds-provisioner/hypsersds-provisioner/)
	curl -Lo $1/minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
	chmod +x $1/minikube
	if [ -z $2 ]
	then
	$1/minikube start --kubernetes-version v1.19.4 --mount --mount-string=$1"/tests/e2e/:/e2e/"
	$1/minikube update-context
	else
	$1/minikube start --insecure-registry $2 --kubernetes-version v1.19.4 --mount --mount-string=$1"/tests/e2e/:/e2e/"
	$1/minikube update-context
	minikube_set_registry $1 $2 $3 $4
	fi
}

function minikube_set_registry() {
	echo "--------- Set docker private registry on minikube ---------"
	$1/minikube ssh "echo $4 | docker login $2 --username $3 --password-stdin"
	$1/minikube ssh cat .docker/config.json > $1/_minikube_registry_config.json
	kubectl create secret generic regcred --from-file=.dockerconfigjson=$1/_minikube_registry_config.json --type=kubernetes.io/dockerconfigjson
    # activate registry and hostpath dir of e2e input files
    sed -i 's/\# registryCredentialName/registryCredentialName/g' $1/tests/e2e/inputs/*.yaml
    sed -i 's/\# testManifestDir/testManifestDir/g' $1/tests/e2e/inputs/*.yaml
}

function minikube_delete() {
	echo "--------- Remove k8s minikube cluster ---------"
	$1/minikube delete
	rm $1/_minikube_registry_config.json
	rm $1/minikube
}

function vagrant_up() {
	echo "--------- Launch VM cluster by vagrant ---------"
	VAGRANT_VAGRANTFILE=$1"/tests/ci/Vagrantfile" vagrant up
}

function vagrant_destroy() {
	echo "--------- Remove ceph vagrant cluster ---------"
	VAGRANT_VAGRANTFILE=$1"/tests/ci/Vagrantfile" vagrant destroy -f
}

function build_image() {
	echo "--------- Build hypersds-provisioner docker image ---------"
	make build
	eval $($1/minikube docker-env)
	if [ -z $2 ]
	then
	make container
	else
	echo $4 | docker login $2 --username $3 --password-stdin
	make container REGISTRY=$2
	fi
}

function e2e_test() {
	echo "--------- Test all e2e testcases ---------"
	ginkgo -v $1"/tests/e2e"
}

function how_to_use() {
	echo "$0 <op> {param1} {param2}...

Available Operations:
k8s_up <base_dir> [reg_endpoint reg_id reg_pw]	Launch minikube cluster with mounting '<base_directory>/tests/e2e/'
						Activate private registry reg_endpoint with reg_id and reg_pw if they exist

k8s_down <base_dir>				Remove k8s cluster launched by minikube at <base_directory>

cluster_up <base_dir>				Launch cluster to install ceph by '<base_directory/tests/ci/Vagrantfile'

cluster_destroy <base_dir>			Remove cluster launched by '<base_directory/tests/ci/Vagrantfile'

build_image <base_dir> [reg_endpoint]		Build pod image
						Push built image to private registry if reg_endpoint exists

test <base_dir>					Test all e2e testcases in '<base_directory>/tests/e2e/'
" >&2
}

case "${1:-}" in
k8s_up)
	minikube_start $2 $3 $4 $5
	;;
k8s_down)
	minikube_delete $2
	;;
cluster_up)
	vagrant_up $2
	;;
cluster_down)
	vagrant_destroy $2
	;;
build_image)
	build_image $2 $3 $4 $5
	;;
test)
	e2e_test $2
	;;
*)
	how_to_use
;;
esac
