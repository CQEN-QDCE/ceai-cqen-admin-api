package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/services"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
	"github.com/gorilla/mux"
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

// GetAllUsers
func (s ServerHandlers) GetUsers(response *apifirst.Response, request *http.Request) error {
	usersList, err := services.GetUsers()

	if err != nil {
		if e, ok := err.(services.ErrorExternalServerError); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusInternalServerError)
			return err
		}
		//TODO else
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(usersList)

	return nil
}

func (s ServerHandlers) GetUserFromUsername(response *apifirst.Response, request *http.Request) error {
	params := mux.Vars(request)
	username := params["username"]

	user, err := services.GetUserFromUsername(username)

	if err != nil {
		if e, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusNotFound)
			return err
		}
		//TODO else
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(user)

	return nil
}

// CreateUser
func (s ServerHandlers) CreateUser(response *apifirst.Response, request *http.Request) error {
	puser := models.User{}
	if err := json.NewDecoder(request.Body).Decode(&puser); err != nil {
		response.SetStatus(http.StatusBadRequest)
		log.Println(err)
		return err
	}

	err := services.CreateUser(puser)

	if err != nil {
		if e, ok := err.(services.ErrorExternalRessourceExist); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusConflict)
			return err
		}

		if e, ok := err.(services.ErrorExternalServerError); ok {
			//TODO email not sent...
			log.Println(e.Error())
			response.SetStatus(http.StatusInternalServerError)
			return err
		}
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
	pUser := models.UserUpdate{}
	if err := json.NewDecoder(request.Body).Decode(&pUser); err != nil {
		response.SetStatus(http.StatusBadRequest)
		log.Println(err)
		return err
	}

	err := services.UpdateUser(username, pUser)

	if err != nil {
		if e, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusNotFound)
			return err
		}

		if e, ok := err.(services.ErrorExternalServerError); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusInternalServerError)
			return err
		}
	}

	response.SetStatus(http.StatusOK)

	return nil
}

func (s ServerHandlers) DeleteUser(response *apifirst.Response, request *http.Request) error {
	//Path params
	params := mux.Vars(request)
	username := params["username"]

	err := services.DeleteUser(username)

	if err != nil {
		if e, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusNotFound)
			return err
		}

		if e, ok := err.(services.ErrorExternalServerError); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusInternalServerError)
			return err
		}
	}

	response.SetStatus(http.StatusOK)

	return nil
}
