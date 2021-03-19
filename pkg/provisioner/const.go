package provisioner

// ProvisionerState is current installation phase of provisioner
type ProvisionerState string

const (
    InstallStarted ProvisionerState = "Initialized"
    BasePkgInstalled ProvisionerState = "BaseInstalled"
    CephadmPkgInstalled ProvisionerState = "AdmInstalled"
    CephBootstrapped ProvisionerState = "Bootstrapped"
    CephBootstrapCommitted ProvisionerState = "Committed"
)

// File name, path, etc
const (
	pathConfigWorkingDir = "/working/config/"
	cephConfNameFromCr   = "ceph_initial.conf"
	pathConfFromAdm      = "/etc/ceph/ceph.conf"
	pathKeyringFromAdm   = "/etc/ceph/ceph.client.admin.keyring"
	cephConfToUpdate     = "conf_to_update.conf"
	cephKeyringToUpdate  = "keyring_to_update.keyring"
)
