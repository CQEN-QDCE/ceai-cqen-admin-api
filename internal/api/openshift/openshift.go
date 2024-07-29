package openshift

import (
	"context"
	"log"
	"os"
	"strconv"

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
var isOCNonPersist bool

const ERR_RESOURCE_NAME_MAY_NOT_BE_EMPTY = "resource name may not be empty"
const ERR_RESOURCE_USER_NOT_FOUND_PREFIX = "users.user.openshift.io \""
const ERR_RESOURCE_GROUP_NOT_FOUND_PREFIX = "groups.user.openshift.io \""
const ERR_RESOURCE_NAMESPACE_NOT_FOUND_PREFIX = "namespaces \""
const ERR_RESOURCE_NOT_FOUND_SUFFIX = "\" not found"

func GetClientConfig() (*rest.Config, error) {
	if ocConfig != nil {
		return ocConfig, nil
	}

	kubeconfig := os.Getenv("KUBECONFIG_PATH")

	ocConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	//Openshift Non Persist?
	strOcNonPersist, existsOcNonPersist := os.LookupEnv("OC_NON_PERSIST")
	if existsOcNonPersist {
		_isOCNonPersist, _ := strconv.ParseBool(strOcNonPersist)
		if _isOCNonPersist {
			isOCNonPersist = true
		}
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

	_user, err := userClient.Users().Get(context.TODO(), username, meta.GetOptions{})
	errorMessage := ERR_RESOURCE_USER_NOT_FOUND_PREFIX + username + ERR_RESOURCE_NOT_FOUND_SUFFIX
	if IsErrorWithMessageAndOcNonPersist(errorMessage, err) {
		return _user, nil
	}

	return _user, err
}

func CreateUser(user *user.User) (*user.User, error) {
	userClient, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	return userClient.Users().Create(context.TODO(), user, GetCreateOptions())
}

func UpdateUser(user *user.User) (*user.User, error) {
	userClient, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	_usr, err := userClient.Users().Update(context.TODO(), user, GetUpdateOptions())
	if IsErrorWithMessageAndOcNonPersist(ERR_RESOURCE_NAME_MAY_NOT_BE_EMPTY, err) {
		return _usr, nil
	}

	return _usr, err
}

func DeleteUser(user *user.User) error {
	userClient, err := GetUserClient()
	if err != nil {
		return err
	}

	_err := userClient.Users().Delete(context.TODO(), user.Name, GetDeleteOptions())
	if IsErrorWithMessageAndOcNonPersist(ERR_RESOURCE_NAME_MAY_NOT_BE_EMPTY, _err) {
		return nil
	}

	return _err
}

func GetGroup(groupName string) (*user.Group, error) {
	userClient, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	_group, err := userClient.Groups().Get(context.TODO(), groupName, meta.GetOptions{})
	errorMessage := ERR_RESOURCE_GROUP_NOT_FOUND_PREFIX + groupName + ERR_RESOURCE_NOT_FOUND_SUFFIX
	if IsErrorWithMessageAndOcNonPersist(errorMessage, err) {
		return _group, nil
	}

	return _group, err
}

func CreateGroup(group *user.Group) (*user.Group, error) {
	userClient, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	return userClient.Groups().Create(context.TODO(), group, GetCreateOptions())
}

func UpdateGroup(group *user.Group) (*user.Group, error) {
	userClient, err := GetUserClient()
	if err != nil {
		return nil, err
	}

	_group, err := userClient.Groups().Update(context.TODO(), group, GetUpdateOptions())
	if IsErrorWithMessageAndOcNonPersist(ERR_RESOURCE_NAME_MAY_NOT_BE_EMPTY, err) {
		return _group, nil
	}

	return _group, err
}

func DeleteGroup(group *user.Group) error {
	userClient, err := GetUserClient()
	if err != nil {
		return err
	}

	return userClient.Groups().Delete(context.TODO(), group.Name, GetDeleteOptions())
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
		_, err = userClient.Groups().Update(context.TODO(), group, GetUpdateOptions())
		if IsErrorWithMessageAndOcNonPersist(ERR_RESOURCE_NAME_MAY_NOT_BE_EMPTY, err) {
			return nil
		}
	}

	return err
}

func RemoveUserFromGroup(userName string, groupName string) error {

	userClient, err := GetUserClient()
	if err != nil {
		return err
	}

	group, err := userClient.Groups().Get(context.TODO(), groupName, meta.GetOptions{})
	if IsErrorWithMessageAndOcNonPersist(ERR_RESOURCE_NAME_MAY_NOT_BE_EMPTY, err) {
		return nil
	}

	inGroup, pos := UserInGroup(userName, group)

	if inGroup {
		if len(group.Users) > 1 {
			group.Users = append(group.Users[:pos], group.Users[pos+1:]...) //Removing an array element in go...
		} else {
			group.Users = []string{} //Replace with a empty array if it only has our user in it
		}

		_, err = userClient.Groups().Update(context.TODO(), group, GetUpdateOptions())
		if IsErrorWithMessageAndOcNonPersist(ERR_RESOURCE_NAME_MAY_NOT_BE_EMPTY, err) {
			return nil
		}
	}

	return err
}

// Check if a user is in a group
func UserInGroup(username string, group *user.Group) (bool, int) {
	for i, user := range group.Users {
		if user == username {
			return true, i
		}
	}

	return false, -1
}

func GetProjects() (*[]project.Project, error) {
	projectClient, err := GetProjectClient()
	if err != nil {
		return nil, err
	}

	projects, err := projectClient.Projects().List(context.TODO(), meta.ListOptions{})

	if err != nil {
		return nil, err
	}

	return &projects.Items, nil
}

func GetProject(projectName string) (*project.Project, error) {
	projectClient, err := GetProjectClient()
	if err != nil {
		return nil, err
	}

	_project, err := projectClient.Projects().Get(context.TODO(), projectName, meta.GetOptions{})
	errorMessage := ERR_RESOURCE_NAMESPACE_NOT_FOUND_PREFIX + projectName + ERR_RESOURCE_NOT_FOUND_SUFFIX
	if IsErrorWithMessageAndOcNonPersist(errorMessage, err) {
		_project = GetDummyProject(projectName)
		return _project, nil
	}

	return _project, err
}

func CreateProject(project *project.ProjectRequest) (*project.Project, error) {
	projectClient, err := GetProjectClient()
	if err != nil {
		return nil, err
	}

	return projectClient.ProjectRequests().Create(context.TODO(), project, GetCreateOptions())
}

func UpdateProject(project *project.Project) (*project.Project, error) {

	projectClient, err := GetProjectClient()
	if err != nil {
		return nil, err
	}

	_project, err := projectClient.Projects().Update(context.TODO(), project, GetUpdateOptions())
	errorMessage := ERR_RESOURCE_NAMESPACE_NOT_FOUND_PREFIX + project.Name + ERR_RESOURCE_NOT_FOUND_SUFFIX
	if IsErrorWithMessageAndOcNonPersist(errorMessage, err) {
		_project = GetDummyProject(project.Name)
		return _project, nil
	}

	return _project, err
}

func DeleteProject(project *project.Project) error {

	projectClient, err := GetProjectClient()
	if err != nil {
		return err
	}

	_err := projectClient.Projects().Delete(context.TODO(), project.Name, GetDeleteOptions())

	errorMessage := ERR_RESOURCE_NAMESPACE_NOT_FOUND_PREFIX + project.Name + ERR_RESOURCE_NOT_FOUND_SUFFIX
	if IsErrorWithMessageAndOcNonPersist(errorMessage, err) {
		return nil
	}

	return _err
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

	return k8sClient.CoreV1().Namespaces().Update(context.TODO(), namespace, GetUpdateOptions())
}

func GetNamespaceRoleBindings(namespace string) (*[]authorization.RoleBinding, error) {
	authorizationClient, err := GetAuthorizationClient()
	if err != nil {
		return nil, err
	}

	roleList, err := authorizationClient.RoleBindings(namespace).List(context.TODO(), meta.ListOptions{FieldSelector: ""})

	if err != nil {
		return nil, err
	}

	return &roleList.Items, nil
}

func CreateRoleBinding(namespace string, roleBinding *authorization.RoleBinding) (*authorization.RoleBinding, error) {
	authorizationClient, err := GetAuthorizationClient()
	if err != nil {
		return nil, err
	}

	return authorizationClient.RoleBindings(namespace).Create(context.TODO(), roleBinding, GetCreateOptions())
}

func DeleteRoleBinding(namespace string, roleBinding *authorization.RoleBinding) error {

	authorizationClient, err := GetAuthorizationClient()
	if err != nil {
		return err
	}

	return authorizationClient.RoleBindings(namespace).Delete(context.TODO(), roleBinding.Name, GetDeleteOptions())
}

func GetCreateOptions() meta.CreateOptions {
	opts := meta.CreateOptions{}
	if isOCNonPersist {
		opts.DryRun = append(opts.DryRun, meta.DryRunAll)
	}
	return opts
}

func GetUpdateOptions() meta.UpdateOptions {
	opts := meta.UpdateOptions{}
	if isOCNonPersist {
		opts.DryRun = append(opts.DryRun, meta.DryRunAll)
	}
	return opts
}

func GetDeleteOptions() meta.DeleteOptions {
	opts := meta.DeleteOptions{}
	if isOCNonPersist {
		opts.DryRun = append(opts.DryRun, meta.DryRunAll)
	}
	return opts
}

/**
Methods to support a non Openshift environment
**/

func IsErrorWithMessageAndOcNonPersist(errorMessage string, err error) bool {
	if err != nil {
		if err.Error() == errorMessage {
			if isOCNonPersist {
				log.Println("IsSpecifiedErrorMessageAndOcNonPersist " + errorMessage + " is true!!!")
				return true
			}
		}
	}
	return false
}

func GetDummyProject(projectName string) *project.Project {
	var dummyProject *project.Project = new(project.Project)

	dummyProject.Name = projectName

	mapAnnotations := map[string]string{
		"openshift.io/description":  "some description",
		"openshift.io/display-name": "some displayname",
	}
	dummyProject.Annotations = mapAnnotations

	return dummyProject
}
