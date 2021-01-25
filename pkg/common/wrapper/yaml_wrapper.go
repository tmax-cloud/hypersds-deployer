package wrapper

import (
    "gopkg.in/yaml.v2"
)

type YamlInterface interface {
    Unmarshal(in []byte, out interface{}) (err error)
}

type yamlStruct struct {
}

func (y *yamlStruct) Unmarshal(in []byte, out interface{}) (err error) {
    return yaml.Unmarshal(in, out)
}

var YamlWrapper YamlInterface

func init() {
    YamlWrapper = &yamlStruct{}
}