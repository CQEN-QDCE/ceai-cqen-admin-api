package services

import (
	"strings"

	scim "github.com/CQEN-QDCE/aws-sso-scim-goclient"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/aws"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/keycloak"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/openshift"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/Nerzal/gocloak/v11"
	userv1 "github.com/openshift/api/user/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ADMIN_ROLE_NAME = "Admin"
const DEV_ROLE_NAME = "Developer"

type UserState struct {
	models.UserWithLabs
	Keycloak  *gocloak.User
	Aws       *scim.User
	Openshift *userv1.User
}

func GetKeycloakAdminGroup() (*gocloak.Group, error) {
	return keycloak.GetGroup(ADMIN_ROLE_NAME)
}

// Gets current User state across all products: Keycloak|AWS|Openshift
func GetUserState(username string) (*UserState, error) {
	var state UserState
	var kerr, aerr, oerr error

	fKeycloak := func() {
		state.Keycloak, kerr = keycloak.GetUser(username)
	}

	fAws := func() {
		state.Aws, aerr = aws.GetUser(username)
	}

	fOpenshift := func() {
		state.Openshift, oerr = openshift.GetUser(username)
	}

	Parallelize(fKeycloak, fAws, fOpenshift)

	if kerr != nil || oerr != nil || aerr != nil {
		var err error

		if kerr != nil {
			err = NewErrorExternalRessourceNotFound(kerr, ERROR_SERVER_KEYCLOAK)
		} else if oerr != nil {
			err = NewErrorExternalRessourceNotFound(oerr, ERROR_SERVER_OPENSHIFT)
		} else if aerr != nil {
			err = NewErrorExternalRessourceNotFound(aerr, ERROR_SERVER_AWS)
		}

		return nil, err
	}

	if state.Keycloak != nil {
		state.UserWithLabs = *mapKeycloakUserWithLabs(state.Keycloak)
	}

	return &state, nil
}

func mapKeycloakUser(kuser *gocloak.User) *models.User {
	var user models.User

	if kuser.Email != nil {
		user.Email = *kuser.Email //email is used as username or id
	}

	if kuser.FirstName != nil {
		user.Firstname = *kuser.FirstName
	}

	if kuser.LastName != nil {
		user.Lastname = *kuser.LastName
	}

	user.Disabled = gocloak.BoolP(!gocloak.PBool(kuser.Enabled))

	//Iterate RealmRoles to find infraRole (Admin || Developer)
	if kuser.RealmRoles != nil {
		for _, role := range *kuser.RealmRoles {
			if role == ADMIN_ROLE_NAME {
				user.Infrarole = ADMIN_ROLE_NAME
			}
		}
	}

	//Non-admin user is a dev, even if not in Developer group?
	if user.Infrarole == "" {
		user.Infrarole = DEV_ROLE_NAME
	}

	//Values in attributes
	if kuser.Attributes != nil {
		attributes := *kuser.Attributes

		organisation, present := attributes["organisation"]
		if present {
			user.Organisation = organisation[0]
		}
	}

	return &user
}

func mapKeycloakUserWithLabs(kuser *gocloak.User) *models.UserWithLabs {
	user := mapKeycloakUser(kuser)

	userWL := models.UserWithLabs{User: *user}

	//Populate laboratories
	var laboratoryRoles []models.LaboratoryRole

	//For now, role assigned on a lab is the same a user has on the whole infra
	if kuser.Groups != nil {
		kGroups, _ := keycloak.GetUserGroups(kuser)

		for _, kGroup := range kGroups {
			if strings.HasPrefix(*kGroup.Path, "/"+KEYCLOAK_LAB_TOP_GROUP) {
				lab, err := MapLaboratory(*kGroup)

				if err == nil {
					laboratoryRoles = append(laboratoryRoles, models.LaboratoryRole{
						Laboratory: *lab,
						Role:       user.Infrarole,
					})
				}
			}
		}

		if len(laboratoryRoles) > 0 {
			userWL.Laboratories = &laboratoryRoles
		}
	}

	return &userWL
}

func CreateUserKeycloak(user *models.User) (string, error) {
	groups := []string{user.Infrarole}
	attributes := map[string][]string{
		"organisation": {user.Organisation},
	}

	kuser := gocloak.User{
		Username:   &user.Email,
		FirstName:  &user.Firstname,
		LastName:   &user.Lastname,
		Email:      &user.Email,
		Groups:     &groups,
		Attributes: &attributes,
	}

	if user.Disabled != nil {
		kuser.Enabled = gocloak.BoolP(!gocloak.PBool(user.Disabled))
	}

	return keycloak.CreateUser(&kuser)
}

func CreateUserAws(user *models.User) (*scim.User, error) {
	if user.Disabled == nil {
		//TODO Default value on unmarshall ?
		user.Disabled = gocloak.BoolP(false)
	}

	auser := scim.NewUser(user.Firstname, user.Lastname, user.Email, !gocloak.PBool(user.Disabled))

	newuser, err := aws.CreateUser(auser)

	if err == nil {
		//Must obtain group id before adding user
		group, err := aws.GetGroup(user.Infrarole)

		if err == nil {
			err = aws.AddUserToGroup(newuser, group)
		}
	}

	return newuser, err
}

func CreateUserOpenshift(user *models.User) (*userv1.User, error) {
	ouser := userv1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: user.Email,
		},
		FullName: user.Firstname + " " + user.Lastname,
	}

	newOuser, err := openshift.CreateUser(&ouser)

	if err == nil {
		err = openshift.AddUserInGroup(user.Email, user.Infrarole)
	}

	return newOuser, err
}

func UpdateUserKeycloak(userState *UserState, pUser *models.UserUpdate) error {
	kUser := userState.Keycloak

	//change info in keycloak user
	if pUser.Disabled != nil {
		kUser.Enabled = gocloak.BoolP(!gocloak.PBool(pUser.Disabled))
	}

	if pUser.Firstname != nil {
		kUser.FirstName = pUser.Firstname
	}

	if pUser.Lastname != nil {
		kUser.LastName = pUser.Lastname
	}

	if pUser.Organisation != nil {
		(*kUser.Attributes)["organisation"][0] = *pUser.Organisation
	}

	err := keycloak.UpdateUser(kUser)

	if err == nil && pUser.Infrarole != nil && *pUser.Infrarole != userState.Infrarole {
		//Remove former infrarole group
		oldGroup, err := keycloak.GetGroup(userState.Infrarole)

		if err == nil {
			keycloak.DeleteUserFromGroup(kUser, oldGroup)

			//Add to new group
			newGroup, err := keycloak.GetGroup(*pUser.Infrarole)

			if err == nil {
				keycloak.AddUserToGroup(kUser, newGroup)
			}
		}
	}

	return err
}

func UpdateUserAws(userState *UserState, pUser *models.UserUpdate) error {
	auser := userState.Aws

	if pUser.Firstname != nil {
		auser.Name.GivenName = *pUser.Firstname
	}

	if pUser.Lastname != nil {
		auser.Name.FamilyName = *pUser.Lastname
	}

	if pUser.Disabled != nil {
		auser.Active = !*pUser.Disabled
	}

	updatedUser, aerr := aws.UpdateUser(auser)

	if aerr == nil && pUser.Infrarole != nil && *pUser.Infrarole != userState.Infrarole {
		oldGroup, aerr := aws.GetGroup(userState.Infrarole)

		if aerr == nil {
			aerr = aws.RemoveUserFromGroup(updatedUser, oldGroup)

			if aerr == nil {
				newGroup, aerr := aws.GetGroup(*pUser.Infrarole)

				if aerr == nil {
					aerr = aws.AddUserToGroup(updatedUser, newGroup)
				}
			}
		}
	}

	return aerr
}

func UpdateUserOpenshift(userState *UserState, pUser *models.UserUpdate) error {
	var oerr error

	oUser := userState.Openshift

	newFullName := ""

	if pUser.Firstname != nil {
		newFullName = newFullName + *pUser.Firstname
	} else {
		newFullName = newFullName + *userState.Keycloak.FirstName
	}

	newFullName = newFullName + " "

	if pUser.Lastname != nil {
		newFullName = newFullName + *pUser.Lastname
	} else {
		newFullName = newFullName + *userState.Keycloak.LastName
	}

	if newFullName != oUser.FullName {
		oUser.FullName = newFullName

		_, oerr = openshift.UpdateUser(oUser)
	}

	if oerr == nil && pUser.Infrarole != nil && *pUser.Infrarole != userState.Infrarole {
		oerr = openshift.AddUserInGroup(userState.Email, *pUser.Infrarole)

		if oerr == nil {
			oerr = openshift.RemoveUserFromGroup(userState.Email, userState.Infrarole)
		}
	}

	return oerr
}

func DeleteUserKeycloak(userState *UserState) error {
	kuser := userState.Keycloak

	return keycloak.DeleteUser(*kuser.ID)
}

func DeleteUserOpenshift(userState *UserState) error {
	//Groups and users in Openshift are loosely coupled so we have to remove the username from the group.
	err := openshift.RemoveUserFromGroup(userState.Email, userState.Infrarole)

	//TODO handle laboratories groups

	if err == nil {
		err = openshift.DeleteUser(userState.Openshift)
	}

	return err
}

func DeleteUserAws(userState *UserState) error {
	return aws.DeleteUser(userState.Aws)
}

func GetUsers() (*[]models.User, error) {
	//Extract all users
	kusers, err := keycloak.GetUsers()
	if err != nil {
		return nil, NewErrorExternalServerError(err, ERROR_SERVER_KEYCLOAK)
	}

	//getUsers do not provide roles and groups
	//For performance extract all users of the admin group and assume the rest has the user role
	kAdminGroup, err := GetKeycloakAdminGroup()
	if err != nil {
		return nil, NewErrorExternalServerError(err, ERROR_SERVER_KEYCLOAK)
	}

	kadmins, err := keycloak.GetGroupMembers(kAdminGroup)
	//Create a dictionary of admin for easy search
	adminsDict := make(map[string]*gocloak.User, len(kadmins))

	for _, kadmin := range kadmins {

		adminsDict[*kadmin.Username] = kadmin
	}

	//Build user list
	usersList := make([]models.User, 0, len(kusers))

	for _, kuser := range kusers {
		//Add admin role to kuser if he is in the list
		if _, ok := adminsDict[*kuser.Username]; ok {
			reamlRole := []string{ADMIN_ROLE_NAME}
			kuser.RealmRoles = &reamlRole
		}

		usersList = append(usersList, *mapKeycloakUser(kuser))
	}

	return &usersList, nil
}

func GetUserFromUsername(username string) (*models.UserWithLabs, error) {
	kuser, err := keycloak.GetUser(username)

	if err != nil {
		return nil, NewErrorExternalRessourceNotFound(err, ERROR_SERVER_KEYCLOAK)
	}

	//Map User
	return mapKeycloakUserWithLabs(kuser), nil
}

func CreateUser(pUser models.User) error {
	var kerr, oerr, aerr error

	kfunc := func() {
		_, kerr = CreateUserKeycloak(&pUser)
	}

	ofunc := func() {
		_, oerr = CreateUserOpenshift(&pUser)
	}

	afunc := func() {
		_, aerr = CreateUserAws(&pUser)
	}

	Parallelize(kfunc, ofunc, afunc)

	if kerr != nil || oerr != nil || aerr != nil {
		var err error

		if kerr != nil {
			err = NewErrorExternalRessourceExist(kerr, ERROR_SERVER_KEYCLOAK)
		} else if oerr != nil {
			err = NewErrorExternalRessourceExist(oerr, ERROR_SERVER_OPENSHIFT)
		} else if aerr != nil {
			err = NewErrorExternalRessourceExist(aerr, ERROR_SERVER_AWS)
		}

		return err
	}

	//Send account init email
	err := keycloak.ExecuteCurrentActionEmail(pUser.Email)

	if err != nil {
		return NewErrorExternalServerError(err, ERROR_SERVER_KEYCLOAK)
		//TODO email not sent error??
	}

	return nil
}

// Idempotent
func UpdateUser(username string, pUser models.UserUpdate) error {
	userState, err := GetUserState(username)

	if err != nil {
		return err
	}

	var kerr, oerr, aerr error

	kfunc := func() {
		kerr = UpdateUserKeycloak(userState, &pUser)
	}

	ofunc := func() {
		oerr = UpdateUserOpenshift(userState, &pUser)
	}

	afunc := func() {
		aerr = UpdateUserAws(userState, &pUser)
	}

	Parallelize(kfunc, ofunc, afunc)

	if kerr != nil || oerr != nil || aerr != nil {
		var err error

		if kerr != nil {
			err = NewErrorExternalServerError(kerr, ERROR_SERVER_KEYCLOAK)
		} else if oerr != nil {
			err = NewErrorExternalServerError(oerr, ERROR_SERVER_OPENSHIFT)
		} else if aerr != nil {
			err = NewErrorExternalServerError(aerr, ERROR_SERVER_AWS)
		}

		return err
	}

	return nil
}

func DeleteUser(username string) error {
	userState, err := GetUserState(username)

	if err != nil {
		return err
	}

	var kerr, oerr, aerr error

	kfunc := func() {
		kerr = DeleteUserKeycloak(userState)
	}

	ofunc := func() {
		oerr = DeleteUserOpenshift(userState)
	}

	afunc := func() {
		aerr = DeleteUserAws(userState)
	}

	Parallelize(kfunc, ofunc, afunc)

	if kerr != nil || oerr != nil || aerr != nil {
		var err error

		if kerr != nil {
			err = NewErrorExternalServerError(kerr, ERROR_SERVER_KEYCLOAK)
		} else if oerr != nil {
			err = NewErrorExternalServerError(oerr, ERROR_SERVER_OPENSHIFT)
		} else if aerr != nil {
			err = NewErrorExternalServerError(aerr, ERROR_SERVER_AWS)
		}

		return err
	}

	return nil
}
