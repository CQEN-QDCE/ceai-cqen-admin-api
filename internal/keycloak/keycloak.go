package keycloak

import (
	"context"
	"errors"
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

	var briefRep = false
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

func GetUser(username string) (*gocloak.User, error) {
	c, err := GetClient()

	if err != nil {
		return nil, err
	}

	var briefRep = false
	ctx := context.Background()
	users, err := (*c.client).GetUsers(
		ctx,
		c.token.AccessToken,
		c.realm,
		gocloak.GetUsersParams{
			BriefRepresentation: &briefRep,
			Username:            &username,
		})

	if err != nil {
		return nil, err
	}

	if len(users) < 1 {
		err := errors.New("Username not found.")
		return nil, err
	}

	return users[0], nil
}

func CreateUser(user *gocloak.User) error {
	c, err := GetClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	_, err = (*c.client).CreateUser(
		ctx,
		c.token.AccessToken,
		c.realm,
		*user)

	return err
}

func GetUserRoles(user *gocloak.User) ([]*gocloak.Role, error) {
	c, err := GetClient()

	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	roles, err := (*c.client).GetCompositeRealmRolesByUserID(
		ctx,
		c.token.AccessToken,
		c.realm,
		*user.ID)

	if err != nil {
		return nil, err
	}

	return roles, nil
}

func GetUserGroups(user *gocloak.User) ([]*gocloak.Group, error) {
	c, err := GetClient()

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var briefRep = true

	groups, err := (*c.client).GetUserGroups(
		ctx,
		c.token.AccessToken,
		c.realm,
		*user.ID,
		gocloak.GetGroupsParams{
			BriefRepresentation: &briefRep,
		})

	if err != nil {
		return nil, err
	}

	return groups, nil
}

func GetGroup(groupName string) (*gocloak.Group, error) {
	c, err := GetClient()

	if err != nil {
		return nil, err
	}

	var briefRep = true
	ctx := context.Background()

	groups, err := (*c.client).GetGroups(
		ctx,
		c.token.AccessToken,
		c.realm,
		gocloak.GetGroupsParams{
			BriefRepresentation: &briefRep,
			Search:              &groupName,
		})

	if err != nil {
		return nil, err
	}

	if len(groups) < 1 {
		err := errors.New("Group not found.")
		return nil, err
	}

	if len(groups) > 1 {
		err := errors.New("Group name is not unique.")
		return nil, err
	}

	return groups[0], nil
}

func GetGroupMembers(groupName string) ([]*gocloak.User, error) {
	c, err := GetClient()

	if err != nil {
		return nil, err
	}

	var briefRep = true
	ctx := context.Background()

	users, err := (*c.client).GetGroupMembers(
		ctx,
		c.token.AccessToken,
		c.realm,
		groupName,
		gocloak.GetGroupsParams{
			BriefRepresentation: &briefRep,
		})

	if err != nil {
		return nil, err
	}

	return users, nil
}
