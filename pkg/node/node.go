package node

type NodeInterface interface {
	SetUserId(userId string) error
	SetUserPw(userPw string) error
	SetHostSpec(hostSpec HostSpecInterface) error

	GetUserId() (string, error)
	GetUserPw() (string, error)
	GetHostSpec() (HostSpecInterface, error)
}

type Node struct {
	userId   string
	userPw   string
	hostSpec HostSpecInterface
}

func (n *Node) SetUserId(userId string) error {
	n.userId = userId
	return nil
}

func (n *Node) SetUserPw(userPw string) error {
	n.userPw = userPw
	return nil
}

func (n *Node) SetHostSpec(hostSpec HostSpecInterface) error {
	n.hostSpec = hostSpec
	return nil
}

func (n *Node) GetUserId() (string, error) {
	return n.userId, nil
}

func (n *Node) GetUserPw() (string, error) {
	return n.userPw, nil
}

func (n *Node) GetHostSpec() (HostSpecInterface, error) {
	return n.hostSpec, nil
}
