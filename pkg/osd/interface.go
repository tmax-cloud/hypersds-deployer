package osd

import (
	"hypersds-provisioner/pkg/common/wrapper"
	"hypersds-provisioner/pkg/service"
)

type DeviceInterface interface {
	SetPaths(paths []string) error
	SetModel(model string) error
	SetSize(size string) error
	SetRotational(rotational bool) error
	SetLimit(limit int) error
	SetVendor(vendor string) error
	SetAll(all bool) error

	GetPaths() ([]string, error)
	GetModel() (string, error)
	GetSize() (string, error)
	GetRotational() (bool, error)
	GetLimit() (int, error)
	GetVendor() (string, error)
	GetAll() (bool, error)
}

type OsdInterface interface {
	SetService(s service.ServiceInterface) error
	SetDataDevices(dataDevices DeviceInterface) error
	GetService() (service.ServiceInterface, error)
	GetDataDevices() (DeviceInterface, error)

	MakeYml(yaml wrapper.YamlInterface) ([]byte, error)
}
