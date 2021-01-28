package node

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
)

var _ = Describe("Node Initializer Test", func() {
	defer GinkgoRecover()

	var (
		ipAddr, userId, userPw, hostName string
		serviceType                      string
		cephSpec                         hypersdsv1alpha1.CephClusterSpec
	)

	Describe("[NewNodesFromCephCr Test]", func() {
		It("is simple case", func() {
			ipAddr = "1.1.1.1"
			userId = "developer1"
			userPw = "abc123!@#"
			hostName = "node1"
			serviceType = HostServiceType

			cephSpec = hypersdsv1alpha1.CephClusterSpec{
				Nodes: []hypersdsv1alpha1.Node{
					{
						IP:       ipAddr,
						UserID:   userId,
						Password: userPw,
						HostName: hostName,
					},
				},
			}

			nodes, err := NodeInitWrapper.NewNodesFromCephCr(cephSpec)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodes)).To(Equal(1))

			createdUserId, err := nodes[0].GetUserId()
			Expect(err).NotTo(HaveOccurred())
			Expect(createdUserId).To(Equal(userId))

			createdUserPw, err := nodes[0].GetUserId()
			Expect(err).NotTo(HaveOccurred())
			Expect(createdUserPw).To(Equal(userId))

			createdHostSpec, err := nodes[0].GetHostSpec()
			Expect(err).NotTo(HaveOccurred())
			Expect(createdHostSpec.GetServiceType()).To(Equal(serviceType))
			Expect(createdHostSpec.GetHostName()).To(Equal(hostName))
			Expect(createdHostSpec.GetAddr()).To(Equal(ipAddr))
			Expect(createdHostSpec.GetLabels()).To(BeEmpty())
			Expect(createdHostSpec.GetStatus()).To(BeEmpty())
		})
	})
})
