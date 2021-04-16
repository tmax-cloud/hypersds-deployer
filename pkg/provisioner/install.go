package provisioner

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	common "hypersds-provisioner/pkg/common/wrapper"
	node "hypersds-provisioner/pkg/node"
	"hypersds-provisioner/pkg/osd"
	util "hypersds-provisioner/pkg/util"
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
	admBootstrapCmd := fmt.Sprintf("cephadm --image %s bootstrap --mon-ip %s --config %s",
		cephImageName, monIp, pathConf)
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

	fmt.Println("[installCephadm] executing curl cephadm")
	curlCephadmCmd := fmt.Sprintf("curl --silent --remote-name --location https://github.com/ceph/ceph/raw/v%s/src/cephadm/cephadm", cephVersion)
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
	admAddRepoCmd := fmt.Sprintf("./cephadm add-repo --version %s", cephVersion)
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
	const mkdirCmd = "mkdir -p /etc/ceph"
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

	fmt.Println("[installBasePackage] executing ntpdate")
	const setNtpCmd = "ntpdate -u time.google.com"
	for _, n := range targetNodeList {
		err = processCmdOnNode(n, setNtpCmd)
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

	return nil
}

func (p *Provisioner) applyOsd() error {
	fmt.Println("[applyOsd] get osds from CephOrch")

	cmd := []string{"orch", "ls", "--service_type", "osd", "--export", "--refresh"}
	output, err := util.RunCephCmd(common.ExecWrapper, cmd...)
	if err != nil {
		return processExecError(err, output)
	}

	var osdsFromOrch []*osd.Osd
	if !strings.Contains("No services reported", output.String()) {
		rawOsdsFromOrch := output.Bytes()

		osdsFromOrch, err = osd.OsdWrapper.NewOsdsFromCephOrch(common.YamlWrapper, rawOsdsFromOrch)
		if err != nil {
			return err
		}
	}
	osdsFromCephCr, err := osd.OsdWrapper.NewOsdsFromCephCr(p.cephCluster)
	if err != nil {
		return err
	}

	var osdMap map[string]*osd.Osd
	var removeOsdMap map[string]bool

	osdMap = make(map[string]*osd.Osd)
	removeOsdMap = make(map[string]bool)

	for _, osdOrch := range osdsFromOrch {
		osdService, err := osdOrch.GetService()
		if err != nil {
			return err
		}
		osdServiceId, err := osdService.GetServiceId()
		if err != nil {
			return err
		}
		osdMap[osdServiceId] = osdOrch
		removeOsdMap[osdServiceId] = true
	}

	fmt.Println("[applyOsd] compare osds between CephCR and CephOrch")

	for _, osdCephCr := range osdsFromCephCr {
		osdService, err := osdCephCr.GetService()
		if err != nil {
			return err
		}
		osdServiceId, err := osdService.GetServiceId()
		if err != nil {
			return err
		}
		osdOrch, exist := osdMap[osdServiceId]
		if exist {
			addDeviceList, removeDeviceList, err := osdOrch.CompareDataDevices(osdCephCr)
			if err != nil {
				return err
			}
			removeOsdMap[osdServiceId] = false
			//todo remove disk ....
			fmt.Printf("[applyOsd] osd service: %s, add: %+q, remove: %+q\n", osdServiceId, addDeviceList, removeDeviceList)
		}

		fmt.Println("[applyOsd] make osd yaml")

		osdFileName := osdServiceId + ".yaml"
		err = osdCephCr.MakeYmlFile(common.YamlWrapper, common.IoUtilWrapper, osdFileName)
		if err != nil {
			return err
		}

		fmt.Printf("[applyOsd] apply osd service: %s\n", osdServiceId)

		applyCmd := []string{"orch", "apply", "-i", osdFileName}
		output, err = util.RunCephCmd(common.ExecWrapper, applyCmd...)
		if err != nil {
			return processExecError(err, output)
		}
	}
	for osdServiceId, value := range removeOsdMap {
		if value {
			osdServiceName := "osd." + osdServiceId

			fmt.Printf("[applyOsd] remove osd service: %s\n", osdServiceName)

			removeCmd := []string{"orch", "rm", osdServiceName}
			output, err = util.RunCephCmd(common.ExecWrapper, removeCmd...)
			if err != nil {
				return processExecError(err, output)
			}
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
