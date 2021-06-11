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

	users, err := (*c).GetUsers()

	if err != nil {
		return nil, err
	}

	return users, nil
}
