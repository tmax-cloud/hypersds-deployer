package provisioner

import (
	common "hypersds-provisioner/pkg/common/wrapper"
	config "hypersds-provisioner/pkg/config"
	node "hypersds-provisioner/pkg/node"
	util "hypersds-provisioner/pkg/util"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    errors "k8s.io/apimachinery/pkg/api/errors"

    "fmt"
    "context"
)

type ProvisionerInterface interface {
    GetState() (ProvisionerState, error)
    GetNodes() ([]node.NodeInterface, error)
    Run() error

    setState(state ProvisionerState) error
    setCephCluster(cephCluster hypersdsv1alpha1.CephClusterSpec) error
    identifyProvisionerState() (ProvisionerState, error)
}

type Provisioner struct {
	cephCluster hypersdsv1alpha1.CephClusterSpec
    state ProvisionerState
}
// Provisioner singleton instance
var provisionerInstance ProvisionerInterface

var (
    cephClusterSpec hypersdsv1alpha1.CephClusterSpec
	nodeList        []node.NodeInterface
	cephConfig      *config.CephConfig
	deployNode      node.NodeInterface
	err             error
)

func (p *Provisioner) GetState() (ProvisionerState, error) {
    return p.state, nil
}

func (p *Provisioner) GetNodes() ([]node.NodeInterface, error) {
    return node.NewNodeWrapper.NewNodesFromCephCr(p.cephCluster)
}

func (p *Provisioner) Run() error {
    // Decide deploying node (currently, first node is deploying node)
    nodeList, err = p.GetNodes()
    if err != nil {
        return err
    }
    deployNode = nodeList[0]

    // Get into switch if Ceph is not installed or failure occured during installation
    provisionerState, err := p.GetState()
    if err != nil {
        return err
    }

    switch provisionerState {
    case InstallStarted:
        // Install base package to all nodes
        err = installBasePackage(nodeList)
        if err != nil {
            return err
        }

        // Set provisioner state to BasePkgInstalled
        err = p.setState(BasePkgInstalled)
        if err != nil {
            return err
        }
        provisionerState, _ = p.GetState()

        fallthrough

    case BasePkgInstalled:
        // Install cephadm package to deploying node
        err = installCephadm(deployNode)
        if err != nil {
            return err
        }

        // Set provisioner state to CephadmPkgInstalled
        err = p.setState(CephadmPkgInstalled)
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
        err = p.setState(CephBootstrapped)
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
        err = p.setState(CephBootstrapCommitted)
        if err != nil {
            return err
        }
        provisionerState, _ = p.GetState()
    }

    /// TODO: Check diff of host and osd, then apply differences
    return nil
}

func (p *Provisioner) setState(state ProvisionerState) error {
    p.state = state
    return nil
}

func (p *Provisioner) setCephCluster(cephCluster hypersdsv1alpha1.CephClusterSpec) error {
    p.cephCluster = cephCluster
    return nil
}

func (p *Provisioner) identifyProvisionerState() (ProvisionerState, error) {
    var state ProvisionerState
    // Decide deploying node (currently, first node is deploying node)
    nodes, err := p.GetNodes()
    if err != nil {
        return "", err
    }
    deployNode := nodes[0]

    const checkDockerActive = "dpkg --list | grep docker-ce"
	output, err := deployNode.RunSshCmd(common.SshWrapper, checkDockerActive)
    if err != nil {
        fmt.Println("fuck is it error?")
        fmt.Println(err)
        return "", err
    }

    fmt.Println(output.String())

    // Check Ceph bootstrap is committed
    committed, err := checkKubeObjectUpdated(common.KubeWrapper)
    if err != nil {
        return "", err
    } else if committed {
        //return CephBootstrapCommitted, nil
        state = CephBootstrapCommitted
    }

    return state, nil
}

// TODO: replace config const to inputs (e.g. K8sConfigMap, etc)
func checkKubeObjectUpdated(kubeWrapper common.KubeInterface) (bool, error) {
    kubeConfig, err := kubeWrapper.InClusterConfig()
    if err != nil {
        return false, err
    }

    clientSet, err := kubeWrapper.NewForConfig(kubeConfig)
    if err != nil {
        return false, err
    }

    // Check ceph.conf is updated to ConfigMap
    configMap, err := clientSet.CoreV1().ConfigMaps(config.K8sNamespace).Get(context.TODO(), config.K8sConfigMap, metav1.GetOptions{})
    if err != nil {
        // configmap must exist
        if errors.IsNotFound(err) {
            // TODO: replace stdout to log out
            fmt.Println("ConfigMap must exist")
            return false, nil
        } else {
            return false, err
        }
    }

    // bootstrap commit has not occurred
    if configMap.Data == nil {
        return false, nil
    }

    // Check client.admin.keyring is updated to Secret
	secret, err := clientSet.CoreV1().Secrets(config.K8sNamespace).Get(context.TODO(), config.K8sSecret, metav1.GetOptions{})
	if err != nil {
        if errors.IsNotFound(err) {
            // TODO: replace stdout to log out
            fmt.Println("Secret must exist")
            return false, nil
        } else {
            return false, err
        }
	}

    if secret.Data == nil {
        return false, nil
    }

    return true, nil
}

type GetProvisionerStruct struct {
}

func (n *GetProvisionerStruct) GetProvisioner() ProvisionerInterface {
    return provisionerInstance
}

func init() {
    // For singleton pattern
    provisionerInstance = &Provisioner{}

	// Unmarshal yaml file to CephCluster CR
	cephClusterSpec, err = util.UtilWrapper.GetCephClusterSpec(common.IoUtilWrapper, common.YamlWrapper)
	if err != nil {
        // TODO: replace stdout to log out
        fmt.Println("[Provisioner] CephClusterSpec Parsing Error")
		panic(err)
	}

    // setCephCluster is only called once, on init, and no one is allowed to modify CephCluster
    err := provisionerInstance.setCephCluster(cephClusterSpec)
    if err != nil {
        // TODO: replace stdout to log out
        fmt.Println("[Provisioner] setCephCluster Error")
		panic(err)
	}

    // identifyProvisionerState is only called once, on init, and no one is allowed to modify ProvisionerState
    provisionerState, err := provisionerInstance.identifyProvisionerState()
    if err != nil {
        // TODO: replace stdout to log out
        fmt.Println("[Provisioner] identifyProvisionerState Error")
		panic(err)
	}

    err = provisionerInstance.setState(provisionerState)
    if err != nil {
        // TODO: replace stdout to log out
        fmt.Println("[Provisioner] setState Error")
		panic(err)
	}
}
