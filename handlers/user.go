package handlers

import (
	"log"
	"net/http"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/keycloak"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
	"github.com/gorilla/mux"
)

// Handlers Interface represents all server handlers.
type UserHandlersInterface interface {

	// (GET /user)
	GetAllUsers(response *apifirst.Response, r *http.Request) error
	CreateUser(response *apifirst.Response, r *http.Request) error
}

type UserHandlers struct {
	Handler UserHandlersInterface
}

// User defines model for User.
type User struct {
	Email     string `json:"email,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Username  string `json:"username,omitempty"`
}

// User defines model for User.
type BogusUser struct {
	Email     string `json:"emailx,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Username  string `json:"username,omitempty"`
}

// GetAllUsers
func (s UserHandlers) GetAllUsers(response *apifirst.Response, request *http.Request) error {
	kusers, err := keycloak.GetUsers()
	if err != nil {
		response.SetStatus(http.StatusInternalServerError)
		log.Println(err)
		return err
	}

	t := make([]User, 0, len(kusers))

	for _, kuser := range kusers {
		var user User

		user.Email = *kuser.Email
		user.Firstname = *kuser.FirstName
		user.Lastname = *kuser.LastName
		user.Username = *kuser.Username

		t = append(t, user)
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(t)

	return nil
}

func (s UserHandlers) GetUserByUsername(response *apifirst.Response, request *http.Request) error {
	params := mux.Vars(request)
	username := params["username"]

	kUser, err := keycloak.GetUser(username)
	if err != nil {
		response.SetStatus(http.StatusNotFound)
		return err
	}

	log.Println(kUser)

	var user User

	user.Email = *kUser.Email
	user.Firstname = *kUser.FirstName
	user.Lastname = *kUser.LastName
	user.Username = *kUser.Username

	response.SetStatus(http.StatusOK)
	response.SetBody(user)

	return err
}

// CreateUser
func (s UserHandlers) CreateUser(response *apifirst.Response, r *http.Request) error {
	var err error

	//TODO Create the user in Keycloak

	response.SetStatus(http.StatusCreated)

	return err
}
