package handlers

import (
	"net/http"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
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
	var err error

	u := User{
		Email:     "user@example.com",
		Firstname: "Bobby",
		Lastname:  "Beaulieu",
		Username:  "bobeau01",
	}

	var t [1]User

	t[0] = u

	response.SetBody(t)

	return err
}

// CreateUser
func (s UserHandlers) CreateUser(response *apifirst.Response, r *http.Request) error {
	var err error

	//TODO Create the user in Keycloak

	response.SetStatus(http.StatusCreated)

	return err
}
