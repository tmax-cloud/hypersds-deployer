package main

import (
	provisioner "hypersds-provisioner/cmd/hypersds-provisioner"
	"os"
)

func main() {
	err := provisioner.Install()
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}
