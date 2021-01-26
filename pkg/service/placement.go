package service

type Placement struct {
	Label       string   `yaml:"label,omitempty"`
	Hosts       []string `yaml:"hosts,omitempty"`
	Count       int      `yaml:"count,omitempty"`
	HostPattern string   `yaml:"host_pattern,omitempty"`
}

func (p *Placement) SetLabel(label string) error {
	p.Label = label
	return nil
}
func (p *Placement) SetHosts(hosts []string) error {
	p.Hosts = hosts
	return nil
}
func (p *Placement) SetCount(count int) error {
	p.Count = count
	return nil
}
func (p *Placement) SetHostPattern(hostPattern string) error {
	p.HostPattern = hostPattern
	return nil
}

func (p Placement) GetLabel() (string, error) {
	return p.Label, nil
}
func (p Placement) GetHosts() ([]string, error) {
	return p.Hosts, nil
}
func (p Placement) GetCount() (int, error) {
	return p.Count, nil
}
func (p Placement) GetHostPattern() (string, error) {
	return p.HostPattern, nil
}
