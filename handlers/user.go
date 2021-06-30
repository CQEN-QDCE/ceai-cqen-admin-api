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
	Disabled     bool   `json:"disabled,omitempty"`
	Email        string `json:"email"`
	Firstname    string `json:"firstname"`
	Infrarole    string `json:"infrarole"`
	Lastname     string `json:"lastname"`
	Organisation string `json:"organisation"`
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

func mapKeycloakUser(kuser *gocloak.User) *User {
	var user User

	user.Email = *kuser.Email //email is used as username or id
	user.Firstname = *kuser.FirstName
	user.Lastname = *kuser.LastName
	user.Disabled = !gocloak.PBool(kuser.Enabled)

	//Iterate RealmRoles to fing infraRole (Admin || Developer)
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

func mapKeycloakUserWithLabs(kuser *gocloak.User, kgroups []*gocloak.Group) *UserWithLabs {
	user := mapKeycloakUser(kuser)

	userWL := UserWithLabs{User: *user}

	//Populate laboratories
	var laboratoryRoles []LaboratoryRole

	//For now, role assigned on a lab is the same a user has on the whole infra
	for _, kgroup := range kgroups {
		if strings.HasPrefix(*kgroup.Path, LAB_TOP_GROUP) {
			laboratoryRoles = append(laboratoryRoles, LaboratoryRole{*kgroup.Name, user.Infrarole})
		}
	}

	if len(laboratoryRoles) > 0 {
		userWL.Laboratories = &laboratoryRoles
	}

	return &userWL
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

		adminsDict[*kadmin.Email] = kadmin
	}

	//Build user list
	usersList := make([]User, 0, len(kusers))

	for _, kuser := range kusers {
		//Add admin role to kuser if he is in the list
		if _, ok := adminsDict[*kuser.Email]; ok {
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

	//Get groups because keycloak won't get them in its User endpoint
	kgroups, err := keycloak.GetUserGroups(kuser)
	if err != nil {
		response.SetStatus(http.StatusInternalServerError)
		log.Println(err)
		return err
	}

	//Get roles because keycloak won't get them either
	kroles, err := keycloak.GetUserRoles(kuser)
	if err != nil {
		response.SetStatus(http.StatusInternalServerError)
		log.Println(err)
		return err
	}

	//Add roles to kuser
	krolesList := make([]string, len(kroles))

	for _, krole := range kroles {
		krolesList = append(krolesList, *krole.Name)
	}

	kuser.RealmRoles = &krolesList

	//Map User
	user := mapKeycloakUserWithLabs(kuser, kgroups)

	response.SetStatus(http.StatusOK)
	response.SetBody(user)

	return nil
}

//TODO Usefull ?
func GetUserState(username string) {
	//Check if user exist in Keycloak|AWS|Openshift
	keycloakExist := false
	awsExist := false
	openshiftExist := false

	fKeycloak := func() {
		kuser, _ := keycloak.GetUser(username)

		if kuser != nil {
			keycloakExist = true
		}
	}

	fAws := func() {
		auser, _ := aws.GetUser(username)

		if auser != nil {
			awsExist = true
		}
	}

	fOpenshift := func() {
		ouser, err := openshift.GetUser(username)

		if ouser != nil && err == nil {
			openshiftExist = true
		}
	}

	Parallelize(fKeycloak, fAws, fOpenshift)

	if keycloakExist || awsExist || openshiftExist {
		log.Println("User already exist")
		return
	}

}

// CreateUser
func (s ServerHandlers) CreateUser(response *apifirst.Response, r *http.Request) error {
	//TODO Function?
	puser := User{}
	if err := json.NewDecoder(r.Body).Decode(&puser); err != nil {
		response.SetStatus(http.StatusBadRequest)
		log.Println(err)
		return err
	}

	var kerr, oerr, aerr error

	kfunc := func() {
		groups := []string{puser.Infrarole}
		attributes := map[string][]string{
			"organisation": {puser.Organisation},
		}

		kuser := gocloak.User{
			Username:   &puser.Email,
			FirstName:  &puser.Firstname,
			LastName:   &puser.Lastname,
			Email:      &puser.Email,
			Enabled:    gocloak.BoolP(!puser.Disabled),
			Groups:     &groups,
			Attributes: &attributes,
		}

		kerr = keycloak.CreateUser(&kuser)
	}

	ofunc := func() {
		ouser := userv1.User{
			ObjectMeta: metav1.ObjectMeta{
				Name: puser.Email,
			},
			FullName: puser.Firstname + " " + puser.Lastname,
		}

		_, oerr = openshift.CreateUser(&ouser)

		if oerr == nil {
			oerr = openshift.AddUserInGroup(puser.Email, puser.Infrarole)
		}
	}

	afunc := func() {
		auser := scim.NewUser(puser.Firstname, puser.Lastname, puser.Email, !puser.Disabled)

		newuser, aerr := aws.CreateUser(auser)

		if aerr == nil {
			//Must obtain group id before adding user
			group, err := aws.GetGroup(puser.Infrarole)

			if err == nil {
				aerr = aws.AddUserToGroup(newuser, group)
			} else {
				aerr = err
			}
		}
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

	response.SetStatus(http.StatusCreated)

	return nil
}
