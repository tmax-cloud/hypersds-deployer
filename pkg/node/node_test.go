package node

import (
	"bytes"

	commonWrapper "hypersds-provisioner/pkg/common/wrapper"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
)

var _ = Describe("Node Test", func() {
	defer GinkgoRecover()

	var (
		testingNode                      Node
		userId, userPw, ipAddr, hostName string
		hostSpec                         HostSpec
		cephSpec                         hypersdsv1alpha1.CephClusterSpec
	)

	Describe("Getter/Setter Test", func() {
		It("is simple getter/setter test case", func() {
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
				ServiceType: HostSpecServiceType,
			}
			err = testingNode.SetHostSpec(&hostSpec)
			Expect(err).NotTo(HaveOccurred())

			changedHostSpec, err := testingNode.GetHostSpec()
			Expect(err).NotTo(HaveOccurred())
			Expect(changedHostSpec).To(Equal(&hostSpec))
		})
	})

	Describe("RunSshCmd Test", func() {
		var (
			mockCtrl *gomock.Controller
			m        *commonWrapper.MockExecInterface
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			m = commonWrapper.NewMockExecInterface(mockCtrl)
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		It("is simple test case", func() {
			testCommand := "hello world"
			m.EXPECT().CommandExecute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
				func(resultStdout, resultStderr *bytes.Buffer, ctx, name interface{}, arg ...string) error {
					resultStdout.WriteString("hello world")
					return nil
				}).AnyTimes()

			result, err := testingNode.RunSshCmd(m, testCommand)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.String()).To(Equal(testCommand))
		})
	})

	Describe("[NewNodesFromCephCr Test]", func() {
		It("is simple test case", func() {
			ipAddr = "1.1.1.1"
			userId = "developer1"
			userPw = "abc123!@#"
			hostName = "node1"

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

			nodes, err := NewNodeWrapper.NewNodesFromCephCr(cephSpec)
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
			Expect(createdHostSpec.GetServiceType()).To(Equal(HostSpecServiceType))
			Expect(createdHostSpec.GetHostName()).To(Equal(hostName))
			Expect(createdHostSpec.GetAddr()).To(Equal(ipAddr))
			Expect(createdHostSpec.GetLabels()).To(BeEmpty())
			Expect(createdHostSpec.GetStatus()).To(BeEmpty())
		})
	})
})
