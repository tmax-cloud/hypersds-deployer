package util

import (
	"bytes"
	"context"
	"time"
    "hypersds-provisioner/pkg/common/wrapper"
)

const (
	cmdTimeout = 20 * time.Minute
)

// unit test를 하기 위해서는 function 내 다른 package function 호출, package 내 다른 function 호출을 unit test시 가짜 function으로 대체할 수 있게 mocking이 가능한 구조로 구현 필요
// RunSSHCmd에서는 다른 package function인 exec.CommandContext를 호출하여 cmd를 받고, cmd.Run()을 통해 최종적으로 exec 실행
// 다른 package의 function을 unit test시 대체하기 위해 호출할 function을 인자로 받아서 구조로 구현 필요
// function에서 다른 package function 호출할 경우, interface를 인자로 받아 interface를 function으로 호출, interface 구현은 interface.go 참조
func RunSSHCmd(exec wrapper.ExecInterface, hostName, hostAddr string, cephQuery ...string) (bytes.Buffer, error) {
	// exec interface  인자로 받음
	var resultStdout, resultStderr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()

	hostInfo := hostName + "@" + hostAddr

	/// Add password option to sshCmd
	sshCmd := []string{hostInfo}
	sshCmd = append(sshCmd, cephQuery...)
	// exec interface 의 function 호출
	err := exec.CommandExecute(&resultStdout, &resultStderr, ctx, "ssh", sshCmd...)

	if err != nil {
		return resultStderr, err
	}

	return resultStdout, nil
}
