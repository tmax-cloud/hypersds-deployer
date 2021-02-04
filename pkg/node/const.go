package node

import (
	"time"
)

const (
	HostSpecServiceType = "host"
	SshCmdTimeout       = 20 * time.Minute
)

type Role int

const (
	SRC Role = iota
	DEST
)
