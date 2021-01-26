package service

type PlacementInterface interface {
	SetLabel(label string) error
	SetHosts(hosts []string) error
	SetCount(count int) error
	SetHostPattern(hostPattern string) error
	GetLabel() (string, error)
	GetHosts() ([]string, error)
	GetCount() (int, error)
	GetHostPattern() (string, error)
}
type ServiceInterface interface {
	SetServiceType(serviceType string) error
	SetServiceId(serviceID string) error
	SetPlacement(placement PlacementInterface) error
	SetUnmanaged(unmanaged bool) error
	GetServiceType() (string, error)
	GetServiceId() (string, error)
	GetPlacement() (PlacementInterface, error)
	GetUnmanaged() (bool, error)
}
