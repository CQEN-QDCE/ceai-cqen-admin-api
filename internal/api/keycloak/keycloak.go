package keycloak

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/Nerzal/gocloak/v11"
)

const CLIENT_TOKEN_TTL = 60

var serviceAccountClient *ServiceAccountClient
var serviceAccountClientTime int64

type ClientCredentials struct {
	Id     string
	Secret string
	Realm  string
}

type ServiceAccountClient struct {
	GoCloakClient *gocloak.GoCloak
	Token         *gocloak.JWT
	Realm         string
}

func GetClientCredentials() *ClientCredentials {
	creds := ClientCredentials{
		Realm:  os.Getenv("KEYCLOAK_REALM"),
		Id:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		Secret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
	}

	return &creds
}

func GetGoCloakClient() *gocloak.GoCloak {
	url := os.Getenv("KEYCLOAK_URL")

	client := gocloak.NewClient(url)

	return &client
}

func GetServiceAccountClient() (*ServiceAccountClient, error) {

	if serviceAccountClient != nil && (time.Now().Unix()-serviceAccountClientTime) < CLIENT_TOKEN_TTL {
		return serviceAccountClient, nil
	}

	goCloakClient := GetGoCloakClient()
	clientCreds := GetClientCredentials()
	grantType := "client_credentials"
	ctx := context.Background()

	token, err := (*goCloakClient).GetToken(ctx, clientCreds.Realm, gocloak.TokenOptions{
		ClientID:     &clientCreds.Id,
		ClientSecret: &clientCreds.Secret,
		GrantType:    &grantType,
	})

	if err != nil {
		return nil, err
	}

	serviceAccountClient = &ServiceAccountClient{
		GoCloakClient: goCloakClient,
		Token:         token,
		Realm:         clientCreds.Realm,
	}

	serviceAccountClientTime = time.Now().Unix()

	return serviceAccountClient, nil
}

func LoginOtp(username string, password string, totp string) (*gocloak.JWT, error) {
	goCloakClient := GetGoCloakClient()
	clientCreds := GetClientCredentials()

	return (*goCloakClient).LoginOtp(context.Background(), clientCreds.Id, clientCreds.Secret, clientCreds.Realm, username, password, totp)
}

func RefreshToken(refreshToken string) (*gocloak.JWT, error) {
	goCloakClient := GetGoCloakClient()
	clientCreds := GetClientCredentials()

	return (*goCloakClient).RefreshToken(context.Background(), refreshToken, clientCreds.Id, clientCreds.Secret, clientCreds.Realm)
}

func GetUsers() ([]*gocloak.User, error) {
	c, err := GetServiceAccountClient()

	if err != nil {
		return nil, err
	}

	var briefRep = false
	ctx := context.Background()

	users, err := (*c.GoCloakClient).GetUsers(
		ctx,
		c.Token.AccessToken,
		c.Realm,
		gocloak.GetUsersParams{
			BriefRepresentation: &briefRep,
		})

	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUser(username string) (*gocloak.User, error) {
	c, err := GetServiceAccountClient()

	if err != nil {
		return nil, err
	}

	var briefRep = false
	ctx := context.Background()
	users, err := (*c.GoCloakClient).GetUsers(
		ctx,
		c.Token.AccessToken,
		c.Realm,
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
	c, err := GetServiceAccountClient()

	if err != nil {
		return "", err
	}

	ctx := context.Background()

	return (*c.GoCloakClient).CreateUser(
		ctx,
		c.Token.AccessToken,
		c.Realm,
		*user)
}

func UpdateUser(user *gocloak.User) error {
	c, err := GetServiceAccountClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	return (*c.GoCloakClient).UpdateUser(
		ctx,
		c.Token.AccessToken,
		c.Realm,
		*user)
}

func DeleteUser(userID string) error {
	c, err := GetServiceAccountClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	return (*c.GoCloakClient).DeleteUser(
		ctx,
		c.Token.AccessToken,
		c.Realm,
		userID)
}

func GetUserRoles(user *gocloak.User) ([]*gocloak.Role, error) {
	c, err := GetServiceAccountClient()

	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	roles, err := (*c.GoCloakClient).GetCompositeRealmRolesByUserID(
		ctx,
		c.Token.AccessToken,
		c.Realm,
		*user.ID)

	if err != nil {
		return nil, err
	}

	return roles, nil
}

func GetUserGroups(user *gocloak.User) ([]*gocloak.Group, error) {
	c, err := GetServiceAccountClient()

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var briefRep = false

	groups, err := (*c.GoCloakClient).GetUserGroups(
		ctx,
		c.Token.AccessToken,
		c.Realm,
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
	c, err := GetServiceAccountClient()

	if err != nil {
		return nil, err
	}

	var briefRep = false
	ctx := context.Background()

	groups, err := (*c.GoCloakClient).GetGroups(
		ctx,
		c.Token.AccessToken,
		c.Realm,
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

func GetGroupById(id string) (*gocloak.Group, error) {
	c, err := GetServiceAccountClient()

	if err != nil {
		return nil, err
	}

	return (*c.GoCloakClient).GetGroup(
		context.Background(),
		c.Token.AccessToken,
		c.Realm,
		id)
}

func GetGroups(parentGroupName *string) (*[]gocloak.Group, error) {
	if parentGroupName == nil {
		parentGroupName = gocloak.StringP("/")
	}

	group, err := GetGroup(*parentGroupName)

	if group != nil && err == nil {
		//For some reasons Keycloak won't always provide subgroups on GetGroups search call
		//So we get group full info with a GetGroup(idGroup) call
		fullGroup, err := GetGroupById(*group.ID)

		if err != nil {
			return nil, err
		}

		return fullGroup.SubGroups, err
	} else {
		err := errors.New("Group name does not exist.")
		return nil, err
	}
}

func GetGroupMembers(group *gocloak.Group) ([]*gocloak.User, error) {
	c, err := GetServiceAccountClient()

	if err != nil {
		return nil, err
	}

	var briefRep = true
	ctx := context.Background()

	users, err := (*c.GoCloakClient).GetGroupMembers(
		ctx,
		c.Token.AccessToken,
		c.Realm,
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
	c, err := GetServiceAccountClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	return (*c.GoCloakClient).AddUserToGroup(
		ctx,
		c.Token.AccessToken,
		c.Realm,
		*user.ID,
		*group.ID,
	)
}

func DeleteUserFromGroup(user *gocloak.User, group *gocloak.Group) error {
	c, err := GetServiceAccountClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	return (*c.GoCloakClient).DeleteUserFromGroup(
		ctx,
		c.Token.AccessToken,
		c.Realm,
		*user.ID,
		*group.ID,
	)
}

func CreateGroup(group *gocloak.Group) error {
	c, err := GetServiceAccountClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	_, err = (*c.GoCloakClient).CreateGroup(
		ctx,
		c.Token.AccessToken,
		c.Realm,
		*group,
	)

	return err
}

func CreateChildGroup(parentGroup *gocloak.Group, group *gocloak.Group) error {
	c, err := GetServiceAccountClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	_, err = (*c.GoCloakClient).CreateChildGroup(
		ctx,
		c.Token.AccessToken,
		c.Realm,
		*parentGroup.ID,
		*group,
	)

	return err
}

func UpdateGroup(group *gocloak.Group) error {
	c, err := GetServiceAccountClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	return (*c.GoCloakClient).UpdateGroup(
		ctx,
		c.Token.AccessToken,
		c.Realm,
		*group,
	)
}

func DeleteGroup(group *gocloak.Group) error {
	c, err := GetServiceAccountClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	return (*c.GoCloakClient).DeleteGroup(
		ctx,
		c.Token.AccessToken,
		c.Realm,
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
	c, err := GetServiceAccountClient()

	if err != nil {
		return err
	}
	ctx := context.Background()

	return (*c.GoCloakClient).ExecuteActionsEmail(
		ctx,
		c.Token.AccessToken,
		c.Realm,
		gocloak.ExecuteActionsEmail{
			UserID:   &userID,
			Lifespan: gocloak.IntP(86400), //24h
			Actions:  &actions,
		})
}
