package provisioner

import (
	common "hypersds-provisioner/pkg/common/wrapper"
	node "hypersds-provisioner/pkg/node"

	"bytes"
	"context"
	"fmt"
	"io"
	"time"
)

func (p *Provisioner) applyHost(yamlWrapper common.YamlInterface, execWrapper common.ExecInterface, ioUtilWrapper common.IoUtilInterface) error {
	// Get host list to apply from CephCluster CR
	cephHostsToApply := []node.HostSpecInterface{}
	cephHostNodesToApply := map[string]node.NodeInterface{}

	nodes, err := p.getNodes()
	if err != nil {
		return err
	}

	for _, node := range nodes {
		host, err := node.GetHostSpec()
		if err != nil {
			return err
		}

		cephHostsToApply = append(cephHostsToApply, host)
		hostName, err := host.GetHostName()
		if err != nil {
			return err
		}
		cephHostNodesToApply[hostName] = node
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()

	var cephDirCreateOutBuf bytes.Buffer
	var cephDirCreateErrBuf bytes.Buffer

	// Create /etc/ceph directory
	const mkdirCmd = "mkdir -p /etc/ceph"
	err = execWrapper.CommandExecute(&cephDirCreateOutBuf, &cephDirCreateErrBuf, ctx, "bash", "-c", mkdirCmd)
	if err != nil {
		fmt.Println("Error: " + cephDirCreateErrBuf.String())
		return err
	}

	// Get ceph conf and keyring from deploy node
	deployNode := nodes[0]
	err = copyFile(deployNode, node.SOURCE, pathConfFromAdm, pathConfFromAdm)
	if err != nil {
		return err
	}

	err = copyFile(deployNode, node.SOURCE, pathKeyringFromAdm, pathKeyringFromAdm)
	if err != nil {
		return err
	}

	// Get current ceph hosts
	var cephadmCurrentHostsOutBuf bytes.Buffer
	var cephadmCurrentHostsErrBuf bytes.Buffer
	const cephHostCheckCmd = "ceph orch host ls yaml"

	fmt.Println("Executing: " + cephHostCheckCmd)
	err = execWrapper.CommandExecute(&cephadmCurrentHostsOutBuf, &cephadmCurrentHostsErrBuf, ctx, "bash", "-c", cephHostCheckCmd)
	if err != nil {
		fmt.Println("Error: " + cephadmCurrentHostsErrBuf.String())
		return err
	}

	fmt.Println("[applyHost] Existing hosts ---")
	fmt.Println(cephadmCurrentHostsOutBuf.String())

	// Extract host specs from ceph orch
	currentHosts := map[string]node.HostSpecInterface{}

	hostReader := bytes.NewReader(cephadmCurrentHostsOutBuf.Bytes())
	decoder := yamlWrapper.NewDecoder(hostReader)

	for {
		host := node.HostSpec{}
		if err = decoder.Decode(&host); err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		hostName, err := host.GetHostName()
		if err != nil {
			return err
		}
		currentHosts[hostName] = &host
	}

	// Compare hosts in CR to Ceph and apply all changes
	for _, hostToApply := range cephHostsToApply {
		hostNameToApply, err := hostToApply.GetHostName()
		if err != nil {
			return err
		}

		if _, exist := currentHosts[hostNameToApply]; exist {
			fmt.Println("Host EXIST!!!" + hostNameToApply)
			continue
		} else {
			// Write hostspec to yml
			hostFileName := fmt.Sprintf("%s%s.yml", pathConfigWorkingDir, hostNameToApply)
			fmt.Println("writing file to ", hostFileName)
			err = hostToApply.MakeYmlFile(yamlWrapper, ioUtilWrapper, hostFileName)
			if err != nil {
				return err
			}

			// generage public key
			var hostAuthGetOutBuf bytes.Buffer
			var hostAuthGetErrBuf bytes.Buffer

			hostAuthGetCmd := fmt.Sprintf("ceph cephadm get-pub-key > %s%s.pub", pathConfigWorkingDir, hostNameToApply)

			fmt.Println("Executing: " + hostAuthGetCmd)
			err = execWrapper.CommandExecute(&hostAuthGetOutBuf, &hostAuthGetErrBuf, ctx, "bash", "-c", hostAuthGetCmd)
			if err != nil {
				fmt.Println("Error: " + hostAuthGetErrBuf.String())
				return err
			}

			// Copy generated key
			var hostAuthApplyOutBuf bytes.Buffer
			var hostAuthApplyErrBuf bytes.Buffer
			nodeId, err := cephHostNodesToApply[hostNameToApply].GetUserId()
			if err != nil {
				return err
			}
			nodeIp, err := hostToApply.GetAddr()
			if err != nil {
				return err
			}
			nodePw, err := cephHostNodesToApply[hostNameToApply].GetUserPw()
			if err != nil {
				return err
			}

			const sshKeyCheckOpt = "-oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null"
			sshPassCmd := fmt.Sprintf("sshpass -f <(printf '%%s\\n' %s)", nodePw)
			hostAuthApplyCmd := fmt.Sprintf("%s ssh-copy-id %s -f -i %s%s.pub %s@%s", sshPassCmd, sshKeyCheckOpt, pathConfigWorkingDir, hostNameToApply, nodeId, nodeIp)

			fmt.Println("Executing: " + hostAuthApplyCmd)
			err = execWrapper.CommandExecute(&hostAuthApplyOutBuf, &hostAuthApplyErrBuf, ctx, "bash", "-c", hostAuthApplyCmd)
			if err != nil {
				fmt.Println("Error: " + hostAuthApplyErrBuf.String())
				return err
			}

			// Apply host
			var hostApplyOutBuf bytes.Buffer
			var hostApplyErrBuf bytes.Buffer
			hostApplyCmd := fmt.Sprintf("ceph orch apply -i %s", hostFileName)

			fmt.Println("Executing: " + hostApplyCmd)
			err = execWrapper.CommandExecute(&hostApplyOutBuf, &hostApplyErrBuf, ctx, "bash", "-c", hostApplyCmd)
			if err != nil {
				fmt.Println("Error: " + hostApplyErrBuf.String())
				return err
			}

			fmt.Println(hostApplyOutBuf.String())
		}
	}

	// Check the result on ceph cluster hosts
	cephadmCurrentHostsOutBuf.Reset()
	cephadmCurrentHostsErrBuf.Reset()

	fmt.Println("Executing: " + cephHostCheckCmd)
	err = execWrapper.CommandExecute(&cephadmCurrentHostsOutBuf, &cephadmCurrentHostsErrBuf, ctx, "bash", "-c", cephHostCheckCmd)
	if err != nil {
		fmt.Println("Error: " + cephadmCurrentHostsErrBuf.String())
		return err
	}

	fmt.Println(cephadmCurrentHostsOutBuf.String())

	return nil
}
