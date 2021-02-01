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
	UpdateConfToK8s(kubeWrapper wrapper.KubeInterface) error
	UpdateKeyringToK8s(kubeWrapper wrapper.KubeInterface) error

	SetCrConf(map[string]string) error
	SetAdmConf(map[string]string) error
	SetAdmSecret(map[string][]byte) error

	GetCrConf() (map[string]string, error)
	GetAdmConf() (map[string]string, error)
	GetAdmSecret(map[string][]byte, error)
}

func (conf *CephConfig) SetCrConf(m map[string]string) error {
	conf.crConf = m
	return nil
}

func (conf *CephConfig) SetAdmConf(m map[string]string) error {
	conf.admConf = m
	return nil
}

func (conf *CephConfig) SetAdmSecret(m map[string][]byte) error {
	conf.admSecret = m
	return nil
}

func (conf *CephConfig) GetCrConf() (map[string]string, error) {
	return conf.crConf, nil
}

func (conf *CephConfig) GetAdmConf() (map[string]string, error) {
	return conf.admConf, nil
}

func (conf *CephConfig) GetAdmSecret() (map[string][]byte, error) {
	return conf.admSecret, nil
}

type CephConfig struct {
	crConf    map[string]string
	admConf   map[string]string
	admSecret map[string][]byte
}

type ConfigInitStruct struct{}

func (c *ConfigInitStruct) NewConfigFromCephCr(cephCr hypersdsv1alpha1.CephClusterSpec) *CephConfig {
	conf := CephConfig{}
	crConf := make(map[string]string)

	for key, value := range cephCr.Config {
		crConf[key] = value
	}

	err := conf.SetCrConf(crConf)
	if err != nil {
		return nil
	}
	return &conf
}

func (conf *CephConfig) ConfigFromAdm(ioUtil wrapper.IoUtilInterface, cephconf string) error {
	admConf := make(map[string]string)

	dat, err := ioUtil.ReadFile(cephconf)
	if err != nil {
		return err
	}

	lines := strings.Split(string(dat[:]), "\n")
	for _, s := range lines {
		kv := strings.Split(s, "=")
		if len(kv) > 1 {
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			admConf[key] = val
		}
	}
	return conf.SetAdmConf(admConf)
}

func (conf *CephConfig) SecretFromAdm(ioUtil wrapper.IoUtilInterface, cephsecret string) error {
	admSecret := make(map[string][]byte)

	dat, err := ioUtil.ReadFile(cephsecret)
	if err != nil {
		return err
	}

	admSecret["keyring"] = dat
	return conf.SetAdmSecret(admSecret)
}

func (conf *CephConfig) MakeIni() string {
	ini := "[global]\n"
	crConf, _ := conf.GetCrConf()
	for key, value := range crConf {
		s1 := fmt.Sprintf("\t%s = %s\n", key, value)
		ini = fmt.Sprintf("%s%s", ini, s1)
	}
	return ini
}

//default:default SA가 configmap put할수 있는 rolebinding 있어야함!
func (conf *CephConfig) UpdateConfToK8s(kubeWrapper wrapper.KubeInterface) error {
	//creates the in-cluster config
	config, err := kubeWrapper.InClusterConfig()
	if err != nil {
		return err
	}

	clientset, err := kubeWrapper.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	admConf, _ := conf.GetAdmConf()
	configMap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "ceph-conf",
		},
		Data: admConf,
	}
	_, err = clientset.CoreV1().ConfigMaps("default").Get(context.TODO(), "ceph-conf", metav1.GetOptions{})
	if err != nil {
		return err
	}
	_, err = clientset.CoreV1().ConfigMaps("default").Update(context.TODO(), &configMap, metav1.UpdateOptions{})
	return err
}

func (conf *CephConfig) UpdateKeyringToK8s(kubeWrapper wrapper.KubeInterface) error {
	//creates the in-cluster config
	config, err := kubeWrapper.InClusterConfig()
	if err != nil {
		return err
	}

	clientset, err := kubeWrapper.NewForConfig(config)
	if err != nil {
		return err
	}

	admSecret, _ := conf.GetAdmSecret()
	secret := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "ceph-secret",
		},
		Data: admSecret,
	}

	_, err = clientset.CoreV1().Secrets("default").Get(context.TODO(), "ceph-secret", metav1.GetOptions{})
	if err != nil {
		return err
	}
	_, err = clientset.CoreV1().Secrets("default").Update(context.TODO(), &secret, metav1.UpdateOptions{})
	return err
}
