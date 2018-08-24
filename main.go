package main

import (
	"fmt"

	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/client/k8srestconfig"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func main() {
	var err error

	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			panic(err)
		}
	}

	var restConfig *rest.Config
	{
		c := k8srestconfig.Config{
			Logger: logger,

			Address:   "https://127.0.0.1:53078",
			InCluster: false,
			TLS: k8srestconfig.TLSClientConfig{
				CAFile:  "/Users/xh3b4sd/.config/opsctl/ginger/certs/1534240785/ca.pem",
				CrtFile: "/Users/xh3b4sd/.config/opsctl/ginger/certs/1534240785/crt.pem",
				KeyFile: "/Users/xh3b4sd/.config/opsctl/ginger/certs/1534240785/key.pem",
			},
		}

		restConfig, err = k8srestconfig.New(c)
		if err != nil {
			panic(err)
		}
	}

	g8sClient, err := versioned.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}

	n := ""
	o := metav1.GetOptions{}
	_, err = g8sClient.ProviderV1alpha1().AWSConfigs(n).Get("", o)
	if errors.IsBadRequest(err) {
		fmt.Printf("IsBadRequest\n")
	} else if errors.IsInvalid(err) {
		fmt.Printf("IsInvalid\n")
	} else if err != nil {
		fmt.Printf("%#v\n", err)
		panic(err)
	}
}
