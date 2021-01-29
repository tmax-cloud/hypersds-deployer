package main

import (
	"fmt"
	provisioner "hypersds-provisioner/cmd/hypersds-provisioner"
	"hypersds-provisioner/pkg/common/wrapper"
	"hypersds-provisioner/pkg/util"
	"os"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
)

func main() {

	var t hypersdsv1alpha1.CephClusterSpec
	t, _ = util.UtilWrapper.GetCephClusterSpec(wrapper.IoUtilWrapper, wrapper.YamlWrapper)
	fmt.Println(t)
	err := provisioner.Run()
	if err != nil {
		os.Exit(1)
	}
}
