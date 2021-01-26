package provisioner

import (
	"hypersds-provisioner/pkg/common/wrapper"
	"hypersds-provisioner/pkg/util"
	"os"
)

var (
	//// Test
	hostName = "tmax"
	hostAddr = "192.168.7.19"
)

func Run() error {
	testCommand := []string{"ls", "-alh"}
	output, err := util.RunSSHCmd(wrapper.ExecWrapper, hostName, hostAddr, testCommand...)

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
