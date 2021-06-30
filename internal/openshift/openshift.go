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

func GetUserClient() (*userclientv1.UserV1Client, error) {
	conf, err := GetClientConfig()

	if err != nil {
		return nil, err
	}

	return userclientv1.NewForConfig(conf)
}

func GetUsers() (*[]userv1.User, error) {
	userV1Client, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	users, err := userV1Client.Users().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return &users.Items, nil
}

func GetUser(username string) (*userv1.User, error) {
	userV1Client, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	return userV1Client.Users().Get(context.TODO(), username, metav1.GetOptions{})
}

func CreateUser(user *userv1.User) (*userv1.User, error) {
	userV1Client, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	return userV1Client.Users().Create(context.TODO(), user, metav1.CreateOptions{})
}

func UpdateUser(user *userv1.User) (*userv1.User, error) {
	userV1Client, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	return userV1Client.Users().Update(context.TODO(), user, metav1.UpdateOptions{})
}

func DeleteUser(user *userv1.User) error {
	userV1Client, err := GetUserClient()
	if err != nil {
		return err
	}

	return userV1Client.Users().Delete(context.TODO(), user.Name, metav1.DeleteOptions{})
}

func AddUserInGroup(userName string, groupName string) error {
	userV1Client, err := GetUserClient()
	if err != nil {
		return err
	}

	group, err := userV1Client.Groups().Get(context.TODO(), groupName, metav1.GetOptions{})

	inGroup, _ := UserInGroup(userName, group)

	if !inGroup {
		group.Users = append(group.Users, userName)
		_, err = userV1Client.Groups().Update(context.TODO(), group, metav1.UpdateOptions{})
	}

	return err
}

//Check if a user is in a group
func UserInGroup(username string, group *userv1.Group) (bool, int) {
	for i, user := range group.Users {
		if user == username {
			return true, i
		}
	}

	return false, -1
}
