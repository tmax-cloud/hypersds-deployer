package node

import (
	"fmt"
	commonWrapper "hypersds-provisioner/pkg/common/wrapper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HostSpec Test", func() {
	defer GinkgoRecover()

	var (
		hostSpec                 HostSpec
		ipAddr, hostName, status string
		labels, labelsToAdd      []string
	)

	BeforeEach(func() {
		ipAddr = "1.1.1.1"
		hostName = "node1"
		status = "CEPHADM_STRAY_HOST"
		labels = []string{"example1", "example2"}
		labelsToAdd = []string{"example3", "example4"}

		hostSpec = HostSpec{
			ServiceType: HostServiceType,
			Addr:        ipAddr,
			HostName:    hostName,
			Labels:      labels,
			Status:      status,
		}
	})

	// TODO: replace to ginkgo table extension
	It("is simple getter/setter case", func() {
		// ServiceType getter/setter test
		err := hostSpec.SetServiceType()
		Expect(err).NotTo(HaveOccurred())

		createdServiceType, err := hostSpec.GetServiceType()
		Expect(err).NotTo(HaveOccurred())
		Expect(createdServiceType).To(Equal(HostServiceType))

		// Addr getter/setter test
		newIpAddr := "1.1.1.2"
		err = hostSpec.SetAddr(newIpAddr)
		Expect(err).NotTo(HaveOccurred())

		changedAddr, err := hostSpec.GetAddr()
		Expect(err).NotTo(HaveOccurred())
		Expect(changedAddr).To(Equal(newIpAddr))

		// HostName getter/setter test
		newHostName := "node2"
		err = hostSpec.SetHostName(newHostName)
		Expect(err).NotTo(HaveOccurred())

		changedHostName, err := hostSpec.GetHostName()
		Expect(err).NotTo(HaveOccurred())
		Expect(changedHostName).To(Equal(newHostName))

		// Labels getter/setter test
		newLabels := []string{"exampleA", "exampleB"}
		err = hostSpec.SetLabels(newLabels)
		Expect(err).NotTo(HaveOccurred())

		changedLabels, err := hostSpec.GetLabels()
		Expect(err).NotTo(HaveOccurred())
		Expect(changedLabels).To(Equal(newLabels))

		// Labels adder test
		err = hostSpec.AddLabels(labelsToAdd...)
		Expect(err).NotTo(HaveOccurred())

		allLabels, err := hostSpec.GetLabels()
		Expect(err).NotTo(HaveOccurred())
		newLabels = append(newLabels, labelsToAdd...)
		Expect(allLabels).To(Equal(newLabels))

		// HostName getter/setter test
		newStatus := "CEPHADM_HOST_CHECK_FAILED"
		err = hostSpec.SetStatus(newStatus)
		Expect(err).NotTo(HaveOccurred())

		changedStatus, err := hostSpec.GetStatus()
		Expect(err).NotTo(HaveOccurred())
		Expect(changedStatus).To(Equal(newStatus))
	})

	It("is simple MakeYml case", func() {
		expectedYmlString := fmt.Sprintf("service_type: %s\naddr: %s\nhostname: %s\nlabels:\n- %s\n- %s\nstatus: %s\n",
			HostServiceType, ipAddr, hostName, labels[0], labels[1], status)

		createdBytes, err := hostSpec.MakeYml(commonWrapper.YamlWrapper)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(createdBytes)).To(Equal(expectedYmlString))
	})
})
