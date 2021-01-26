package node

import(
    wrapper "hypersds-provisioner/pkg/common/wrapper"
)

type HostSpec struct{
    ServiceType string `yaml:"service_type"`    // use const var
    Addr string `yaml:"addr"`
    HostName string `yaml:"hostname"`
    Labels []string `yaml:"labels,omitempty"`
    Status string `yaml:status,omitempty`
}

func (hs *HostSpec) SetServiceType() error {
    hs.ServiceType = HostServiceType
    return nil
}

func (hs *HostSpec) SetAddr(addr string) error {
    hs.Addr = addr
    return nil
}

func (hs *HostSpec) SetHostName(hostName string) error {
    hs.HostName = hostName
    return nil
}

func (hs *HostSpec) SetLabels(labels []string) error {
    hs.Labels = append([]string{}, labels...)
    return nil
}

func (hs *HostSpec) AddLabels(label ...string) error {
    hs.Labels = append(hs.Labels, label...)
    return nil
}

func (hs *HostSpec) SetStatus(status string) error {
    hs.Status = status
    return nil
}

func (hs *HostSpec) GetServiceType() (string, error) {
    return HostServiceType, nil
}

// TODO: add error process (ex. if hostname is "")
func (hs *HostSpec) GetHostName() (string, error) {
    return hs.HostName, nil
}

// TODO: add error process (ex. if addr is not IP format)
func (hs *HostSpec) GetAddr() (string, error) {
    return hs.Addr, nil
}

func (hs *HostSpec) GetLabels() ([]string, error) {
    return hs.Labels, nil
}

func (hs *HostSpec) GetStatus() (string, error) {
    return hs.Status, nil
}

func (hs *HostSpec) MakeYml(wr wrapper.YamlInterface) ([]byte, error) {
    return wr.Marshal(hs)
}
