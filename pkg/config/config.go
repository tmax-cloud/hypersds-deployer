package config

import (
	"context"
	"fmt"
	"hypersds-provisioner/pkg/common/wrapper"
	"strings"

	"github.com/juju/errors"
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CephConfigInterface interface {
	ConfigFromAdm(wrapper.IoUtilInterface, string) error
	MakeIni() string
	PutConfigToK8s() error

	SetCrConf(map[string]string) error
	SetAdmConf(map[string]string) error

	GetCrConf() (map[string]string, error)
	GetAdmConf() (map[string]string, error)
}

func (cconf *CephConfig) SetCrConf(m map[string]string) error {
	cconf.CrConf = m
	return nil
}

func (cconf *CephConfig) SetAdmConf(m map[string]string) error {
	cconf.AdmConf = m
	return nil
}

func (cconf *CephConfig) GetCrConf() (map[string]string, error) {
	return cconf.CrConf, nil
}

func (cconf *CephConfig) GetAdmConf() (map[string]string, error) {
	return cconf.AdmConf, nil
}

type CephConfig struct {
	CrConf  map[string]string
	AdmConf map[string]string
}

func NewConfigFromCephCr(CephCr hypersdsv1alpha1.CephClusterSpec) *CephConfig {
	cconf := CephConfig{}
	cconf.CrConf = make(map[string]string)

	for key, value := range CephCr.Config {
		cconf.CrConf[key] = value
	}
	return &cconf
}

func (cconf *CephConfig) ConfigFromAdm(IoUtil wrapper.IoUtilInterface, cephconf string) error {
	if cconf.AdmConf == nil {
		cconf.AdmConf = make(map[string]string)
	}

	dat, err := IoUtil.ReadFile(cephconf)
	if err != nil {
		return err
	}

	lines := strings.Split(string(dat[:]), "\n")
	for _, s := range lines {
		kv := strings.Split(s, "=")
		if len(kv) > 1 {
			key := strings.Trim(kv[0], " ")
			val := strings.Trim(kv[1], " ")
			cconf.AdmConf[key] = val
		}
	}
	return nil
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
		panic(err)
	}

	clientset, err := kubeWrapper.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	configMap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "ceph-cluster",
		},
		Data: cconf.AdmConf,
	}

	if _, err := clientset.CoreV1().ConfigMaps("default").Get(context.TODO(), "ceph-cluster", metav1.GetOptions{}); errors.IsNotFound(err) {
		_, err = clientset.CoreV1().ConfigMaps("default").Create(context.TODO(), &configMap, metav1.CreateOptions{})
	} else {
		_, err = clientset.CoreV1().ConfigMaps("default").Update(context.TODO(), &configMap, metav1.UpdateOptions{})
	}
	return err
}
