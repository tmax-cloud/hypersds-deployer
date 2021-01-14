package util

import (
	"os/exec"
	"time"
	"bytes"
	"context"
)

const (
	cmdTimeout = 20*time.Minute
)

func RunSSHCmd(hostName, hostAddr string, cephQuery ...string) (bytes.Buffer, error) {
	var resultStdout, resultStderr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()

	hostInfo := hostName + "@" + hostAddr

	/// Add password option to sshCmd
	sshCmd := []string{hostInfo}
	sshCmd = append(sshCmd, cephQuery...)

	cmd := exec.CommandContext(ctx, "ssh", sshCmd...)
	cmd.Stdout = &resultStdout
	cmd.Stderr = &resultStderr

	err := cmd.Run()

	if err != nil {
		return resultStderr, err
	}

	return resultStdout, nil
}
