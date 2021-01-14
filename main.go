package main

import (
	"hypersds-provisioner/cmd/hypersds-provisioner"
	"os"
)

func main() {
	err := provisioner.Run()
	if err != nil {
		os.Exit(1)
	}
}
