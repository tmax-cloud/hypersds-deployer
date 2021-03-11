package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"io/ioutil"

	"gopkg.in/yaml.v2"

	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

type BootstrapInput struct {
	CephConfigMapName       string `yaml:"cephConfigMapName" default:"ceph-conf"`
	CephKeyringName         string `yaml:"cephKeyringName" default:"ceph-secret"`
	CephRoleName            string `yaml:"cephRoleName" default:"ceph-role"`
	CephRoleBindingName     string `yaml:"cephRoleBindingName" default:"ceph-rolebinding"`
	CephServiceAccountName  string `yaml:"cephServiceAccountName" default:"default"`
	CephNamespace           string `yaml:"cephNamespace" default:"default"`
	CephProvisionerImage    string `yaml:"cephProvisionerImage" default:"hypersds-provisioner:test"`
	CephProvisionerNodeName string `yaml:"cephProvisionerNodeName,omitempty"`
	RegistryCredentialName  string `yaml:"registryCredentialName,omitempty"`
	TestManifestDir         string `yaml:"testManifestDir,omitempty"`
}

const (
	//testWorkspaceDir    = "/e2e"
	inputDir    = "inputs"   // directory to use in test. required.
	hostPathDir = "manifest" // directory to use in test. required.
	inputFile   = "bootstrap.yaml"
)

var _ = Describe("[E2e] Bootstrap Test", func() {
	defer GinkgoRecover()

	var (
		err              error
		clientSet        *kubernetes.Clientset
		bootstrapInput   BootstrapInput
		testWorkspaceDir string
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

		// Open bootstrap input yaml file
		testWorkspaceDir, err = os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		inputFilePath := filepath.Join(testWorkspaceDir, inputDir, inputFile)
		fmt.Println("Opening file ", inputFilePath)
		source, err := ioutil.ReadFile(inputFilePath)
		Expect(err).NotTo(HaveOccurred())
		fmt.Println(source)

		err = yaml.Unmarshal(source, &bootstrapInput)
		Expect(err).NotTo(HaveOccurred())

		cephConfCm := corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      bootstrapInput.CephConfigMapName,
				Namespace: bootstrapInput.CephNamespace,
			},
		}

		createdCm, err := clientSet.CoreV1().ConfigMaps(bootstrapInput.CephNamespace).Create(context.TODO(), &cephConfCm, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())
		fmt.Println("return type:", reflect.TypeOf(createdCm))

		fmt.Println("cm result: ", createdCm)

		cephKeyringSecret := corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      bootstrapInput.CephKeyringName,
				Namespace: bootstrapInput.CephNamespace,
			},
		}

		createdSecret, err := clientSet.CoreV1().Secrets(bootstrapInput.CephNamespace).Create(context.TODO(), &cephKeyringSecret, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		fmt.Println("return type:", reflect.TypeOf(createdSecret))
		fmt.Println("secret result: ", createdSecret)

		provisionerRole := rbacv1.Role{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Role",
				APIVersion: "rbac.authorization.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      bootstrapInput.CephRoleName,
				Namespace: bootstrapInput.CephNamespace,
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs:     []string{"get", "update"},
					APIGroups: []string{""},
					Resources: []string{"configmaps", "secrets"},
				},
			},
		}

		createdRole, err := clientSet.RbacV1().Roles(bootstrapInput.CephNamespace).Create(context.TODO(), &provisionerRole, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		fmt.Println("return type:", reflect.TypeOf(createdRole))
		fmt.Println("secret result: ", createdRole)

		provisionerRoleBinding := rbacv1.RoleBinding{
			TypeMeta: metav1.TypeMeta{
				Kind:       "RoleBinding",
				APIVersion: "rbac.authorization.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      bootstrapInput.CephRoleBindingName,
				Namespace: bootstrapInput.CephNamespace,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      bootstrapInput.CephServiceAccountName,
					Namespace: bootstrapInput.CephNamespace,
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     bootstrapInput.CephRoleName,
			},
		}

		createdRoleBinding, err := clientSet.RbacV1().RoleBindings(bootstrapInput.CephNamespace).Create(context.TODO(), &provisionerRoleBinding, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		fmt.Println("return type:", reflect.TypeOf(createdRoleBinding))
		fmt.Println("secret result: ", createdRoleBinding)
	})

	AfterEach(func() {
		deletePolicy := metav1.DeletePropagationForeground
		err = clientSet.CoreV1().ConfigMaps(bootstrapInput.CephNamespace).Delete(context.TODO(), bootstrapInput.CephConfigMapName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		Expect(err).NotTo(HaveOccurred())

		err = clientSet.CoreV1().Secrets(bootstrapInput.CephNamespace).Delete(context.TODO(), bootstrapInput.CephKeyringName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		Expect(err).NotTo(HaveOccurred())

		err = clientSet.RbacV1().Roles(bootstrapInput.CephNamespace).Delete(context.TODO(), bootstrapInput.CephRoleName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		Expect(err).NotTo(HaveOccurred())

		err = clientSet.RbacV1().RoleBindings(bootstrapInput.CephNamespace).Delete(context.TODO(), bootstrapInput.CephRoleBindingName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		Expect(err).NotTo(HaveOccurred())
	})

	It("is simple e2e test case", func() {
		if bootstrapInput.TestManifestDir != "" {
			err = runProvisionerContainer(clientSet,
				bootstrapInput.CephNamespace,
				bootstrapInput.CephProvisionerImage,
				bootstrapInput.TestManifestDir,
				bootstrapInput.RegistryCredentialName,
				bootstrapInput.CephProvisionerNodeName)

			// Check bootstrap successfully completed
			Expect(err).NotTo(HaveOccurred())
		} else {
			testManifestDir := filepath.Join(testWorkspaceDir, inputDir, hostPathDir)
			err = runProvisionerContainer(clientSet,
				bootstrapInput.CephNamespace,
				bootstrapInput.CephProvisionerImage,
				testManifestDir,
				bootstrapInput.RegistryCredentialName,
				bootstrapInput.CephProvisionerNodeName)

			// Check bootstrap successfully completed
			Expect(err).NotTo(HaveOccurred())
		}

		// Check ConfigMap and Secret are successfully updated
		cephConfCm, err := clientSet.CoreV1().ConfigMaps(bootstrapInput.CephNamespace).Get(context.TODO(), bootstrapInput.CephConfigMapName, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		cmData := cephConfCm.Data
		Expect(cmData).NotTo(BeEmpty())

		cephKeyringSecret, err := clientSet.CoreV1().Secrets(bootstrapInput.CephNamespace).Get(context.TODO(), bootstrapInput.CephKeyringName, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		secretData := cephKeyringSecret.Data
		Expect(secretData).NotTo(BeEmpty())
	})
})
