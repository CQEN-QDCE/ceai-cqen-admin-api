package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/services"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
	"github.com/gorilla/mux"
)

type LaboratoryHandlersInterface interface {

	// (GET /laboratory)
	GetLaboratories(response *apifirst.ResponseWriter, request *http.Request) error

	// (GET /laboratory/{laboratoryid})
	GetLaboratoryFromId(response *apifirst.ResponseWriter, request *http.Request) error

	// (POST /laboratory)
	CreateLaboratory(response *apifirst.ResponseWriter, request *http.Request) error

	// (PUT /laboratory/{laboratoryid})
	UpdateLaboratory(response *apifirst.ResponseWriter, request *http.Request) error

	// (DELETE /laboratory/{laboratoryid}/user)
	RemoveLaboratoryUsers(response *apifirst.ResponseWriter, request *http.Request) error

	// (PUT /laboratory/{laboratoryid}/user)
	AddLaboratoryUsers(response *apifirst.ResponseWriter, request *http.Request) error

	// (PUT /laboratory/{laboratoryid}/openshift/{projectid})
	AttachOpenshiftProjectToLaboratory(response *apifirst.ResponseWriter, request *http.Request) error

	// (DELETE /laboratory/{laboratoryid}/openshift/{projectid})
	DetachOpenshiftProjectFromLaboratory(response *apifirst.ResponseWriter, request *http.Request) error

	// (DELETE /laboratory/{laboratoryid}/aws/{accountid})
	DetachAwsAccountFromLaboratory(response *apifirst.ResponseWriter, request *http.Request) error

	// (PUT /laboratory/{laboratoryid}/aws/{accountid})
	AttachAwsAccountToLaboratory(response *apifirst.ResponseWriter, request *http.Request) error
}

func (s ServerHandlers) GetLaboratories(response *apifirst.ResponseWriter, request *http.Request) error {
	labsList, err := services.GetLaboratories()

	if err != nil {
		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}
		//TODO else
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(labsList)

	return nil
}

func (s ServerHandlers) GetLaboratoryFromId(response *apifirst.ResponseWriter, request *http.Request) error {
	params := mux.Vars(request)
	laboratoryid := params["laboratoryid"]

	lab, err := services.GetLaboratoryFromId(laboratoryid)

	if err != nil {
		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(lab)

	return nil
}

func (s ServerHandlers) CreateLaboratory(response *apifirst.ResponseWriter, request *http.Request) error {
	pLab := models.Laboratory{}
	if err := json.NewDecoder(request.Body).Decode(&pLab); err != nil {
		response.SetStatus(http.StatusBadRequest)
		return err
	}

	err := services.CreateLaboratory(&pLab)

	if err != nil {
		if _, ok := err.(services.ErrorExternalRessourceExist); ok {
			response.SetStatus(http.StatusConflict)
			return err
		}
	}

	response.SetStatus(http.StatusCreated)

	return nil
}

func (s ServerHandlers) UpdateLaboratory(response *apifirst.ResponseWriter, request *http.Request) error {
	//Path params
	params := mux.Vars(request)
	laboratoryId := params["laboratoryid"]

	//Body param
	pLab := models.LaboratoryUpdate{}
	if err := json.NewDecoder(request.Body).Decode(&pLab); err != nil {
		response.SetStatus(http.StatusBadRequest)
		return err
	}

	err := services.UpdateLaboratory(laboratoryId, &pLab)

	if err != nil {
		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}
	}

	response.SetStatus(http.StatusOK)

	return nil
}

func (s ServerHandlers) AddLaboratoryUsers(response *apifirst.ResponseWriter, request *http.Request) error {
	//Path param
	params := mux.Vars(request)
	laboratoryId := params["laboratoryid"]

	//Body param
	var usernameList []string
	if err := json.NewDecoder(request.Body).Decode(&usernameList); err != nil {
		response.SetStatus(http.StatusBadRequest)
		return err
	}

	err := services.AddLaboratoryUsers(laboratoryId, usernameList)

	if err != nil {
		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}
	}

	response.SetStatus(http.StatusCreated)

	return nil

}

func (s ServerHandlers) RemoveLaboratoryUsers(response *apifirst.ResponseWriter, request *http.Request) error {
	//Path param
	params := mux.Vars(request)
	laboratoryId := params["laboratoryid"]

	//Body param
	var usernameList []string
	if err := json.NewDecoder(request.Body).Decode(&usernameList); err != nil {
		response.SetStatus(http.StatusBadRequest)
		return err
	}

	err := services.AddLaboratoryUsers(laboratoryId, usernameList)

	if err != nil {
		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}
	}

	response.SetStatus(http.StatusOK)

	return nil
}

func (s ServerHandlers) AttachOpenshiftProjectToLaboratory(response *apifirst.ResponseWriter, request *http.Request) error {
	//Path param
	params := mux.Vars(request)
	laboratoryId := params["laboratoryid"]
	projectId := params["projectid"]

	err := services.AttachOpenshiftProjectToLaboratory(laboratoryId, projectId)

	if err != nil {
		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}
	}

	response.SetStatus(http.StatusOK)

	return nil
}

func (s ServerHandlers) DetachOpenshiftProjectFromLaboratory(response *apifirst.ResponseWriter, request *http.Request) error {
	//Path param
	params := mux.Vars(request)
	laboratoryId := params["laboratoryid"]
	projectId := params["projectid"]

	err := services.DetachOpenshiftProjectFromLaboratory(laboratoryId, projectId)

	if err != nil {
		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}
	}

	response.SetStatus(http.StatusOK)

	return nil
}

func (s ServerHandlers) AttachAwsAccountToLaboratory(response *apifirst.ResponseWriter, request *http.Request) error {
	params := mux.Vars(request)
	laboratoryId := params["laboratoryid"]
	accountId := params["accountid"]

	err := services.AttachAwsAccountToLaboratory(laboratoryId, accountId)

	if err != nil {
		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}
	}

	response.SetStatus(http.StatusOK)

	return nil
}

func (s ServerHandlers) DetachAwsAccountFromLaboratory(response *apifirst.ResponseWriter, request *http.Request) error {
	params := mux.Vars(request)
	laboratoryId := params["laboratoryid"]
	accountId := params["accountid"]

	err := services.DetachAwsAccountFromLaboratory(laboratoryId, accountId)

	if err != nil {
		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if _, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			response.SetStatus(http.StatusNotFound)
			return err
		}
	}

	response.SetStatus(http.StatusOK)

	return nil
}
