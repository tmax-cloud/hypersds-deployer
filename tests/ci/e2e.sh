#!/bin/bash

function minikube_start() {
	echo "--------- Launch k8s cluster by minikube ---------"
	curl -Lo $1/minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
	chmod +x $1/minikube
	# working path is $GITHUB_WORKSPACE (e.g. /home/runner/hypersds-provisioner/hypsersds-provisioner/)
	$1/minikube start --kubernetes-version v1.19.4 --mount --mount-string=$1"/tests/e2e/:/e2e/"
	minikube update-context
}

function minikube_delete() {
	echo "--------- Remove k8s minikube cluster ---------"
	$1/minikube delete
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
	make container
}

function e2e_test() {
	echo "--------- Test all e2e testcases ---------"
	ginkgo -v $1"/tests/e2e"
}

function how_to_use() {
	echo "$0 <op> {param1} {param2}...

Available Operations:
k8s_up <base_directory>			launch minikube cluster with mounting '<base_directory>/tests/e2e/'
k8s_down <base_directory>		remove k8s cluster launched by minikube at <base_directory>
cluster_up <base_directory>		launch cluster to install ceph by '<base_directory/tests/ci/Vagrantfile'
cluster_destroy <base_directory>	remove cluster launched by '<base_directory/tests/ci/Vagrantfile'
build_image <base_directory>		build pod image
test <base_directory>			test all e2e testcases in '<base_directory>/tests/e2e/'
" >&2
}

case "${1:-}" in
k8s_up)
	minikube_start $2
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
	build_image $2
	;;
test)
	e2e_test $2
	;;
*)
	how_to_use
;;
esac
