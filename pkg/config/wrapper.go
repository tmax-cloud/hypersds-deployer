package config

import hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"

type ConfigInitInterface interface {
	NewConfigFromCephCr(hypersdsv1alpha1.CephClusterSpec) *CephConfig
}

type configInitStruct struct{}

func (c *configInitStruct) NewConfigFromCephCr(cr hypersdsv1alpha1.CephClusterSpec) *CephConfig {
	return NewConfigFromCephCr(cr)
}

var ConfigInitWrapper ConfigInitInterface

func init() {
	ConfigInitWrapper = &configInitStruct{}
}
