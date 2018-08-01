package main

import (
	"encoding/json"
	"fmt"

	providerv1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/client/k8srestconfig"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

const (
	SelfLink = "/apis/provider.giantswarm.io/v1alpha1/namespaces/default/awsconfigs/6aben/status"
)

type Patch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type Status struct {
	Cluster providerv1alpha1.StatusCluster `json:"cluster" yaml:"cluster"`
}

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

			Address:   "https://127.0.0.1:52565",
			InCluster: false,
			TLS: k8srestconfig.TLSClientConfig{
				CAFile:  "/Users/xh3b4sd/.config/opsctl/gauss/certs/1508360771/ca.pem",
				CrtFile: "/Users/xh3b4sd/.config/opsctl/gauss/certs/1508360771/crt.pem",
				KeyFile: "/Users/xh3b4sd/.config/opsctl/gauss/certs/1508360771/key.pem",
			},
		}

		restConfig, err = k8srestconfig.New(c)
		if err != nil {
			panic(err)
		}
	}

	var restClient rest.Interface
	{
		g8sClient, err := versioned.NewForConfig(restConfig)
		if err != nil {
			panic(err)
		}
		restClient = g8sClient.ProviderV1alpha1().RESTClient()
	}

	obj, err := restClient.Get().AbsPath(SelfLink).Do().Get()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", obj)

	var patches []Patch

	patches = append(patches, Patch{
		Op:   "add",
		Path: "/status",
		Value: Status{
			Cluster: providerv1alpha1.StatusCluster{
				Conditions: []providerv1alpha1.StatusClusterCondition{},
				Nodes:      []providerv1alpha1.StatusClusterNode{},
				Versions:   []providerv1alpha1.StatusClusterVersion{},
			},
		},
	})

	b, err := json.Marshal(patches)
	if err != nil {
		panic(err)
	}

	err = restClient.Patch(types.JSONPatchType).AbsPath(SelfLink).Body(b).Do().Error()
	if errors.IsConflict(err) {
		panic(err)
	} else if errors.IsResourceExpired(err) {
		panic(err)
	} else if err != nil {
		panic(err)
	}

	fmt.Printf("\n")

	obj, err = restClient.Get().AbsPath(SelfLink).Do().Get()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", obj)
}
