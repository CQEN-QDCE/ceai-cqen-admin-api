package keycloak

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/Nerzal/gocloak/v9"
)

const CLIENT_TOKEN_TTL = 60

var client *Client
var clientTime int64

type Client struct {
	client *gocloak.GoCloak
	token  *gocloak.JWT
	realm  string
}

func GetClient() (*Client, error) {

	if client != nil && (time.Now().Unix()-clientTime) < CLIENT_TOKEN_TTL {
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

	clientTime = time.Now().Unix()

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

	//Get groups because Keycloak won't get them in its User endpoint
	groups, err := GetUserGroups(users[0])
	if err != nil {
		return nil, err
	}

	groupList := make([]string, len(groups))
	for _, group := range groups {
		groupList = append(groupList, *group.Path)
	}

	users[0].Groups = &groupList

	//Get roles because Keycloak won't get them either
	roles, err := GetUserRoles(users[0])
	if err != nil {
		return nil, err
	}

	roleList := make([]string, len(groups))
	for _, role := range roles {
		roleList = append(roleList, *role.Name)
	}

	users[0].RealmRoles = &roleList

	return users[0], nil
}

//Create a new User, returns it's ID
func CreateUser(user *gocloak.User) (string, error) {
	c, err := GetClient()

	if err != nil {
		return "", err
	}

	ctx := context.Background()

	return (*c.client).CreateUser(
		ctx,
		c.token.AccessToken,
		c.realm,
		*user)
}

func UpdateUser(user *gocloak.User) error {
	c, err := GetClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	return (*c.client).UpdateUser(
		ctx,
		c.token.AccessToken,
		c.realm,
		*user)
}

func DeleteUser(userID string) error {
	c, err := GetClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	return (*c.client).DeleteUser(
		ctx,
		c.token.AccessToken,
		c.realm,
		userID)
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

//Recursive
func FindSubgroup(group *gocloak.Group, subgroupName string) *gocloak.Group {
	if group != nil {
		if *group.Name == subgroupName {
			return group
		} else if group.SubGroups != nil {
			for _, subgroup := range *group.SubGroups {
				if foundGroup := FindSubgroup(&subgroup, subgroupName); foundGroup != nil {
					return foundGroup
				}
			}
		}
	}

	return nil
}

func GetGroup(groupName string) (*gocloak.Group, error) {
	c, err := GetClient()

	if err != nil {
		return nil, err
	}

	var briefRep = false
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

	//If group is a subgroup Keycloak will return the whole tree containing the group
	group := FindSubgroup(groups[0], groupName)

	if group == nil {
		return nil, errors.New("Group not found in tree.")
	}

	return group, nil
}

func GetGroups(parentGroupName *string) (*[]gocloak.Group, error) {
	if parentGroupName == nil {
		parentGroupName = gocloak.StringP("/")
	}

	group, err := GetGroup(*parentGroupName)

	if group != nil {
		return group.SubGroups, err
	} else {
		err := errors.New("Group name does not exist.")
		return nil, err
	}
}

func GetGroupMembers(group *gocloak.Group) ([]*gocloak.User, error) {
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
		*group.ID,
		gocloak.GetGroupsParams{
			BriefRepresentation: &briefRep,
		})

	if err != nil {
		return nil, err
	}

	return users, nil
}

//Idempotent
func AddUserToGroup(user *gocloak.User, group *gocloak.Group) error {
	c, err := GetClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	return (*c.client).AddUserToGroup(
		ctx,
		c.token.AccessToken,
		c.realm,
		*user.ID,
		*group.ID,
	)
}

func DeleteUserFromGroup(user *gocloak.User, group *gocloak.Group) error {
	c, err := GetClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	return (*c.client).DeleteUserFromGroup(
		ctx,
		c.token.AccessToken,
		c.realm,
		*user.ID,
		*group.ID,
	)
}

func CreateGroup(group *gocloak.Group) error {
	c, err := GetClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	_, err = (*c.client).CreateGroup(
		ctx,
		c.token.AccessToken,
		c.realm,
		*group,
	)

	return err
}

func CreateChildGroup(parentGroup *gocloak.Group, group *gocloak.Group) error {
	c, err := GetClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	_, err = (*c.client).CreateChildGroup(
		ctx,
		c.token.AccessToken,
		c.realm,
		*parentGroup.ID,
		*group,
	)

	return err
}

func UpdateGroup(group *gocloak.Group) error {
	c, err := GetClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	return (*c.client).UpdateGroup(
		ctx,
		c.token.AccessToken,
		c.realm,
		*group,
	)
}

func DeleteGroup(group *gocloak.Group) error {
	c, err := GetClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	return (*c.client).DeleteGroup(
		ctx,
		c.token.AccessToken,
		c.realm,
		*group.ID,
	)
}

//Send an email with a link to complete all required actions for a user
func ExecuteCurrentActionEmail(username string) error {
	user, err := GetUser(username)

	if err != nil {
		return err
	}

	return ExecuteActionEmail(*user.ID, *user.RequiredActions)
}

func ExecuteActionEmail(userID string, actions []string) error {
	c, err := GetClient()

	if err != nil {
		return err
	}
	ctx := context.Background()

	return (*c.client).ExecuteActionsEmail(
		ctx,
		c.token.AccessToken,
		c.realm,
		gocloak.ExecuteActionsEmail{
			UserID:   &userID,
			Lifespan: gocloak.IntP(86400), //24h
			Actions:  &actions,
		})
}
