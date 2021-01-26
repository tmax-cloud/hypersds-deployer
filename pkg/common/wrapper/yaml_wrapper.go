package wrapper

import (
	"gopkg.in/yaml.v2"
)

type YamlInterface interface {
	Unmarshal(in []byte, out interface{}) (err error)
	Marshal(in interface{}) (out []byte, err error)
}

type yamlStruct struct {
}

func (y *yamlStruct) Unmarshal(in []byte, out interface{}) (err error) {
	return yaml.Unmarshal(in, out)
}
func (y *yamlStruct) Marshal(in interface{}) (out []byte, err error) {
	return yaml.Marshal(in)
}

var YamlWrapper YamlInterface

func init() {
	YamlWrapper = &yamlStruct{}
}
