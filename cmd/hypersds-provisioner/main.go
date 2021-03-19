package main

import (
	provisioner "hypersds-provisioner/pkg/provisioner"
	"os"
)

func main() {
	provisioner, err := provisioner.NewProvisionerWrapper.NewProvisioner()
	if err != nil {
		panic(err)
	}

    err = provisioner.Run()
    if err != nil {
        panic(err)
    }

	os.Exit(0)
}
