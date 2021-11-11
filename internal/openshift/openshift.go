package openshift

import (
	"context"
	"os"

	authorization "github.com/openshift/api/authorization/v1"
	project "github.com/openshift/api/project/v1"
	user "github.com/openshift/api/user/v1"
	authorizationclient "github.com/openshift/client-go/authorization/clientset/versioned/typed/authorization/v1"
	projectclient "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	userclient "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sclient "k8s.io/client-go/kubernetes"
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

func GetUserClient() (*userclient.UserV1Client, error) {
	conf, err := GetClientConfig()

	if err != nil {
		return nil, err
	}

	return userclient.NewForConfig(conf)
}

func GetProjectClient() (*projectclient.ProjectV1Client, error) {
	conf, err := GetClientConfig()

	if err != nil {
		return nil, err
	}

	return projectclient.NewForConfig(conf)
}

func GetAuthorizationClient() (*authorizationclient.AuthorizationV1Client, error) {
	conf, err := GetClientConfig()

	if err != nil {
		return nil, err
	}

	return authorizationclient.NewForConfig(conf)
}

func GetK8sClient() (*k8sclient.Clientset, error) {
	conf, err := GetClientConfig()

	if err != nil {
		return nil, err
	}

	return k8sclient.NewForConfig(conf)
}

func GetUsers() (*[]user.User, error) {
	userClient, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	users, err := userClient.Users().List(context.TODO(), meta.ListOptions{})

	if err != nil {
		return nil, err
	}

	return &users.Items, nil
}

func GetUser(username string) (*user.User, error) {
	userClient, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	return userClient.Users().Get(context.TODO(), username, meta.GetOptions{})
}

func CreateUser(user *user.User) (*user.User, error) {
	userClient, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	return userClient.Users().Create(context.TODO(), user, meta.CreateOptions{})
}

func UpdateUser(user *user.User) (*user.User, error) {
	userClient, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	return userClient.Users().Update(context.TODO(), user, meta.UpdateOptions{})
}

func DeleteUser(user *user.User) error {
	userClient, err := GetUserClient()
	if err != nil {
		return err
	}

	return userClient.Users().Delete(context.TODO(), user.Name, meta.DeleteOptions{})
}

func GetGroup(groupName string) (*user.Group, error) {
	userClient, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	return userClient.Groups().Get(context.TODO(), groupName, meta.GetOptions{})
}

func CreateGroup(group *user.Group) (*user.Group, error) {
	userClient, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	return userClient.Groups().Create(context.TODO(), group, meta.CreateOptions{})
}

func UpdateGroup(group *user.Group) (*user.Group, error) {
	userClient, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	return userClient.Groups().Update(context.TODO(), group, meta.UpdateOptions{})
}

func DeleteGroup(group *user.Group) error {
	userClient, err := GetUserClient()
	if err != nil {
		return err
	}

	return userClient.Groups().Delete(context.TODO(), group.Name, meta.DeleteOptions{})
}

func AddUserInGroup(userName string, groupName string) error {
	userClient, err := GetUserClient()
	if err != nil {
		return err
	}

	group, err := userClient.Groups().Get(context.TODO(), groupName, meta.GetOptions{})

	inGroup, _ := UserInGroup(userName, group)

	if !inGroup {
		group.Users = append(group.Users, userName)
		_, err = userClient.Groups().Update(context.TODO(), group, meta.UpdateOptions{})
	}

	return err
}

func RemoveUserFromGroup(userName string, groupName string) error {
	userClient, err := GetUserClient()
	if err != nil {
		return err
	}

	group, err := userClient.Groups().Get(context.TODO(), groupName, meta.GetOptions{})

	inGroup, pos := UserInGroup(userName, group)

	if inGroup {
		if len(group.Users) > 1 {
			group.Users = append(group.Users[:pos], group.Users[pos+1:]...) //Removing an array element in go...
		} else {
			group.Users = []string{} //Replace with a empty array if it only has our user in it
		}

		_, err = userClient.Groups().Update(context.TODO(), group, meta.UpdateOptions{})
	}

	return err
}

//Check if a user is in a group
func UserInGroup(username string, group *user.Group) (bool, int) {
	for i, user := range group.Users {
		if user == username {
			return true, i
		}
	}

	return false, -1
}

func GetProject(projectName string) (*project.Project, error) {
	projectClient, err := GetProjectClient()
	if err != nil {
		return nil, err
	}

	return projectClient.Projects().Get(context.TODO(), projectName, meta.GetOptions{})
}

func CreateProject(project *project.Project) (*project.Project, error) {
	projectClient, err := GetProjectClient()
	if err != nil {
		return nil, err
	}

	return projectClient.Projects().Create(context.TODO(), project, meta.CreateOptions{})
}

func UpdateProject(project *project.Project) (*project.Project, error) {
	projectClient, err := GetProjectClient()
	if err != nil {
		return nil, err
	}

	return projectClient.Projects().Update(context.TODO(), project, meta.UpdateOptions{})
}

func GetNamespace(projectName string) (*core.Namespace, error) {
	k8sClient, err := GetK8sClient()
	if err != nil {
		return nil, err
	}

	return k8sClient.CoreV1().Namespaces().Get(context.TODO(), projectName, meta.GetOptions{})
}

func UpdateNamespace(namespace *core.Namespace) (*core.Namespace, error) {
	k8sClient, err := GetK8sClient()
	if err != nil {
		return nil, err
	}

	return k8sClient.CoreV1().Namespaces().Update(context.TODO(), namespace, meta.UpdateOptions{})
}

func GetNamespaceRoleBindings(namespace string) (*[]authorization.RoleBinding, error) {
	authorizationClient, err := GetAuthorizationClient()
	if err != nil {
		return nil, err
	}

	roleList, err := authorizationClient.RoleBindings(namespace).List(context.TODO(), meta.ListOptions{FieldSelector: ""})

	return &roleList.Items, nil
}

func CreateRoleBinding(namespace string, roleBinding *authorization.RoleBinding) (*authorization.RoleBinding, error) {
	authorizationClient, err := GetAuthorizationClient()
	if err != nil {
		return nil, err
	}

	return authorizationClient.RoleBindings(namespace).Create(context.TODO(), roleBinding, meta.CreateOptions{})
}

func DeleteRoleBinding(namespace string, roleBinding *authorization.RoleBinding) error {
	authorizationClient, err := GetAuthorizationClient()
	if err != nil {
		return err
	}

	return authorizationClient.RoleBindings(namespace).Delete(context.TODO(), roleBinding.Name, meta.DeleteOptions{})
}
