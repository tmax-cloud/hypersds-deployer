package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"context"
	"flag"
	"fmt"
	"path/filepath"
	"reflect"
)

const (
	cephConfName    = "ceph-conf"
	cephKeyringName = "ceph-secret"
)

var _ = Describe("Bootstrap Test", func() {
	defer GinkgoRecover()

	var (
		err       error
		clientSet *kubernetes.Clientset
	)

	BeforeEach(func() {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		kubeConfigWithFlag, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		Expect(err).NotTo(HaveOccurred())

		clientSet, err = kubernetes.NewForConfig(kubeConfigWithFlag)
		Expect(err).NotTo(HaveOccurred())

		cephConfCm := corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: cephConfName,
			},
		}

		createdCm, err := clientSet.CoreV1().ConfigMaps("default").Create(context.TODO(), &cephConfCm, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())
		fmt.Println("return type:", reflect.TypeOf(createdCm))

		fmt.Println("cm result: ", createdCm)

		cephKeyringSecret := corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: cephKeyringName,
			},
		}

		createdSecret, err := clientSet.CoreV1().Secrets("default").Create(context.TODO(), &cephKeyringSecret, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())
		fmt.Println("return type:", reflect.TypeOf(createdSecret))

		fmt.Println("secret result: ", createdSecret)
	})

	AfterEach(func() {
		deletePolicy := metav1.DeletePropagationForeground
		err = clientSet.CoreV1().ConfigMaps("default").Delete(context.TODO(), cephConfName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		Expect(err).NotTo(HaveOccurred())

		err = clientSet.CoreV1().Secrets("default").Delete(context.TODO(), cephKeyringName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		Expect(err).NotTo(HaveOccurred())
	})

	It("is simple e2e test case", func() {
		err = runProvisionerContainer(clientSet, "master1")

		// Check bootstrap successfully completed
		Expect(err).NotTo(HaveOccurred())

		// Check ConfigMap and Secret are successfully updated
		cephConfCm, err := clientSet.CoreV1().ConfigMaps("default").Get(context.TODO(), cephConfName, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		cmData := cephConfCm.Data
		Expect(cmData).NotTo(BeEmpty())

		cephKeyringSecret, err := clientSet.CoreV1().Secrets("default").Get(context.TODO(), cephKeyringName, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		secretData := cephKeyringSecret.Data
		Expect(secretData).NotTo(BeEmpty())
	})
})
