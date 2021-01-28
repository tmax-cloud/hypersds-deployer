package node

import (
	commonWrapper "hypersds-provisioner/pkg/common/wrapper"

	"errors"
)

type HostSpecInterface interface {
	SetServiceType() error
	SetAddr(addr string) error
	SetHostName(hostName string) error
	SetLabels(labels []string) error
	AddLabels(labels ...string) error
	SetStatus(status string) error

	GetServiceType() (string, error)
	GetHostName() (string, error)
	GetAddr() (string, error)
	GetLabels() ([]string, error)
	GetStatus() (string, error)

	MakeYml(wrapper commonWrapper.YamlInterface) ([]byte, error)
}

// variables are required to be importable so that yaml wrapper marshal/unmarshal them
type HostSpec struct {
	ServiceType string   `yaml:"service_type"` // use const var
	Addr        string   `yaml:"addr"`
	HostName    string   `yaml:"hostname"`
	Labels      []string `yaml:"labels,omitempty"`
	Status      string   `yaml:"status,omitempty"`
}

func (hs *HostSpec) SetServiceType() error {
	hs.ServiceType = HostServiceType
	return nil
}

// TODO: add error process (ex. if addr is not IP format)
func (hs *HostSpec) SetAddr(addr string) error {
	hs.Addr = addr
	return nil
}

func (hs *HostSpec) SetHostName(hostName string) error {
	if hostName == "" {
		return errors.New("HostName must not be empty string")
	}

	hs.HostName = hostName
	return nil
}

func (hs *HostSpec) SetLabels(labels []string) error {
	hs.Labels = append([]string{}, labels...)
	return nil
}

func (hs *HostSpec) AddLabels(labels ...string) error {
	hs.Labels = append(hs.Labels, labels...)
	return nil
}

func (hs *HostSpec) SetStatus(status string) error {
	hs.Status = status
	return nil
}

func (hs *HostSpec) GetServiceType() (string, error) {
	return HostServiceType, nil
}

func (hs *HostSpec) GetHostName() (string, error) {
	return hs.HostName, nil
}

func (hs *HostSpec) GetAddr() (string, error) {
	return hs.Addr, nil
}

func (hs *HostSpec) GetLabels() ([]string, error) {
	return hs.Labels, nil
}

func (hs *HostSpec) GetStatus() (string, error) {
	return hs.Status, nil
}

func (hs *HostSpec) MakeYml(wrapper commonWrapper.YamlInterface) ([]byte, error) {
	return wrapper.Marshal(hs)
}
