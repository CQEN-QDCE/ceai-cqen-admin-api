package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	scim "github.com/CQEN-QDCE/aws-sso-scim-goclient"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/aws"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/keycloak"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/openshift"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
	"github.com/Nerzal/gocloak/v8"
	"github.com/gorilla/mux"
	userv1 "github.com/openshift/api/user/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type UserHandlersInterface interface {

	// (GET /user)
	GetUsers(response *apifirst.Response, r *http.Request) error
	// (GET /user/{username})
	GetUserFromUsername(response *apifirst.Response, request *http.Request) error
	// (POST /user)
	CreateUser(response *apifirst.Response, r *http.Request) error
	// (PUT /user/{username})
	UpdateUser(response *apifirst.Response, r *http.Request) error
}

// User defines model for User.
type User struct {
	Disabled     *bool  `json:"disabled,omitempty"`
	Email        string `json:"email"`
	Firstname    string `json:"firstname"`
	Infrarole    string `json:"infrarole"`
	Lastname     string `json:"lastname"`
	Organisation string `json:"organisation"`
}

type UserUpdate struct {
	Disabled     *bool   `json:"disabled,omitempty"`
	Firstname    *string `json:"firstname,omitempty"`
	Infrarole    *string `json:"infrarole,omitempty"`
	Lastname     *string `json:"lastname,omitempty"`
	Organisation *string `json:"organisation,omitempty"`
}

// UserWithLabs defines model for UserWithLabs.
type UserWithLabs struct {
	// Embedded struct due to allOf(#/components/schemas/User)
	User `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Laboratories *[]LaboratoryRole `json:"laboratories,omitempty"`
}

// LaboratoryRole defines model for LaboratoryRole.
type LaboratoryRole struct {
	Laboratory string `json:"laboratory"`
	Role       string `json:"role"`
}

type UserState struct {
	UserWithLabs
	Keycloak  *gocloak.User
	Aws       *scim.User
	Openshift *userv1.User
}

func mapKeycloakUser(kuser *gocloak.User) *User {
	var user User

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
		user.Organisation = attributes["organisation"][0]
	}

	return &user
}

func mapKeycloakUserWithLabs(kuser *gocloak.User) *UserWithLabs {
	user := mapKeycloakUser(kuser)

	userWL := UserWithLabs{User: *user}

	//Populate laboratories
	var laboratoryRoles []LaboratoryRole

	//For now, role assigned on a lab is the same a user has on the whole infra
	if kuser.Groups != nil {
		for _, group := range *kuser.Groups {
			if strings.HasPrefix(group, LAB_TOP_GROUP) {
				laboratoryRoles = append(laboratoryRoles, LaboratoryRole{group, user.Infrarole})
			}
		}

		if len(laboratoryRoles) > 0 {
			userWL.Laboratories = &laboratoryRoles
		}
	}

	return &userWL
}

//Gets current User state across all products: Keycloak|AWS|Openshift
//TODO Errors
func GetUserState(username string) *UserState {
	var state UserState

	fKeycloak := func() {
		state.Keycloak, _ = keycloak.GetUser(username)
	}

	fAws := func() {
		state.Aws, _ = aws.GetUser(username)
	}

	fOpenshift := func() {
		state.Openshift, _ = openshift.GetUser(username)
	}

	Parallelize(fKeycloak, fAws, fOpenshift)

	if state.Keycloak != nil {
		state.UserWithLabs = *mapKeycloakUserWithLabs(state.Keycloak)
	}

	return &state
}

func CreateUserKeycloak(user *User) (string, error) {
	groups := []string{user.Infrarole}
	attributes := map[string][]string{
		"organisation": {user.Organisation},
	}

	kuser := gocloak.User{
		Username:   &user.Email,
		FirstName:  &user.Firstname,
		LastName:   &user.Lastname,
		Email:      &user.Email,
		Enabled:    gocloak.BoolP(!gocloak.PBool(user.Disabled)),
		Groups:     &groups,
		Attributes: &attributes,
	}

	return keycloak.CreateUser(&kuser)
}

func CreateUserAws(user *User) (*scim.User, error) {
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

func CreateUserOpenshift(user *User) (*userv1.User, error) {
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

func UpdateUserKeycloak(userState *UserState, pUser *UserUpdate) error {
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

func UpdateUserAws(userState *UserState, pUser *UserUpdate) error {
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

func UpdateUserOpenshift(userState *UserState, pUser *UserUpdate) error {
	var oerr error

	oUser := userState.Openshift

	fullName := *pUser.Firstname + " " + *pUser.Lastname
	if fullName != oUser.FullName {
		oUser.FullName = fullName

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

// GetAllUsers
func (s ServerHandlers) GetUsers(response *apifirst.Response, request *http.Request) error {
	//Extract all users
	kusers, err := keycloak.GetUsers()
	if err != nil {
		response.SetStatus(http.StatusInternalServerError)
		log.Println(err)
		return err
	}

	//getUsers do not provide roles and groups
	//For performance extract all users of the admin group and assume the rest has the user role
	kAdminGroup, err := GetKeycloakAdminGroup()
	if err != nil {
		response.SetStatus(http.StatusInternalServerError)
		log.Println(err)
		return err
	}

	kadmins, err := keycloak.GetGroupMembers(*kAdminGroup.ID)
	//Create a dictionary of admin for easy search
	adminsDict := make(map[string]*gocloak.User, len(kadmins))

	for _, kadmin := range kadmins {

		adminsDict[*kadmin.Username] = kadmin
	}

	//Build user list
	usersList := make([]User, 0, len(kusers))

	for _, kuser := range kusers {
		//Add admin role to kuser if he is in the list
		if _, ok := adminsDict[*kuser.Username]; ok {
			reamlRole := []string{ADMIN_ROLE_NAME}
			kuser.RealmRoles = &reamlRole
		}

		usersList = append(usersList, *mapKeycloakUser(kuser))
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(usersList)

	return nil
}

func (s ServerHandlers) GetUserFromUsername(response *apifirst.Response, request *http.Request) error {
	params := mux.Vars(request)
	username := params["username"]

	kuser, err := keycloak.GetUser(username)
	if err != nil {
		response.SetStatus(http.StatusNotFound)
		return err
	}

	//Map User
	user := mapKeycloakUserWithLabs(kuser)

	response.SetStatus(http.StatusOK)
	response.SetBody(user)

	return nil
}

// CreateUser
func (s ServerHandlers) CreateUser(response *apifirst.Response, request *http.Request) error {
	puser := User{}
	if err := json.NewDecoder(request.Body).Decode(&puser); err != nil {
		response.SetStatus(http.StatusBadRequest)
		log.Println(err)
		return err
	}

	var kerr, oerr, aerr error

	kfunc := func() {
		_, kerr = CreateUserKeycloak(&puser)
	}

	ofunc := func() {
		_, oerr = CreateUserOpenshift(&puser)
	}

	afunc := func() {
		_, aerr = CreateUserAws(&puser)
	}

	Parallelize(kfunc, ofunc, afunc)

	//TODO Error map
	if kerr != nil {
		log.Println("Keycloak error: " + kerr.Error())
		response.SetStatus(http.StatusConflict)
		return kerr
	}

	if oerr != nil {
		log.Println("Openshift error: " + oerr.Error())
		response.SetStatus(http.StatusConflict)
		return oerr
	}

	if aerr != nil {
		log.Println("AWS error: " + aerr.Error())
		response.SetStatus(http.StatusConflict)
		return aerr
	}

	//Send account init email
	err := keycloak.ExecuteCurrentActionEmail(puser.Email)

	if aerr != nil {
		log.Println(err.Error())
		response.SetStatus(http.StatusInternalServerError)
		return aerr
	}

	response.SetStatus(http.StatusCreated)

	return nil
}

//Idempotent
func (s ServerHandlers) UpdateUser(response *apifirst.Response, request *http.Request) error {
	//Path params
	params := mux.Vars(request)
	username := params["username"]

	//Body param
	pUser := UserUpdate{}
	if err := json.NewDecoder(request.Body).Decode(&pUser); err != nil {
		response.SetStatus(http.StatusBadRequest)
		log.Println(err)
		return err
	}

	userState := GetUserState(username)

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

	//TODO Error map
	if kerr != nil {
		log.Println("Keycloak error: " + kerr.Error())
		response.SetStatus(http.StatusConflict)
		return kerr
	}

	if oerr != nil {
		log.Println("Openshift error: " + oerr.Error())
		response.SetStatus(http.StatusConflict)
		return oerr
	}

	if aerr != nil {
		log.Println("AWS error: " + aerr.Error())
		response.SetStatus(http.StatusConflict)
		return aerr
	}

	response.SetStatus(http.StatusOK)

	return nil
}

func (s ServerHandlers) DeleteUser(response *apifirst.Response, request *http.Request) error {
	//Path params
	params := mux.Vars(request)
	username := params["username"]

	userState := GetUserState(username)

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

	//TODO Error map
	if kerr != nil {
		log.Println("Keycloak error: " + kerr.Error())
		response.SetStatus(http.StatusConflict)
		return kerr
	}

	if oerr != nil {
		log.Println("Openshift error: " + oerr.Error())
		response.SetStatus(http.StatusConflict)
		return oerr
	}

	if aerr != nil {
		log.Println("AWS error: " + aerr.Error())
		response.SetStatus(http.StatusConflict)
		return aerr
	}

	response.SetStatus(http.StatusOK)

	return nil
}
