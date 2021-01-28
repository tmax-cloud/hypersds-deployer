package config

import (
	"context"
	"fmt"
	"hypersds-provisioner/pkg/common/wrapper"
	"strings"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CephConfigInterface interface {
	ConfigFromAdm(wrapper.IoUtilInterface, string) error
	SecretFromAdm(wrapper.IoUtilInterface, string) error
	MakeIni() string
	PutConfigToK8s(kubeWrapper wrapper.KubeInterface) error
	PutSecretToK8s(kubeWrapper wrapper.KubeInterface) error

	SetCrConf(map[string]string) error
	SetAdmConf(map[string]string) error
	SetAdmSecret(map[string][]byte) error

	GetCrConf() (map[string]string, error)
	GetAdmConf() (map[string]string, error)
	GetAdmSecret(map[string][]byte, error)
}

func (cconf *CephConfig) SetCrConf(m map[string]string) error {
	cconf.CrConf = m
	return nil
}

func (cconf *CephConfig) SetAdmConf(m map[string]string) error {
	cconf.AdmConf = m
	return nil
}

func (cconf *CephConfig) SetAdmSecret(m map[string][]byte) error {
	cconf.AdmSecret = m
	return nil
}

func (cconf *CephConfig) GetCrConf() (map[string]string, error) {
	return cconf.CrConf, nil
}

func (cconf *CephConfig) GetAdmConf() (map[string]string, error) {
	return cconf.AdmConf, nil
}

func (cconf *CephConfig) GetAdmSecret() (map[string][]byte, error) {
	return cconf.AdmSecret, nil
}

type CephConfig struct {
	CrConf    map[string]string
	AdmConf   map[string]string
	AdmSecret map[string][]byte
}

func NewConfigFromCephCr(CephCr hypersdsv1alpha1.CephClusterSpec) *CephConfig {
	cconf := CephConfig{}
	CrConf := make(map[string]string)

	for key, value := range CephCr.Config {
		CrConf[key] = value
	}

	cconf.SetCrConf(CrConf)
	return &cconf
}

func (cconf *CephConfig) ConfigFromAdm(IoUtil wrapper.IoUtilInterface, cephconf string) error {
	AdmConf := make(map[string]string)

	dat, err := IoUtil.ReadFile(cephconf)
	if err != nil {
		return err
	}

	lines := strings.Split(string(dat[:]), "\n")
	for _, s := range lines {
		kv := strings.Split(s, "=")
		if len(kv) > 1 {
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			AdmConf[key] = val
		}
	}
	return cconf.SetAdmConf(AdmConf)
}

func (cconf *CephConfig) SecretFromAdm(IoUtil wrapper.IoUtilInterface, cephsecret string) error {
	AdmSecret := make(map[string][]byte)

	dat, err := IoUtil.ReadFile(cephsecret)
	if err != nil {
		return err
	}

	AdmSecret["keyring"] = dat
	return cconf.SetAdmSecret(AdmSecret)
}

func (cconf *CephConfig) MakeIni() string {
	ini := "[global]\n"
	for key, value := range cconf.CrConf {
		s1 := fmt.Sprintf("\t%s = %s\n", key, value)
		ini = fmt.Sprintf("%s%s", ini, s1)
	}
	return ini
}

//default:default SA가 configmap put할수 있는 rolebinding 있어야함!
func (cconf *CephConfig) PutConfigToK8s(kubeWrapper wrapper.KubeInterface) error {
	//creates the in-cluster config
	config, err := kubeWrapper.InClusterConfig()
	if err != nil {
		return err
	}

	clientset, err := kubeWrapper.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	AdmConf, _ := cconf.GetAdmConf()
	configMap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "ceph-conf",
		},
		Data: AdmConf,
	}
	_, err = clientset.CoreV1().ConfigMaps("default").Get(context.TODO(), "ceph-conf", metav1.GetOptions{})
	if err == nil {
		_, err = clientset.CoreV1().ConfigMaps("default").Update(context.TODO(), &configMap, metav1.UpdateOptions{})

	}
	return err
}

func (cconf *CephConfig) PutSecretToK8s(kubeWrapper wrapper.KubeInterface) error {
	//creates the in-cluster config
	config, err := kubeWrapper.InClusterConfig()
	if err != nil {
		return err
	}

	clientset, err := kubeWrapper.NewForConfig(config)
	if err != nil {
		return err
	}

	AdmSecret, _ := cconf.GetAdmSecret()
	secret := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "ceph-secret",
		},
		Data: AdmSecret,
	}

	_, err = clientset.CoreV1().Secrets("default").Get(context.TODO(), "ceph-secret", metav1.GetOptions{})
	if err == nil {
		_, err = clientset.CoreV1().Secrets("default").Update(context.TODO(), &secret, metav1.UpdateOptions{})
	}
	return err
}
