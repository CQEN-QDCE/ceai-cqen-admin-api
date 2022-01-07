package handlers

import (
	"net/http"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/services"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
	"github.com/gorilla/mux"
)

type AwsHandlersInterface interface {
	// (GET /aws/account)
	GetAwsAccounts(response *apifirst.Response, request *http.Request) error

	// (GET /aws/account/{accountid})
	GetAwsAccount(response *apifirst.Response, request *http.Request) error
}

func (s ServerHandlers) GetAwsAccounts(response *apifirst.Response, request *http.Request) error {
	accounts, err := services.GetAwsAccounts()

	if err != nil {
		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}
		//TODO else
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(accounts)

	return nil
}

func (s ServerHandlers) GetAwsAccount(response *apifirst.Response, request *http.Request) error {
	params := mux.Vars(request)
	accountId := params["accountid"]

	account, err := services.GetAwsAccount(accountId)

	if err != nil {
		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}
		//TODO else
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(account)

	return nil

}
