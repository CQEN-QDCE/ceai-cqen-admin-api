package aws

import (
	"net/http"
	"os"

	scim "github.com/CQEN-QDCE/aws-sso-scim-goclient"
)

var client *scim.Client

func GetScimClient() (*scim.Client, error) {
	if client != nil {
		return client, nil
	}

	endpoint := os.Getenv("SCIM_ENDPOINT")
	token := os.Getenv("SCIM_TOKEN")

	awsClient, err := scim.NewClient(
		&http.Client{},
		&scim.Config{
			Endpoint: endpoint,
			Token:    token,
		})

	if err != nil {
		return nil, err
	}

	client := &awsClient

	return client, nil
}

func GetUsers() ([]*scim.User, error) {
	c, err := GetScimClient()

	if err != nil {
		return nil, err
	}

	return (*c).GetUsers()
}

func GetUser(username string) (*scim.User, error) {
	c, err := GetScimClient()

	if err != nil {
		return nil, err
	}

	//In AWS SCIM usernames are emails
	return (*c).FindUserByEmail(username)
}

func CreateUser(user *scim.User) (*scim.User, error) {
	c, err := GetScimClient()

	if err != nil {
		return nil, err
	}

	return (*c).CreateUser(user)
}

func UpdateUser(user *scim.User) (*scim.User, error) {
	c, err := GetScimClient()

	if err != nil {
		return nil, err
	}

	return (*c).UpdateUser(user)
}

func DeleteUser(user *scim.User) error {
	c, err := GetScimClient()

	if err != nil {
		return err
	}

	return (*c).DeleteUser(user)
}

func GetGroup(groupname string) (*scim.Group, error) {
	c, err := GetScimClient()

	if err != nil {
		return nil, err
	}

	return (*c).FindGroupByDisplayName(groupname)
}

func CreateGroup(group *scim.Group) (*scim.Group, error) {
	c, err := GetScimClient()

	if err != nil {
		return nil, err
	}

	return (*c).CreateGroup(group)
}

func AddUserToGroup(user *scim.User, group *scim.Group) error {
	c, err := GetScimClient()

	if err != nil {
		return err
	}

	return (*c).AddUserToGroup(user, group)
}
