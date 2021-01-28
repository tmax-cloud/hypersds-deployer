package node

import (
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
)

type NodeInitInterface interface {
	NewNodesFromCephCr(cephSpec hypersdsv1alpha1.CephClusterSpec) ([]NodeInterface, error)
}

type NodeInitStruct struct {
}

// TODO: error handling
func (nis *NodeInitStruct) NewNodesFromCephCr(cephSpec hypersdsv1alpha1.CephClusterSpec) ([]NodeInterface, error) {
	var nodes []NodeInterface

	for _, nodeInCephSpec := range cephSpec.Nodes {
		var n Node
		err := n.SetUserId(nodeInCephSpec.UserID)
		if err != nil {
			panic(err)
		}
		err = n.SetUserPw(nodeInCephSpec.Password)
		if err != nil {
			panic(err)
		}

		var hostSpec HostSpec
		err = hostSpec.SetServiceType()
		if err != nil {
			panic(err)
		}

		err = hostSpec.SetAddr(nodeInCephSpec.IP)
		if err != nil {
			panic(err)
		}

		err = hostSpec.SetHostName(nodeInCephSpec.HostName)
		if err != nil {
			panic(err)
		}

		err = n.SetHostSpec(&hostSpec)
		if err != nil {
			panic(err)
		}

		nodes = append(nodes, &n)
	}

	return nodes, nil
}

var NodeInitWrapper NodeInitInterface

func init() {
	NodeInitWrapper = &NodeInitStruct{}
}
