package openshift

import (
	"context"
	"os"

	userv1 "github.com/openshift/api/user/v1"
	userclientv1 "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var ocConfig *rest.Config

func GetClientConfig() (*rest.Config, error) {
	if ocConfig != nil {
		return ocConfig, nil
	}

	kubeconfig := os.Getenv("KUBECONFIG_PATH")

	ocConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return ocConfig, nil
}

func GetUsers() (*[]userv1.User, error) {
	conf, err := GetClientConfig()

	if err != nil {
		return nil, err
	}

	userV1Client, err := userclientv1.NewForConfig(conf)
	if err != nil {
		//return err
	}

	users, err := userV1Client.Users().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return &users.Items, nil
}
