package node

import (
	common "hypersds-provisioner/pkg/common/wrapper"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"

	"bytes"
	"context"
	"fmt"
)

type NodeInterface interface {
	SetUserId(userId string) error
	SetUserPw(userPw string) error
	SetHostSpec(hostSpec HostSpecInterface) error

	GetUserId() (string, error)
	GetUserPw() (string, error)
	GetHostSpec() (HostSpecInterface, error)

	RunSshCmd(exec common.ExecInterface, cmdQuery string) (bytes.Buffer, error)

	// if role is DEST, copy file from container to this node
	// if role is SRC, copy file from this node to container
	RunScpCmd(exec common.ExecInterface, srcFile, destFile string, role Role) (bytes.Buffer, error)
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

// executing commands: sshpass -f <(printf '%s\n' userPw) ssh -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null userId@ipAddr cmdQuery
// TODO: replace sshpass command to go ssh pkg
func (n *Node) RunSshCmd(exec common.ExecInterface, cmdQuery string) (bytes.Buffer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), SshCmdTimeout)
	defer cancel()

	userPw, err := n.GetUserPw()
	if err != nil {
		return bytes.Buffer{}, err
	}

	userId, err := n.GetUserId()
	if err != nil {
		return bytes.Buffer{}, err
	}
	nodeHostSpec, err := n.GetHostSpec()
	if err != nil {
		return bytes.Buffer{}, err
	}
	ipAddr, err := nodeHostSpec.GetAddr()
	if err != nil {
		return bytes.Buffer{}, err
	}

	const sshKeyCheckOpt = "-oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null"
	sshCmd := fmt.Sprintf("sshpass -f <(printf '%%s\\n' %[1]s) ssh %[2]s %[3]s@%[4]s '%[5]s'", userPw, sshKeyCheckOpt, userId, ipAddr, cmdQuery)
	parameters := []string{"-c"}
	parameters = append(parameters, sshCmd)

	var resultStdout, resultStderr bytes.Buffer
	err = exec.CommandExecute(&resultStdout, &resultStderr, ctx, "bash", parameters...)

	if err != nil {
		return resultStderr, err
	}

	return resultStdout, nil
}

// TODO: replace sshpass command to go ssh pkg
func (n *Node) RunScpCmd(exec common.ExecInterface, srcFile, destFile string, role Role) (bytes.Buffer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), SshCmdTimeout)
	defer cancel()

	userPw, err := n.GetUserPw()
	if err != nil {
		return bytes.Buffer{}, err
	}

	userId, err := n.GetUserId()
	if err != nil {
		return bytes.Buffer{}, err
	}
	nodeHostSpec, err := n.GetHostSpec()
	if err != nil {
		return bytes.Buffer{}, err
	}
	ipAddr, err := nodeHostSpec.GetAddr()
	if err != nil {
		return bytes.Buffer{}, err
	}

	const sshKeyCheckOpt = "-oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null"

	var scpCmd string
	// provisioner sends srcFile to this node as destFile
	if role == DESTINATION {
		scpCmd = fmt.Sprintf("sshpass -f <(printf '%%s\\n' %[1]s) scp %[2]s %[3]s %[4]s@%[5]s:/%[6]s",
			userPw, sshKeyCheckOpt, srcFile, userId, ipAddr, destFile)

		// this node sends srcFile to provisioner as destFile
	} else {
		scpCmd = fmt.Sprintf("sshpass -f <(printf '%%s\\n' %[1]s) scp %[2]s %[4]s@%[5]s:%[3]s %[6]s",
			userPw, sshKeyCheckOpt, srcFile, userId, ipAddr, destFile)
	}

	parameters := []string{"-c"}
	parameters = append(parameters, scpCmd)

	var resultStdout, resultStderr bytes.Buffer
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
			return nil, err
		}
		err = n.SetUserPw(nodeInCephSpec.Password)
		if err != nil {
			return nil, err
		}

		var hostSpec HostSpec
		err = hostSpec.SetServiceType()
		if err != nil {
			return nil, err
		}

		err = hostSpec.SetAddr(nodeInCephSpec.IP)
		if err != nil {
			return nil, err
		}

		err = hostSpec.SetHostName(nodeInCephSpec.HostName)
		if err != nil {
			return nil, err
		}

		err = n.SetHostSpec(&hostSpec)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, &n)
	}

	return nodes, nil
}
