package config

import (
	"hypersds-provisioner/pkg/common/wrapper"
	"os"

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
		ioMock   *wrapper.MockIoUtilInterface
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		kubeMock = wrapper.NewMockKubeInterface(mockCtrl)
		ioMock = wrapper.NewMockIoUtilInterface(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("[ConfigFromAdm Test]", func() {
		It("Parse ceph.conf to AdmConfig", func() {
			ioMock.EXPECT().ReadFile(gomock.Any()).DoAndReturn(
				func(filename string) ([]byte, error) {
					conf := []byte("[global]\n\tfsid = b29fd\n\tmon_host = [0.0.0.0]")
					return conf, nil
				}).AnyTimes()
			testConfig := CephConfig{}
			AdmConfig := map[string]string{
				"fsid":     "b29fd",
				"mon_host": "[0.0.0.0]",
			}
			err := testConfig.ConfigFromAdm(ioMock, "ceph.conf")
			Expect(err).NotTo(HaveOccurred())
			Expect(testConfig.GetAdmConf()).To(Equal(AdmConfig))
		})
	})

	Describe("[SecretFromAdm Test]", func() {
		It("Parse keyring to AdmSecret", func() {
			ioMock.EXPECT().ReadFile(gomock.Any()).DoAndReturn(
				func(filename string) ([]byte, error) {
					secret := []byte("[client.admin]\n\tkey = b29fd")
					return secret, nil
				}).AnyTimes()
			testConfig := CephConfig{}
			AdmSecret := map[string][]byte{
				"keyring": []byte("[client.admin]\n\tkey = b29fd"),
			}
			err := testConfig.SecretFromAdm(ioMock, "keyring")
			Expect(err).NotTo(HaveOccurred())
			Expect(testConfig.GetAdmSecret()).To(Equal(AdmSecret))
		})
	})

	Describe("[MakeIniFile Test]", func() {
		It("Make Ini file from Map", func() {
			ioMock.EXPECT().WriteFile(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
				func(fileName string, data []byte, fileMode os.FileMode) error {
					return nil
				}).AnyTimes()

			testConfig := CephConfig{
				crConf: map[string]string{
					"debug_osd": "20/20",
				},
			}
			//ini := "[global]\n\tdebug_osd = 20/20\n"
			err := testConfig.MakeIniFile(ioMock, "ceph.conf")
			Expect(err).NotTo(HaveOccurred())
			//Expect(retini).To(Equal(ini))
		})
	})

	Describe("[INCOMPLETE][UpdateConfToK8s Test]", func() {
		It("should return nil", func() {
			kubeMock.EXPECT().InClusterConfig().Return(nil, nil).AnyTimes()

			kubeMock.EXPECT().NewForConfig(gomock.Any()).DoAndReturn(
				func(c interface{}) (kubernetes.Interface, error) {
					fakeClient := fake.NewSimpleClientset()
					fakeClient.PrependReactor("get", "configmaps", func(action testing.Action) (bool, runtime.Object, error) {
						return true, nil, nil
					})

					fakeClient.PrependReactor("create", "configmaps", func(action testing.Action) (bool, runtime.Object, error) {
						return true, nil, nil
					})

					fakeClient.PrependReactor("update", "configmaps", func(action testing.Action) (bool, runtime.Object, error) {
						return true, nil, nil
					})
					return fakeClient, nil
				}).AnyTimes()

			testConfig := CephConfig{
				crConf: map[string]string{
					"test": "test",
				},
			}
			err := testConfig.UpdateConfToK8s(kubeMock)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("[INCOMPLETE][UpdateKeyringToK8s Test]", func() {
		It("should return nil", func() {
			kubeMock.EXPECT().InClusterConfig().Return(nil, nil).AnyTimes()

			kubeMock.EXPECT().NewForConfig(gomock.Any()).DoAndReturn(
				func(c interface{}) (kubernetes.Interface, error) {
					fakeClient := fake.NewSimpleClientset()
					fakeClient.PrependReactor("get", "secrets", func(action testing.Action) (bool, runtime.Object, error) {
						return true, nil, nil
					})
					fakeClient.PrependReactor("update", "secrets", func(action testing.Action) (bool, runtime.Object, error) {
						return true, nil, nil
					})
					return fakeClient, nil
				}).AnyTimes()
			testConfig := CephConfig{
				admSecret: map[string][]byte{
					"keyring": {0, 0, 0, 0},
				},
			}
			err := testConfig.UpdateKeyringToK8s(kubeMock)
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
