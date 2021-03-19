package provisioner

import (
	"bytes"
	"fmt"
	"os"

	common "hypersds-provisioner/pkg/common/wrapper"
	node "hypersds-provisioner/pkg/node"
)

func updateCephClusterToOp() error {
	fmt.Println("\n----------------Start to update conf and keyring to operator---------------")

	const pathConfToUpdate = pathConfigWorkingDir + cephConfToUpdate
	err = cephConfig.ConfigFromAdm(common.IoUtilWrapper, pathConfToUpdate)
	if err != nil {
		return err
	}

	const pathKeyringToUpdate = pathConfigWorkingDir + cephKeyringToUpdate
	err = cephConfig.SecretFromAdm(common.IoUtilWrapper, pathKeyringToUpdate)
	if err != nil {
		return err
	}

	err = cephConfig.UpdateConfToK8s(common.KubeWrapper)
	if err != nil {
		return err
	}

	err = cephConfig.UpdateKeyringToK8s(common.KubeWrapper)

	return err
}

func bootstrapCephadm(targetNode node.NodeInterface) error {
	fmt.Println("\n----------------Start to bootstrap ceph---------------")

	fmt.Println("[bootstrapCephadm] copying conf file")
	const pathConf = pathConfigWorkingDir + cephConfNameFromCr
	err = copyFile(targetNode, node.DESTINATION, pathConf, pathConf)
	if err != nil {
		return err
	}

	deployNodeHostSpec, err := targetNode.GetHostSpec()
	if err != nil {
		return err
	}

	monIp, err := deployNodeHostSpec.GetAddr()
	if err != nil {
		return err
	}

	fmt.Println("[bootstrapCephadm] executing bootstrap")
	admBootstrapCmd := fmt.Sprintf("cephadm bootstrap --mon-ip %s --config %s", monIp, pathConf)
	err = processCmdOnNode(targetNode, admBootstrapCmd)
	if err != nil {
		return err
	}

	fmt.Println("[bootstrapCephadm] checking status")
	const admHealthCheckCmd = "cephadm shell -- ceph -s"
	err = processCmdOnNode(targetNode, admHealthCheckCmd)

	return err
}

func installCephadm(targetNode node.NodeInterface) error {
	fmt.Println("\n----------------Start to install cephadm---------------")

	// TODO: Specify release version
	fmt.Println("[installCephadm] executing curl cephadm")
	curlCephadmCmd := "curl --silent --remote-name --location https://github.com/ceph/ceph/raw/octopus/src/cephadm/cephadm"
	err = processCmdOnNode(targetNode, curlCephadmCmd)
	if err != nil {
		return err
	}

	fmt.Println("[installCephadm] executing chmod")
	const chmodCmd = "chmod +x cephadm"
	err = processCmdOnNode(targetNode, chmodCmd)
	if err != nil {
		return err
	}

	// TODO: Specify release version
	fmt.Println("[installCephadm] executing cephadm add-repo")
	admAddRepoCmd := "./cephadm add-repo --release octopus"
	err = processCmdOnNode(targetNode, admAddRepoCmd)
	if err != nil {
		return err
	}

	// does something in command need to be changed, related to cephadm version?
	fmt.Println("[installCephadm] executing curl cephadm gpg key")
	addCephadmRepoCmd := "curl https://download.ceph.com/keys/release.asc | gpg --no-default-keyring --keyring /tmp/fix.gpg --import - && gpg --no-default-keyring --keyring /tmp/fix.gpg --export > /etc/apt/trusted.gpg.d/ceph.release.gpg && rm /tmp/fix.gpg"
	err = processCmdOnNode(targetNode, addCephadmRepoCmd)
	if err != nil {
		return err
	}

	fmt.Println("[installCephadm] executing cephadm apt-get update")
	const aptUpdateCmd = "apt-get update"
	err = processCmdOnNode(targetNode, aptUpdateCmd)
	if err != nil {
		return err
	}

	fmt.Println("[installCephadm] executing cephadm install")
	const admInstallCmd = "./cephadm install"
	err = processCmdOnNode(targetNode, admInstallCmd)
	if err != nil {
		return err
	}

	fmt.Println("[installCephadm] executing mkdir")
	const mkdirCmd = "mkdir /etc/ceph"
	err = processCmdOnNode(targetNode, mkdirCmd)

	return err
}

func installBasePackage(targetNodeList []node.NodeInterface) error {
	fmt.Println("\n----------------Start to install base package---------------")

	fmt.Println("[installBasePackage] executing apt-get update")
	const aptUpdateCmd = "apt-get update"
	for _, n := range targetNodeList {
		err = processCmdOnNode(n, aptUpdateCmd)
		if err != nil {
			return err
		}
	}

	// use standard verison in OS
	fmt.Println("[installBasePackage] executing apt-get install ...")
	const installPkgCmd = "apt-get install -y apt-transport-https ca-certificates curl software-properties-common ntpdate chrony"
	for _, n := range targetNodeList {
		err = processCmdOnNode(n, installPkgCmd)
		if err != nil {
			return err
		}
	}

	fmt.Println("[installBasePackage] executing curl docker ...")
	const addDockerGpgKeyCmd = "curl -s https://download.docker.com/linux/ubuntu/gpg | apt-key add - &>/dev/null"
	for _, n := range targetNodeList {
		err = processCmdOnNode(n, addDockerGpgKeyCmd)
		if err != nil {
			return err
		}
	}

	fmt.Println("[installBasePackage] executing add-apt-repo docker ...")
	const addDockerRepoCmd = `add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"`
	for _, n := range targetNodeList {
		err = processCmdOnNode(n, addDockerRepoCmd)
		if err != nil {
			return err
		}
	}

	fmt.Println("[installBasePackage] executing apt-get install docker-ce")
	const installDockerCmd = "apt-get update && apt-get -y install docker-ce"
	for _, n := range targetNodeList {
		err = processCmdOnNode(n, installDockerCmd)
		if err != nil {
			return err
		}
	}

	fmt.Println("[installBasePackage] executing sysctl docker")
	const restartDockerCmd = "systemctl restart docker"
	for _, n := range targetNodeList {
		err = processCmdOnNode(n, restartDockerCmd)
		if err != nil {
			return err
		}
	}

	fmt.Println("[installBasePackage] executing ntpdate")
	const setNtpCmd = "ntpdate -u time.google.com"
	for _, n := range targetNodeList {
		err = processCmdOnNode(n, setNtpCmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func processCmdOnNode(targetNode node.NodeInterface, command string) error {
	output, err := targetNode.RunSshCmd(common.SshWrapper, command)
	return processExecError(err, output)
}

func copyFile(targetNode node.NodeInterface, role node.Role, srcFile, destFile string) error {
	output, err := targetNode.RunScpCmd(common.ExecWrapper, srcFile, destFile, role)
	return processExecError(err, output)
}

func processExecError(errExec error, output bytes.Buffer) error {
	if errExec != nil {
		if output.Bytes() != nil {
			_, err := output.WriteTo(os.Stderr)
			if err != nil {
				// TODO: combine errExec and err
				return err
			}
		}
		return errExec
	} else {
		_, err := output.WriteTo(os.Stdout)

		return err
	}
}
