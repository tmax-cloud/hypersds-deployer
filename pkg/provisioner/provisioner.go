package provisioner

import (
	common "hypersds-provisioner/pkg/common/wrapper"
	config "hypersds-provisioner/pkg/config"
	node "hypersds-provisioner/pkg/node"
	util "hypersds-provisioner/pkg/util"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
)

type ProvisionerInterface interface {
    GetState() (ProvisionerState, error)
    SetState(state ProvisionerState) error
    Run() error
}

type Provisioner struct {
    state ProvisionerState
}

type NewProvisionerStruct struct {
}

func (p *Provisioner) GetState() (ProvisionerState, error) {
    return p.state, nil
}

func (p *Provisioner) SetState(state ProvisionerState) error {
    p.state = state
    return nil
}

var (
	cephClusterSpec hypersdsv1alpha1.CephClusterSpec
	nodeList        []node.NodeInterface
	cephConfig      *config.CephConfig
	deployNode      node.NodeInterface
	err             error
)

func (p *Provisioner) Run() error {
	// Unmarshal yaml file to CephCluster CR
	cephClusterSpec, err = util.UtilWrapper.GetCephClusterSpec(common.IoUtilWrapper, common.YamlWrapper)
	if err != nil {
		return err
	}

    // Get into switch if Ceph is not installed or failure occured during installation
    provisionerState, err := p.GetState()
    if err != nil {
        return err
    }

    switch provisionerState {
    case InstallStarted:
	    // Get nodes info of Ceph hosts
        nodeList, err = node.NewNodeWrapper.NewNodesFromCephCr(cephClusterSpec)
        if err != nil {
            return err
        }

        // Install base package to all nodes
        err = installBasePackage(nodeList)
        if err != nil {
            return err
        }

        // Set provisioner state to BasePkgInstalled
        err = p.SetState(BasePkgInstalled)
        if err != nil {
            return err
        }
        provisionerState, _ = p.GetState()

        fallthrough

    case BasePkgInstalled:
        // Install cephadm package to deploy node
        deployNode = nodeList[0]
        err = installCephadm(deployNode)
        if err != nil {
            return err
        }

        // Set provisioner state to CephadmPkgInstalled
        err = p.SetState(CephadmPkgInstalled)
        if err != nil {
            return err
        }
        provisionerState, _ = p.GetState()

        fallthrough

    case CephadmPkgInstalled:
	    // Extract initial conf file of Ceph
        cephConfig, err = config.ConfigInitWrapper.NewConfigFromCephCr(cephClusterSpec)
        if err != nil {
            return err
        }
        const pathConfFromCr = pathConfigWorkingDir + cephConfNameFromCr
        err = cephConfig.MakeIniFile(common.IoUtilWrapper, pathConfFromCr)
        if err != nil {
            return err
        }

        // Bootstrap ceph on deploy node with cephadm
        err = bootstrapCephadm(deployNode)
        if err != nil {
            return err
        }

        // Set provisioner state to CephBootstrapped
        err = p.SetState(CephBootstrapped)
        if err != nil {
            return err
        }
        provisionerState, _ = p.GetState()

        fallthrough

    case CephBootstrapped:
	    // Copy conf and keyring from deploy node
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

        // Update conf and keyring to ConfigMap and Secret
        err = updateCephClusterToOp()
        if err != nil {
            return err
        }

        // Set provisioner state to CephBootstrapCommitted
        err = p.SetState(CephBootstrapCommitted)
        if err != nil {
            return err
        }
        provisionerState, _ = p.GetState()
    }

    /// TODO: Check diff of host and osd, then apply differences
    return nil
}

func (n *NewProvisionerStruct) NewProvisioner() (ProvisionerInterface, error) {
    var provisioner ProvisionerInterface
    err := provisioner.SetState(InstallStarted)
    return provisioner, err
}
