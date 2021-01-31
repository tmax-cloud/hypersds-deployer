package node

import (
	common "hypersds-provisioner/pkg/common/wrapper"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"

	"bytes"
	"context"
)

type NodeInterface interface {
	SetUserId(userId string) error
	SetUserPw(userPw string) error
	SetHostSpec(hostSpec HostSpecInterface) error

	GetUserId() (string, error)
	GetUserPw() (string, error)
	GetHostSpec() (HostSpecInterface, error)

	RunSshCmd(exec common.ExecInterface, cmdQuery string) (bytes.Buffer, error)
}

type Node struct {
	userId   string
	userPw   string
	hostSpec HostSpecInterface
}

type NewNodeStruct struct {
}

func (n *Node) SetUserId(userId string) error {
	n.userId = userId
	return nil
}

func (n *Node) SetUserPw(userPw string) error {
	n.userPw = userPw
	return nil
}

func (n *Node) SetHostSpec(hostSpec HostSpecInterface) error {
	n.hostSpec = hostSpec
	return nil
}

func (n *Node) GetUserId() (string, error) {
	return n.userId, nil
}

func (n *Node) GetUserPw() (string, error) {
	return n.userPw, nil
}

func (n *Node) GetHostSpec() (HostSpecInterface, error) {
	return n.hostSpec, nil
}

// executing commands: sshpass -f <(printf '%s\n' userPw) ssh userId@ipAddr cmdQuery
// TODO: replace ssh commands to go ssh pkg
func (n *Node) RunSshCmd(exec common.ExecInterface, cmdQuery string) (bytes.Buffer, error) {
	var resultStdout, resultStderr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), SshCmdTimeout)
	defer cancel()

	userPw, err := n.GetUserPw()
	if err != nil {
		panic(err)
	}

	sshCmd := "sshpass -f <(printf '%s\\n' " + userPw + ") "

	userId, err := n.GetUserId()
	if err != nil {
		panic(err)
	}
	nodeHostSpec, err := n.GetHostSpec()
	if err != nil {
		panic(err)
	}
	ipAddr, err := nodeHostSpec.GetAddr()
	if err != nil {
		panic(err)
	}

	hostInfo := userId + "@" + ipAddr
	sshKeyCheckOpt := "-oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null "
	sshCmd += "ssh " + sshKeyCheckOpt + hostInfo + " " + cmdQuery

	parameters := []string{"-c"}
	parameters = append(parameters, sshCmd)
	err = exec.CommandExecute(&resultStdout, &resultStderr, ctx, "bash", parameters...)

	if err != nil {
		return resultStderr, err
	}

	return resultStdout, nil
}

// TODO: error handling
func (newNodeStruct *NewNodeStruct) NewNodesFromCephCr(cephSpec hypersdsv1alpha1.CephClusterSpec) ([]NodeInterface, error) {
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
