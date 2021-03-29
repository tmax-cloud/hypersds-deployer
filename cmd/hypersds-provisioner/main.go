package main

import (
	provisioner "hypersds-provisioner/pkg/provisioner"
	"os"
)

func main() {
	provisionerInstance := provisioner.GetProvisionerWrapper.GetProvisioner()

	err := provisionerInstance.Run()
	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
