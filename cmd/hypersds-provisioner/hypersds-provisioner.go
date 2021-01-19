package provisioner

import (
	"fmt"
	"hypersds-provisioner/pkg/util"
	"os"
)

var (
	//// Test
	hostName = "tmax"
	hostAddr = "192.168.7.19"
)

func Run() error {
	////Test2
	util.TestStruct()
	//// Test
	testCommand := []string{"ls", "-alh"}
	fmt.Printf("qqqqq")
	output, err := util.RunSSHCmd(util.ExecWrapper, hostName, hostAddr, testCommand...)

	if err != nil {
		output.WriteTo(os.Stderr)
		return err
	} else {
		output.WriteTo(os.Stdout)
	}

	return nil
}
