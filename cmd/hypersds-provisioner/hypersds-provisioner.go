package provisioner

import (
	"hypersds-provisioner/pkg/util"
	"os"
)

var (
	//// Test
	hostName = "k8s"
	hostAddr = "192.168.50.90"
)

func Run() error {
	//// Test
	cephcluster, err := util.ParseYaml()
	if err != nil {
		return err
	}

	err = util.GenerateConfFile(cephcluster)
	if err != nil {
		return err
	}

	testCommand := []string{"ls", "-alh"}

	output, err := util.RunSSHCmd(hostName, hostAddr, testCommand...)

	if err != nil {
		output.WriteTo(os.Stderr)
		return err
	} else {
		output.WriteTo(os.Stdout)
	}
	

	return nil
}
