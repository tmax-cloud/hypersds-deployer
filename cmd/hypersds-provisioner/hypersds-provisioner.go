package provisioner

import (
	"bytes"
	"fmt"
	"os"

	common "hypersds-provisioner/pkg/common/wrapper"
	config "hypersds-provisioner/pkg/config"
	node "hypersds-provisioner/pkg/node"
	util "hypersds-provisioner/pkg/util"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
)

const (
	pathConfigWorkingDir = "/"
	cephConfNameFromCr   = "ceph_initial.conf"
	pathConfFromAdm      = "/etc/ceph/ceph.conf"
	pathKeyringFromAdm   = "/etc/ceph/ceph.client.admin.keyring"
	cephConfToUpdate     = "conf_to_update.conf"
	cephKeyringToUpdate  = "keyring_to_update.keyring"
)

var (
	cephClusterSpec hypersdsv1alpha1.CephClusterSpec
	nodeList        []node.NodeInterface
	cephConfig      *config.CephConfig
	deployNode      node.NodeInterface
	err             error
)

func Install() error {
	// 1. Unmarshal yaml file to CephCluster CR
	cephClusterSpec, err = util.UtilWrapper.GetCephClusterSpec(common.IoUtilWrapper, common.YamlWrapper)
	if err != nil {
		return err
	}

	// 2. Get nodes info of Ceph hosts
	nodeList, err = node.NewNodeWrapper.NewNodesFromCephCr(cephClusterSpec)
	if err != nil {
		return err
	}

	// 3. Extract initial conf file of Ceph
	cephConfig, err = config.ConfigInitWrapper.NewConfigFromCephCr(cephClusterSpec)
	if err != nil {
		return err
	}

	const pathConfFromCr = pathConfigWorkingDir + cephConfNameFromCr
	err = cephConfig.MakeIniFile(common.IoUtilWrapper, pathConfFromCr)
	if err != nil {
		return err
	}

	// 4. Install required packages to all nodes
	err = installBasePackage(nodeList)
	if err != nil {
		return err
	}

	deployNode = nodeList[0]

	// 5. Install cephadm to deploy node
	err = installCephadm(deployNode)
	if err != nil {
		return err
	}

	// 6. Bootstrap cephadm on deploy node
	err = bootstrapCephadm(deployNode)
	if err != nil {
		return err
	}

	// 7. Copy conf and keyring from deploy node
	const pathConfToUpdate = pathConfigWorkingDir + cephConfToUpdate
	err = copyFile(deployNode, node.SOURCE, pathConfFromAdm, pathConfToUpdate)
	if err != nil {
		return err
	}
	const pathKeyringToUpdate = pathConfigWorkingDir + cephKeyringToUpdate
	err = copyFile(deployNode, node.SOURCE, pathKeyringFromAdm, pathKeyringToUpdate)
	if err != nil {
		return err
	}

	// 8. Update conf and keyring to ConfigMap and Secret
	err = updateCephClusterToOp()

	return err
}

func updateCephClusterToOp() error {
	fmt.Println("----------------Start to update conf and keyring to operator---------------")

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
	fmt.Println("----------------Start to bootstrap ceph---------------")

	fmt.Println("[bootstrapCephadm] copying conf file")
	const pathConfFromCr = pathConfigWorkingDir + cephConfNameFromCr
	err = copyFile(targetNode, node.DESTINATION, pathConfFromCr, cephConfNameFromCr)
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
	admBootstrapCmd := fmt.Sprintf("cephadm bootstrap --mon-ip %s --config %s", monIp, pathConfFromCr)

	err = processCmdOnNode(targetNode, admBootstrapCmd)

	return err
}

func installCephadm(targetNode node.NodeInterface) error {
	fmt.Println("----------------Start to install cephadm---------------")

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
	fmt.Println("----------------Start to install base package---------------")

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
	output, err := targetNode.RunSshCmd(common.ExecWrapper, command)
	return processExecError(err, output)
}

func copyFile(targetNode node.NodeInterface, role node.Role, srcFile, destFile string) error {
	output, err := targetNode.RunScpCmd(common.ExecWrapper, srcFile, destFile, role)
	return processExecError(err, output)
}

func processExecError(errExec error, output bytes.Buffer) error {
	if errExec != nil {
		// case that RunScpCmd failed SSH and successed to return the stderr result
		if output.Bytes() != nil {
            //////////
            bufString := bytes.NewBufferString("\n------ stderr start--------\n")
            _, err := bufString.Write(output.Bytes())
            if err != nil {
                fmt.Println("added error, stderr shit1")
                return err
            } else {
                fmt.Println("added error, stderr shit2")
                return errExec
            }

            _, err = bufString.WriteString("\n------ stderr done--------")
            if err != nil {
                fmt.Println("added error, shit3")
                return err
            } else {
                fmt.Println("added error, shit4")
                return errExec
            }
            //////////

			//_, err := output.WriteTo(os.Stderr)
			_, err = bufString.WriteTo(os.Stderr)
			if err != nil {
                fmt.Println("added error, shit5")
				return err
			} else {
                fmt.Println("added error, shit6")
				return errExec
			}

			// case that RunScpCmd failed before calling SSH
		} else {
			return errExec
		}

		// case that RunScpCmd succeeded SSH
	} else {
        //////////
        bufString := bytes.NewBufferString("\n------ stdout start--------\n")
        numWrt, err := bufString.Write(output.Bytes())
        fmt.Println("written bytes: ", numWrt)
        if err != nil {
            fmt.Println("added error, shit3")
            return err
        }
        _, err = bufString.WriteString("\n------ stdout done--------\n")
        if err != nil {
            fmt.Println("added error, shit4")
            return err
        }
        //////////

		//_, err := output.WriteTo(os.Stdout)
		_, err = bufString.WriteTo(os.Stdout)
		if err != nil {
            fmt.Println("added error, shit5")
			return err
		} else {
            fmt.Println("added error, shit6")
			return nil
		}
	}
}
