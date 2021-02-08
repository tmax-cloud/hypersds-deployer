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

	"context"
	"flag"
	"fmt"
	"path/filepath"
	"reflect"
)

const (
	cephConfName        = "ceph-conf"
	cephKeyringName     = "ceph-secret"
	cephRoleName        = "ceph-role"
	cephRoleBindingName = "ceph-rolebinding"
	// TODO: change to own SA and NS
	cephServiceAccountName = "default"
	cephNamespace          = "default"
)

var _ = Describe("Bootstrap Test", func() {
	defer GinkgoRecover()

	var (
		err       error
		clientSet *kubernetes.Clientset
		nodeName  string // required in multinode k8s environment
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
				Name:      cephConfName,
				Namespace: cephNamespace,
			},
		}

		createdCm, err := clientSet.CoreV1().ConfigMaps(cephNamespace).Create(context.TODO(), &cephConfCm, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())
		fmt.Println("return type:", reflect.TypeOf(createdCm))

		fmt.Println("cm result: ", createdCm)

		cephKeyringSecret := corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      cephKeyringName,
				Namespace: cephNamespace,
			},
		}

		createdSecret, err := clientSet.CoreV1().Secrets(cephNamespace).Create(context.TODO(), &cephKeyringSecret, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		fmt.Println("return type:", reflect.TypeOf(createdSecret))
		fmt.Println("secret result: ", createdSecret)

		provisionerRole := rbacv1.Role{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Role",
				APIVersion: "rbac.authorization.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      cephRoleName,
				Namespace: cephNamespace,
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs:     []string{"get", "update"},
					APIGroups: []string{""},
					Resources: []string{"configmaps", "secrets"},
				},
			},
		}

		createdRole, err := clientSet.RbacV1().Roles(cephNamespace).Create(context.TODO(), &provisionerRole, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		fmt.Println("return type:", reflect.TypeOf(createdRole))
		fmt.Println("secret result: ", createdRole)

		provisionerRoleBinding := rbacv1.RoleBinding{
			TypeMeta: metav1.TypeMeta{
				Kind:       "RoleBinding",
				APIVersion: "rbac.authorization.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      cephRoleBindingName,
				Namespace: cephNamespace,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      cephServiceAccountName,
					Namespace: cephNamespace,
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     cephRoleName,
			},
		}

		createdRoleBinding, err := clientSet.RbacV1().RoleBindings(cephNamespace).Create(context.TODO(), &provisionerRoleBinding, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		fmt.Println("return type:", reflect.TypeOf(createdRoleBinding))
		fmt.Println("secret result: ", createdRoleBinding)
	})

	AfterEach(func() {
		deletePolicy := metav1.DeletePropagationForeground
		err = clientSet.CoreV1().ConfigMaps(cephNamespace).Delete(context.TODO(), cephConfName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		Expect(err).NotTo(HaveOccurred())

		err = clientSet.CoreV1().Secrets(cephNamespace).Delete(context.TODO(), cephKeyringName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		Expect(err).NotTo(HaveOccurred())

		err = clientSet.RbacV1().Roles(cephNamespace).Delete(context.TODO(), cephRoleName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		Expect(err).NotTo(HaveOccurred())

		err = clientSet.RbacV1().RoleBindings(cephNamespace).Delete(context.TODO(), cephRoleBindingName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		Expect(err).NotTo(HaveOccurred())
	})

	It("is simple e2e test case", func() {
		// XXX: Change it as one's environment
		nodeName = "master1"
		err = runProvisionerContainer(clientSet, nodeName)

		// Check bootstrap successfully completed
		Expect(err).NotTo(HaveOccurred())

		// Check ConfigMap and Secret are successfully updated
		cephConfCm, err := clientSet.CoreV1().ConfigMaps(cephNamespace).Get(context.TODO(), cephConfName, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		cmData := cephConfCm.Data
		Expect(cmData).NotTo(BeEmpty())

		cephKeyringSecret, err := clientSet.CoreV1().Secrets(cephNamespace).Get(context.TODO(), cephKeyringName, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		secretData := cephKeyringSecret.Data
		Expect(secretData).NotTo(BeEmpty())
	})
})
