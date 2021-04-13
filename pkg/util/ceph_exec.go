package util

import (
	"bytes"
	"context"
	"hypersds-provisioner/pkg/common/wrapper"
	"time"
)

const (
	CephExecCmdTimeout = 1 * time.Minute
	CephCmdTimeout     = "50"
)

func RunCephCmd(exec wrapper.ExecInterface, cmdQuery ...string) (bytes.Buffer, error) {
	cmdQuery = append(cmdQuery, "--connect-timeout", CephCmdTimeout)

	var resultStdout, resultStderr bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), CephExecCmdTimeout)
	defer cancel()

	err := exec.CommandExecute(&resultStdout, &resultStderr, ctx, "ceph", cmdQuery...)

	if err != nil {
		return resultStderr, err
	}

	return resultStdout, nil
}
