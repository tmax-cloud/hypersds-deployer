package main

import (
	provisioner "hypersds-provisioner/cmd/hypersds-provisioner"
	"os"
)

func main() {
	err := provisioner.Run()
	if err != nil {
		os.Exit(1)
	}
}
