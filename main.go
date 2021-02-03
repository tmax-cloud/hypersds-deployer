package main

import (
	provisioner "hypersds-provisioner/cmd/hypersds-provisioner"
)

func main() {
	err := provisioner.Install()
	if err != nil {
		panic(err)
	}
}
