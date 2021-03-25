package provisioner

type GetProvisionerInterface interface {
	GetProvisioner() ProvisionerInterface
}

var GetProvisionerWrapper GetProvisionerInterface

func init() {
	GetProvisionerWrapper = &GetProvisionerStruct{}
}
