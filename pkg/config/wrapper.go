package config

import hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"

type ConfigInitInterface interface {
	NewConfigFromCephCr(hypersdsv1alpha1.CephClusterSpec) (*CephConfig, error)
}

var ConfigInitWrapper ConfigInitInterface

func init() {
	ConfigInitWrapper = &ConfigInitStruct{}
}
