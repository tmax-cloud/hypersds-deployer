package config

import (
	"hypersds-provisioner/pkg/common/wrapper"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/testing"
)

var _ = Describe("Config Test", func() {
	defer GinkgoRecover()

	var (
		mockCtrl *gomock.Controller
		kubeMock *wrapper.MockKubeInterface
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		kubeMock = wrapper.NewMockKubeInterface(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("[PutConfigToK8s Test]", func() {
		It("should return nil", func() {
			kubeMock.EXPECT().InClusterConfig().Return(nil, nil).AnyTimes()

			kubeMock.EXPECT().NewForConfig(gomock.Any()).DoAndReturn(
				func(c interface{}) (kubernetes.Interface, error) {
					fakeClient := fake.NewSimpleClientset()
					fakeClient.PrependReactor("get", "configmaps", func(action testing.Action) (bool, runtime.Object, error) {
						return true, nil, nil
					})
					fakeClient.PrependReactor("create", "configmap", func(action testing.Action) (bool, runtime.Object, error) {
						return true, nil, nil
					})
					fakeClient.PrependReactor("update", "configmap", func(action testing.Action) (bool, runtime.Object, error) {
						return true, nil, nil
					})
					return fakeClient, nil
				}).AnyTimes()

			testConfig := CephConfig{
				CrConf: map[string]string{
					"test": "test",
				},
			}
			err := testConfig.PutConfigToK8s(kubeMock)
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
