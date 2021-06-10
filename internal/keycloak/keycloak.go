package keycloak

import (
	"context"
	"os"

	"github.com/Nerzal/gocloak/v8"
)

var client *Client

type Client struct {
	client *gocloak.GoCloak
	token  *gocloak.JWT
	realm  string
}

func GetClient() (*Client, error) {

	if client != nil {
		return client, nil
	}

	clientId := os.Getenv("KEYCLOAK_CLIENT_ID")
	secret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	realm := os.Getenv("KEYCLOAK_REALM")

	url := os.Getenv("KEYCLOAK_URL")

	grantType := "client_credentials"

	kcclient := gocloak.NewClient(url)
	ctx := context.Background()

	token, err := kcclient.GetToken(ctx, realm, gocloak.TokenOptions{
		ClientID:     &clientId,
		ClientSecret: &secret,
		GrantType:    &grantType,
	})
	if err != nil {
		return nil, err
	}

	client = &Client{
		client: &kcclient,
		token:  token,
		realm:  realm,
	}

	return client, nil
}

func GetUsers() ([]*gocloak.User, error) {
	c, err := GetClient()

	if err != nil {
		return nil, err
	}

	var briefRep = true
	ctx := context.Background()

	users, err := (*c.client).GetUsers(
		ctx,
		c.token.AccessToken,
		c.realm,
		gocloak.GetUsersParams{
			BriefRepresentation: &briefRep,
		})

	if err != nil {
		return nil, err
	}

	return users, nil
}
