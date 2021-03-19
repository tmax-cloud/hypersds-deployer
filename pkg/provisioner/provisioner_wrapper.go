package provisioner

type NewProvisionerInterface interface {
	NewProvisioner() (ProvisionerInterface, error)
}

var NewProvisionerWrapper NewProvisionerInterface

func init() {
	NewProvisionerWrapper = &NewProvisionerStruct{}
}
