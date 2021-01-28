package node

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//"fmt"
)

var _ = Describe("Node Test", func() {
	defer GinkgoRecover()

	var (
		testingNode    Node
		userId, userPw string
		hostSpec       HostSpec
	)

	It("is simple getter/setter case", func() {
		// userId getter/setter test
		userId = "shellwedance"
		err := testingNode.SetUserId(userId)
		Expect(err).NotTo(HaveOccurred())

		changedUserId, err := testingNode.GetUserId()
		Expect(err).NotTo(HaveOccurred())
		Expect(changedUserId).To(Equal(userId))

		// userPw getter/setter test
		userPw = "123abc!@#"
		err = testingNode.SetUserPw(userPw)
		Expect(err).NotTo(HaveOccurred())

		changedUserPw, err := testingNode.GetUserPw()
		Expect(err).NotTo(HaveOccurred())
		Expect(changedUserPw).To(Equal(userPw))

		// hostSpec getter/setter test
		hostSpec = HostSpec{
			ServiceType: HostServiceType,
		}
		err = testingNode.SetHostSpec(&hostSpec)
		Expect(err).NotTo(HaveOccurred())

		changedHostSpec, err := testingNode.GetHostSpec()
		Expect(err).NotTo(HaveOccurred())
		Expect(changedHostSpec).To(Equal(&hostSpec))
	})
})
