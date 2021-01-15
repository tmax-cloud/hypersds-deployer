package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/user"
)

var _ = Describe("SSH Command Test", func() {
	defer GinkgoRecover()
	Describe("[RunSSHCmd Test]", func() {
		It("should return same username", func() {
			currentUser, err := user.Current()
			testCommand := []string{"printf", "$USER"}

			result, err := RunSSHCmd(currentUser.Username, "localhost", testCommand...)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.String()).To(Equal(currentUser.Username))
		})
	})

})
