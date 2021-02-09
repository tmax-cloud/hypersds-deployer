package config

import (
	"context"
	"fmt"
	"hypersds-provisioner/pkg/common/wrapper"
	"strings"

	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	K8sNamespace = "default"
	k8sConfigMap = "ceph-conf"
	K8sSecret    = "ceph-secret"
)

type CephConfigInterface interface {
	ConfigFromAdm(wrapper.IoUtilInterface, string) error
	SecretFromAdm(wrapper.IoUtilInterface, string) error
	MakeIniFile(wrapper.IoUtilInterface, string) error
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

func (c *ConfigInitStruct) NewConfigFromCephCr(cephCr hypersdsv1alpha1.CephClusterSpec) (*CephConfig, error) {
	conf := CephConfig{}
	crConf := make(map[string]string)

	for key, value := range cephCr.Config {
		crConf[key] = value
	}

	err := conf.SetCrConf(crConf)
	return &conf, err
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

func (conf *CephConfig) MakeIniFile(ioUtil wrapper.IoUtilInterface, fileName string) error {
	ini := "[global]\n"
	crConf, _ := conf.GetCrConf()
	for key, value := range crConf {
		s1 := fmt.Sprintf("\t%s = %s\n", key, value)
		ini = fmt.Sprintf("%s%s", ini, s1)
	}

	buf := []byte(ini)
	err := ioUtil.WriteFile(fileName, buf, 0644)

	return err
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

	configMap, err := clientset.CoreV1().ConfigMaps(K8sNamespace).Get(context.TODO(), k8sConfigMap, metav1.GetOptions{})
	if err != nil {
		return err
	}

	admConf, _ := conf.GetAdmConf()
	configMap.Data = admConf
	_, err = clientset.CoreV1().ConfigMaps(K8sNamespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
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

	secret, err := clientset.CoreV1().Secrets(K8sNamespace).Get(context.TODO(), K8sSecret, metav1.GetOptions{})
	if err != nil {
		return err
	}

	admSecret, _ := conf.GetAdmSecret()
	secret.Data = admSecret
	_, err = clientset.CoreV1().Secrets(K8sNamespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	return err
}
