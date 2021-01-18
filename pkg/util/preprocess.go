package util

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type MON struct {
	Count int `yaml:"count"`
}

type OSDService struct {
	HostName string   `yaml:"hostname"`
	Devices  []string `yaml:"devices"`
}

type CephNode struct {
	Ip       string `yaml:"ip"`
	UserID   string `yaml:"userid"`
	Passwd   string `yaml:"password"`
	HostName string `yaml:"hostname"`
}

type CephCluster struct {
	Mon         MON               `yaml:"mon"`
	OSDServices []OSDService      `yaml:"osd"`
	CephNodes   []CephNode        `yaml:"nodes"`
	Config      map[string]string `yaml:"config"`
}

func ParseYaml() (CephCluster, error) {
	filename := "test.yaml"
	var cephcluster = CephCluster{}
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return CephCluster{}, err
	}
	err = yaml.Unmarshal(source, &cephcluster)
	if err != nil {
		return CephCluster{}, err
	}

	return cephcluster, nil
}

func GenerateConfFile(cc CephCluster) error {
	f1, err := os.Create("ceph.conf")
	if err != nil {
		return err
	}
	defer f1.Close()
	fmt.Fprintf(f1, "	[global]\n")
	for k, v := range cc.Config {
		fmt.Fprintf(f1, "%s = %s\n", k, v)
	}
	return nil
}
