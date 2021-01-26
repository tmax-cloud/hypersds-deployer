package node

import(
    wrapper "hypersds-provisioner/pkg/common/wrapper"
)

type HostSpecInterface interface {
    SetServiceType() error
    SetHostName(hostName string) error
    SetAddr(addr string) error
    SetLabels(labels []string) error
    SetStatus(status string) error

    GetServiceType() (string, error)
    GetHostName() (string, error)
    GetAddr() (string, error)
    GetLabels() ([]string, error)
    GetStatus() (string, error)

    MakeYml(wr wrapper.YamlInterface) ([]byte, error)
}

type NodeInterface interface {
    SetUserId(userId string) error
    SetUserPw(userPw string) error
    SetHostSpec(hostSpec HostSpecInterface) error

    GetUserId() (string, error)
    GetUserPw() (string, error)
    GetHostSpec() (HostSpecInterface, error)
}
