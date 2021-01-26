package osd

import (
	"fmt"
	"hypersds-provisioner/pkg/common/wrapper"
	"hypersds-provisioner/pkg/service"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
)

type Osd struct {
	DataDevices     DeviceInterface `yaml:"data_devices,omitempty"`
	DbDevices       DeviceInterface `yaml:"db_devices,omitempty"`
	WalDevices      DeviceInterface `yaml:"wal_devices,omitempty"`
	JournalDevices  DeviceInterface `yaml:"journal_devices,omitempty"`
	DataDirectories []string        `yaml:"data_directories,omitempty"`
	OsdsPerDevice   int             `yaml:"osds_per_device,omitempty"`
	Objectstore     string          `yaml:"objectstore,omitempty"`
	Encrypted       bool            `yaml:"encrypted,omitempty"`

	service service.ServiceInterface `yaml:"-"`
	/*
	    				 db_slots=None,  # type: Optional[int]
	                    wal_slots=None,  # type: Optional[int]
	                    osd_id_claims=None,  # type: Optional[Dict[str, List[str]]]
	                    block_db_size=None,  # type: Union[int, str, None]
	                    block_wal_size=None,  # type: Union[int, str, None]
	                    journal_size=None,  # type: Union[int, str, None]
	                    service_type=None,  # type: Optional[str]
	                    unmanaged=False,  # type: bool
	                    filter_logic='AND',  # type: str
	   				 preview_only=False,  # type: bool
	*/
}

func (o *Osd) SetService(s service.ServiceInterface) error {
	o.service = s
	return nil
}

func (o *Osd) SetDataDevices(dataDevices DeviceInterface) error {
	o.DataDevices = dataDevices
	return nil
}

func (o Osd) GetService() (service.ServiceInterface, error) {
	return o.service, nil
}

func (o Osd) GetDataDevices() (DeviceInterface, error) {
	return o.DataDevices, nil
}

func NewOsdFromCephCr(osdSpec hypersdsv1alpha1.CephClusterOsdSpec) *Osd {
	var hosts []string
	osd := Osd{}
	dataDevices := Device{}
	s := service.Service{}
	placement := service.Placement{}

	// set Placement, Service
	hosts = append(hosts, osdSpec.HostName)
	_ = placement.SetHosts(hosts)

	_ = s.SetPlacement(&placement)
	_ = s.SetServiceType("osd")
	_ = s.SetServiceId("osd_" + osdSpec.HostName)

	osd.SetService(&s)

	// set device
	_ = dataDevices.SetPaths(osdSpec.Devices)
	_ = osd.SetDataDevices(&dataDevices)

	return &osd
}

func (osd *Osd) MakeYml(yaml wrapper.YamlInterface) ([]byte, error) {
	s, _ := osd.GetService()
	serviceYaml, _ := yaml.Marshal(s)
	oYaml, _ := yaml.Marshal(osd)
	osdYaml := append(serviceYaml, oYaml...)
	return osdYaml, nil
}

func TestOsd() {
	osdSpec := hypersdsv1alpha1.CephClusterOsdSpec{
		HostName: "test", Devices: []string{"ha"},
	}
	osd := NewOsdFromCephCr(osdSpec)
	tempYaml, _ := osd.MakeYml(wrapper.YamlWrapper)
	fmt.Print(string(tempYaml))
}

/*

func GetOsdsFromCephCR(cephSpec hypersdsv1alpha1.CephClusterSpec) []*Osd {
	var osdList []*Osd
	for i := 0; i < len(cephSpec.Osd); i++ {
		osd := NewOsdFromCephCR(cephSpec.Osd[i])
		osdList = append(osdList, osd)
	}
	return osdList
}
*/
