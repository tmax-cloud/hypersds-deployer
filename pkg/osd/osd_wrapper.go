package osd

import (
	"hypersds-provisioner/pkg/common/wrapper"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
)

type OsdStructInterface interface {
	NewOsdsFromCephCr(cephSpec hypersdsv1alpha1.CephClusterSpec) ([]*Osd, error)
	NewOsdsFromCephOrch(yaml wrapper.YamlInterface, rawOsdsFromOrch []byte) ([]*Osd, error)
}

var OsdWrapper OsdStructInterface

func init() {
	OsdWrapper = &osdStruct{}
}
