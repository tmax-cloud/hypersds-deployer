package util

import (
	"fmt"
	"hypersds-provisioner/pkg/common/wrapper"
	"log"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
)

var _ = Describe("CephSpec", func() {
	defer GinkgoRecover()
	var (
		mockCtrl *gomock.Controller
		i        *wrapper.MockIoUtilInterface
	)
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		i = wrapper.NewMockIoUtilInterface(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("[Run getCephClusterSpec]", func() {
		It("just test", func() {
			log.Printf("[TEST] getCephClusterSpec")
			i.EXPECT().ReadFile(gomock.Any()).DoAndReturn(
				func(filename string) ([]byte, error) {
					source := []byte{109, 111, 110, 58, 10, 32, 32, 99, 111, 117, 110, 116, 58, 32, 49, 10, 110, 111, 100, 101, 115, 58, 10, 32, 32, 45, 32, 104, 111, 115, 116, 110, 97, 109, 101, 58, 32, 34, 117, 110, 105, 116, 84, 101, 115, 116, 34, 10, 99, 111, 110, 102, 105, 103, 58, 32, 10, 32, 32, 115, 105, 109, 112, 108, 101, 58, 32, 34, 105, 115, 32, 103, 111, 111, 100, 34, 10}
					return source, nil
				}).AnyTimes()
			var t hypersdsv1alpha1.CephClusterSpec
			t, _ = getCephClusterSpec(i, wrapper.YamlWrapper)
			Expect(fmt.Sprint(t)).To(Equal("{{1} [] [{   unitTest}] map[simple:is good]}"))
		})
	})

})
