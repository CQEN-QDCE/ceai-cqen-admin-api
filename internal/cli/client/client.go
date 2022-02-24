package client

import (
	"context"
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
	"github.com/getkin/kin-openapi/openapi3"
)

func GetClient() (*apifirst.Client, error) {
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}

	OpenAPIDoc, err := loader.LoadFromFile("/home/gagf01/projets/ceai-cqen-admin-api/api/openapi-v1.yaml")
	if err != nil {
		return nil, fmt.Errorf("Error loading OpenAPI Spec file: %s", err.Error())
	}

	if err = OpenAPIDoc.Validate(ctx); err != nil {
		return nil, fmt.Errorf("Invalid OpenAPI Spec file: " + err.Error())
	}

	client := apifirst.NewClient(OpenAPIDoc)

	return client, nil
}

func GetClientToUrl(serverUrl string) (*apifirst.Client, error) {
	client, err := GetClient()

	if err != nil {
		return nil, err
	}

	err = client.SetServerUrl(serverUrl)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetAuthenticatedClient() (*apifirst.Client, error) {
	session, err := GetSession()

	if err != nil {
		return nil, err
	}

	client, err := GetClient()

	if err != nil {
		return nil, err
	}

	err = client.SetServerUrl(session.ServerUrl)

	if err != nil {
		return nil, err
	}

	client.SetAuthorization(session.Token.AccessToken)

	return client, nil
}
