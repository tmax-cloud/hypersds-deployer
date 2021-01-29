package util

import (
	"hypersds-provisioner/pkg/common/wrapper"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
)

const (
	FilePathPrefix = "/root/"
	CephCrYamlName = "cluster.yaml"
)

func getCephClusterSpec(ioutil wrapper.IoUtilInterface, yaml wrapper.YamlInterface) (hypersdsv1alpha1.CephClusterSpec, error) {
	fileName := FilePathPrefix + CephCrYamlName
	source, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	var t hypersdsv1alpha1.CephClusterSpec
	err = yaml.Unmarshal(source, &t)
	if err != nil {
		panic(err)
	}
	return t, nil
}
