package service

// https://github.com/ceph/ceph/blob/master/src/python-common/ceph/deployment/service_spec.py

type Service struct {
	ServiceType string    `yaml:"service_type"`
	ServiceId   string    `yaml:"service_id"`
	Placement   Placement `yaml:"placement,omitempty"`
	Unmanaged   bool      `yaml:"unmanaged,omitempty"`
	//previewed_only : maybe not used
}

func (s *Service) SetServiceType(serviceType string) error {
	s.ServiceType = serviceType
	return nil
}
func (s *Service) SetServiceId(serviceId string) error {
	s.ServiceId = serviceId
	return nil
}
func (s *Service) SetPlacement(placement Placement) error {
	s.Placement = placement
	return nil
}
func (s *Service) SetUnmanaged(unmanaged bool) error {
	s.Unmanaged = unmanaged
	return nil
}

func (s Service) GetServiceType() (string, error) {
	return s.ServiceType, nil
}
func (s Service) GetServiceId() (string, error) {
	return s.ServiceId, nil
}
func (s Service) GetPlacement() (Placement, error) {
	return s.Placement, nil
}
func (s Service) GetUnmanaged() (bool, error) {
	return s.Unmanaged, nil
}
