package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/keycloak"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
)

type KeycloakHandlersInterface interface {

	// (POST /keycloak/token)
	GetKeycloakAccessToken(response *apifirst.ResponseWriter, request *http.Request) error

	// (POST /keycloak/token/refresh)
	RefreshKeycloakAccessToken(response *apifirst.ResponseWriter, request *http.Request) error
}

func (s ServerHandlers) GetKeycloakAccessToken(response *apifirst.ResponseWriter, request *http.Request) error {
	pCreds := models.KeycloakCredentials{}
	if err := json.NewDecoder(request.Body).Decode(&pCreds); err != nil {
		response.SetStatus(http.StatusBadRequest)
		return err
	}

	jwt, err := keycloak.LoginOtp(pCreds.Username, pCreds.Password, pCreds.Totp)

	if err != nil {
		response.SetStatus(http.StatusUnauthorized)
		return err
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(jwt)

	return nil
}

func (s ServerHandlers) RefreshKeycloakAccessToken(response *apifirst.ResponseWriter, request *http.Request) error {
	var pRefreshToken string
	if err := json.NewDecoder(request.Body).Decode(&pRefreshToken); err != nil {
		response.SetStatus(http.StatusBadRequest)
		return err
	}

	jwt, err := keycloak.RefreshToken(pRefreshToken)

	if err != nil {
		response.SetStatus(http.StatusBadRequest)
		return err
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(jwt)

	return nil
}
