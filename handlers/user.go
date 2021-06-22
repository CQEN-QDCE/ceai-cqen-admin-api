package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/keycloak"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
	"github.com/Nerzal/gocloak/v8"
	"github.com/gorilla/mux"
)

const LAB_TOP_GROUP = "/Laboratories/"

const ADMIN_ROLE_NAME = "Admin"
const DEV_ROLE_NAME = "Developer"

// Handlers Interface represents all server handlers.
type UserHandlersInterface interface {

	// (GET /user)
	GetUsers(response *apifirst.Response, r *http.Request) error
	CreateUser(response *apifirst.Response, r *http.Request) error
}

type UserHandlers struct {
	Handler UserHandlersInterface
}

// User defines model for User.
type User struct {
	Disabled     bool   `json:"disabled,omitempty"`
	Email        string `json:"email"`
	Firstname    string `json:"firstname"`
	Infrarole    string `json:"infrarole,omitempty"`
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

func GetKeycloakAdminGroup() (*gocloak.Group, error) {
	return keycloak.GetGroup(ADMIN_ROLE_NAME)
}

// GetAllUsers
func (s UserHandlers) GetUsers(response *apifirst.Response, request *http.Request) error {
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

func (s UserHandlers) GetUserFromUsername(response *apifirst.Response, request *http.Request) error {
	params := mux.Vars(request)
	username := params["username"]

	kuser, err := keycloak.GetUser(username)
	if err != nil {
		response.SetStatus(http.StatusNotFound)
		return err
	}

	//Get groups
	kgroups, err := keycloak.GetUserGroups(kuser)
	if err != nil {
		response.SetStatus(http.StatusInternalServerError)
		log.Println(err)
		return err
	}

	//Get Roles, because keycloak won't get them in its User endpoint
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

// CreateUser
func (s UserHandlers) CreateUser(response *apifirst.Response, r *http.Request) error {
	var err error

	//TODO Create the user in Keycloak

	response.SetStatus(http.StatusCreated)

	return err
}
