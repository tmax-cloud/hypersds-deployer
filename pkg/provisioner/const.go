package provisioner

// ProvisionerState is current installation phase of provisioner
type ProvisionerState string

const (
	InstallStarted         ProvisionerState = "Initialized"
	BasePkgInstalled       ProvisionerState = "BaseInstalled"
	CephadmPkgInstalled    ProvisionerState = "AdmInstalled"
	CephBootstrapped       ProvisionerState = "Bootstrapped"
	CephBootstrapCommitted ProvisionerState = "Committed"
)

// File name, path, etc
const (
	pathConfigWorkingDir = "/working/config/"
	cephConfToUpdate     = "conf_to_update.conf"
	cephKeyringToUpdate  = "keyring_to_update.keyring"
	cephConfNameFromCr   = "ceph_initial.conf"
	pathConfToUpdate     = pathConfigWorkingDir + cephConfToUpdate
	pathKeyringToUpdate  = pathConfigWorkingDir + cephKeyringToUpdate
	pathConfFromCr       = pathConfigWorkingDir + cephConfNameFromCr
	pathConfFromAdm      = "/etc/ceph/ceph.conf"
	pathKeyringFromAdm   = "/etc/ceph/ceph.client.admin.keyring"
)
