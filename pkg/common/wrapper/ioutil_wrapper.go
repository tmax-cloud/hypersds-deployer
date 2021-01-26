package wrapper

import (
	"io/ioutil"
)

type IoUtilInterface interface {
	ReadFile(filename string) ([]byte, error)
}

type ioUtilStruct struct {
}

func (i *ioUtilStruct) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

var IoUtilWrapper IoUtilInterface

func init() {
	IoUtilWrapper = &ioUtilStruct{}
}
