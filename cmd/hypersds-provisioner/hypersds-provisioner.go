package provisioner

import (
	"hypersds-provisioner/pkg/common/wrapper"
	node "hypersds-provisioner/pkg/node"
	"os"
)

var (
	//// Test
	userId   = "k8s"
	userPw   = "k8s"
	ipAddr   = "192.168.50.92"
	hostName = "worker2"
)

func Run() error {
	n1 := node.Node{}
	hostSpec1 := node.HostSpec{
		ServiceType: "host",
		Addr:        ipAddr,
		HostName:    hostName,
	}
	err := n1.SetUserId(userId)
	if err != nil {
		panic(err)
	}
	err = n1.SetUserPw(userPw)
	if err != nil {
		panic(err)
	}
	err = n1.SetHostSpec(&hostSpec1)
	if err != nil {
		panic(err)
	}

	testCommand := "ifconfig"
	output, err := n1.RunSshCmd(wrapper.ExecWrapper, testCommand)

	if err != nil {
		_, err2 := output.WriteTo(os.Stderr)
		if err2 != nil {
			return err2
		}

		return err
	} else {
		_, err2 := output.WriteTo(os.Stdout)
		if err2 != nil {
			return err2
		}
	}

	return nil
}
