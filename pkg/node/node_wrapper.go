package node

import (
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
)

type NewNodeInterface interface {
	NewNodesFromCephCr(cephSpec hypersdsv1alpha1.CephClusterSpec) ([]NodeInterface, error)
}

var NewNodeWrapper NewNodeInterface

func init() {
	NewNodeWrapper = &NewNodeStruct{}
}
