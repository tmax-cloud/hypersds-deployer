package util

import (
	"hypersds-provisioner/pkg/common/wrapper"
    hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
)

type UtilInterface interface {
    GetCephClusterSpec(ioutil wrapper.IoUtilInterface, yaml wrapper.YamlInterface) (hypersdsv1alpha1.CephClusterSpec, error)
}

type utilStruct struct {
}

func (u *utilStruct) GetCephClusterSpec(ioutil wrapper.IoUtilInterface, yaml wrapper.YamlInterface) (hypersdsv1alpha1.CephClusterSpec, error){
    return getCephClusterSpec(ioutil, yaml)
}

var UtilWrapper UtilInterface

func init() {
	UtilWrapper = &utilStruct{}
}
