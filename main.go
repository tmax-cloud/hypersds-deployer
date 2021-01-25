package main

import (
	provisioner "hypersds-provisioner/cmd/hypersds-provisioner"
    hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
    "hypersds-provisioner/pkg/util"
    "hypersds-provisioner/pkg/common/wrapper"
	"os"
    "fmt"
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
