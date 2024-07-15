package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/services"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
	"github.com/gorilla/mux"
)

type UserHandlersInterface interface {

	// (GET /user)
	GetUsers(response *apifirst.ResponseWriter, r *http.Request) error
	// (GET /user/{username})
	GetUserFromUsername(response *apifirst.ResponseWriter, request *http.Request) error
	// (POST /user)
	CreateUser(response *apifirst.ResponseWriter, r *http.Request) error
	// (PUT /user/{username})
	UpdateUser(response *apifirst.ResponseWriter, r *http.Request) error
	// (DELETE /user/{username}/credential/{credentialType})
	ResetUserCredential(response *apifirst.ResponseWriter, request *http.Request) error
	// (POST /user/{username}/actionEmail
	SendRequiredActionEmail(response *apifirst.ResponseWriter, request *http.Request) error
}

// GetAllUsers
func (s ServerHandlers) GetUsers(response *apifirst.ResponseWriter, request *http.Request) error {
	usersList, err := services.GetUsers()

	if err != nil {
		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		response.SetStatus(http.StatusInternalServerError)
		return err
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(usersList)

	return nil
}

func (s ServerHandlers) GetUserFromUsername(response *apifirst.ResponseWriter, request *http.Request) error {
	params := mux.Vars(request)
	username := params["username"]

	user, err := services.GetUserFromUsername(username)

	if err != nil {
		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}

		response.SetStatus(http.StatusInternalServerError)
		return err
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(user)

	return nil
}

// CreateUser
func (s ServerHandlers) CreateUser(response *apifirst.ResponseWriter, request *http.Request) error {
	puser := models.User{}
	if err := json.NewDecoder(request.Body).Decode(&puser); err != nil {
		response.SetStatus(http.StatusBadRequest)
		return err
	}

	err := services.CreateUser(puser)

	if err != nil {
		if _, ok := err.(services.ErrorExternalRessourceExist); ok {
			response.SetStatus(http.StatusConflict)
			return err
		}

		if _, ok := err.(services.ErrorExternalServerError); ok {
			//TODO email not sent...
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		response.SetStatus(http.StatusInternalServerError)
		return err
	}

	response.SetStatus(http.StatusCreated)

	return nil
}

// Idempotent
func (s ServerHandlers) UpdateUser(response *apifirst.ResponseWriter, request *http.Request) error {
	//Path params
	params := mux.Vars(request)
	username := params["username"]

	//Body param
	pUser := models.UserUpdate{}
	if err := json.NewDecoder(request.Body).Decode(&pUser); err != nil {
		response.SetStatus(http.StatusBadRequest)
		return err
	}

	err := services.UpdateUser(username, pUser)

	if err != nil {
		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}

		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		response.SetStatus(http.StatusInternalServerError)
		return err
	}

	response.SetStatus(http.StatusOK)

	return nil
}

func (s ServerHandlers) DeleteUser(response *apifirst.ResponseWriter, request *http.Request) error {
	//Path params
	params := mux.Vars(request)
	username := params["username"]

	err := services.DeleteUser(username)

	if err != nil {
		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}

		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		response.SetStatus(http.StatusInternalServerError)
		return err
	}

	response.SetStatus(http.StatusOK)

	return nil
}

func (s ServerHandlers) ResetUserCredential(response *apifirst.ResponseWriter, request *http.Request) error {
	//Path params
	params := mux.Vars(request)
	username := params["username"]
	credentialType := params["credentialType"]

	credTypeIndex := map[string]string{
		"password": services.CREDENTIAL_PW,
		"otp":      services.CREDENTIAL_OTP,
		"all":      services.CREDENTIAL_ALL,
	}

	err := services.ResetUserCredential(username, credTypeIndex[credentialType])

	if err != nil {
		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}

		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		response.SetStatus(http.StatusInternalServerError)
		return err
	}

	response.SetStatus(http.StatusOK)

	return nil
}

func (s ServerHandlers) SendRequiredActionEmail(response *apifirst.ResponseWriter, request *http.Request) error {
	//Path params
	params := mux.Vars(request)
	username := params["username"]

	err := services.SendCurrentActionEmail(username)

	if err != nil {
		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}

		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		response.SetStatus(http.StatusInternalServerError)
		return err
	}

	response.SetStatus(http.StatusOK)

	return nil
}
