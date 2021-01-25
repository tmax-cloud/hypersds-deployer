package config

import (
	"fmt"
	"io/ioutil"
	"strings"
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	hypersdsv1alpha1 "github.com/tmax-cloud/hypersds-operator/api/v1alpha1"
)


type cephconf interface {
	ConfigFromCephCr(hypersdsv1aplha1.CephCluster) *CephConfig
}

type cconf struct {}

func (c *cconf) ConfigFromCephCr(cr hyersdsv1aplha1.CephCluster) *CephConfig{
	return CephConfig.ConfigFromCephCr(cr)
}


type CephConfigInterface interface {
	ConfigFromAdm(IoUtilInterface, string) error
	MakeIni() (string)
	PutConfigToK8s() (error) 

	SetCrConf(map[string]string) error
	SetAdmConf(map[string]string) error

	GetCrConf() (map[string]string, error)
	GetAdmConf() (map[string]string, error)
}

func (cconf *CephConfig) SetCrConf(m map[string]string) error {
	cconf.CrConf = m
	return nil
}

func (cconf *CephConfig)SetAdmConf(m map[string]string) error {
	cconf.AdmConf = m
	return nil
}

func (cconf *CephConfig)GetCrConf() (map[string]string, error) {
	return cconf.CrConf, nil
}

func (cconf *CephConfig)GetAdmConf() (map[string]string, error) {
	return cconf.AdmConf, nil
}



type CephConfig struct {
	CrConf map[string]string
	AdmConf map[string]string
}

func ConfigFromCephCr(CephCr hypersdsv1alpha1.CephCluster) *CephConfig{
	cconf := CephConfig{}
	cconf.Conf = make(map[string]string)

	for key, value := range CephCr.Config {
		cconf.Conf[key] = value
	}
	return &cconf
}

func (cconf *CephConfig) ConfigFromAdm(IoUtil IoUtilInterface, cephconf string) error {
	if cconf.AdmConf == nil {
		cconf.AdmConf = make(map[string]string)
	}

	dat, err := IoUtil.ReadFile(cephconf)
	if err != nil {
		return err
	}

	lines := strings.Split(string(dat[:]),"\n")
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

func (cconf *CephConfig) MakeIni() string{
	ini := "[global]\n"
	for key, value := range cconf.Conf {
		s1 := fmt.Sprintf("\t%s = %s\n",key,value)
		ini = fmt.Sprintf("%s%s",ini, s1)
	}
	return ini
}

//default:default SA가 configmap put할수 있는 rolebinding 있어야함!
func (cconf *CephConfig) PutConfigToK8s() error{
	//creates the in-cluster config
	config, err := rest.InCluserConfig()
	if err != nil {
		panic(err)
	}

	clusterset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	configMap := corev1.ConfigMap {
		TypeMeta: metav1.TypeMeta {
			Kind: "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta {
			Name: "ceph-cluster"
		},
		Data: cconf.AdmConf,
	}

	if _, err := clientset.CoreV1.ConfigMaps("default").Get(context.TODO, "ceph-cluster", metav1.GetOptions{}); errors.IsNotFound(err) {
		_, err = clientset.CoreV1.ConfigMap("default").Create(context.TODO(),&configMap, metav1.CreateOptions{})
	} else {
		_, err = clientset.CoreV1.ConfigMap("default").Update(context.TODO(), &configMap, metav1.UpdateOptions{})
	}
	return err
}
