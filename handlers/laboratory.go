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

type LaboratoryHandlersInterface interface {

	// (GET /laboratory)
	GetLaboratories(response *apifirst.Response, request *http.Request) error

	// (GET /laboratory/{laboratoryid})
	GetLaboratoryFromId(response *apifirst.Response, request *http.Request) error

	// (POST /laboratory)
	CreateLaboratory(response *apifirst.Response, request *http.Request) error

	// (PUT /laboratory/{laboratoryid})
	UpdateLaboratory(response *apifirst.Response, request *http.Request) error

	// (DELETE /laboratory/{laboratoryid}/user)
	RemoveLaboratoryUsers(response *apifirst.Response, request *http.Request) error

	// (PUT /laboratory/{laboratoryid}/user)
	AddLaboratoryUsers(response *apifirst.Response, request *http.Request) error

	// (PUT /laboratory/{laboratoryid}/openshift/{projectid})
	AttachOpenshiftProjectToLaboratory(response *apifirst.Response, request *http.Request) error

	// (DELETE /laboratory/{laboratoryid}/openshift/{projectid})
	DetachOpenshiftProjectFromLaboratory(response *apifirst.Response, request *http.Request) error
}

func (s ServerHandlers) GetLaboratories(response *apifirst.Response, request *http.Request) error {
	labsList, err := services.GetLaboratories()

	if err != nil {
		if e, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusNotFound)
			return err
		}
		//TODO else
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(labsList)

	return nil
}

func (s ServerHandlers) GetLaboratoryFromId(response *apifirst.Response, request *http.Request) error {
	params := mux.Vars(request)
	laboratoryid := params["laboratoryid"]

	lab, err := services.GetLaboratoryFromId(laboratoryid)

	if err != nil {
		if e, ok := err.(services.ErrorExternalServerError); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if e, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusNotFound)
			return err
		}
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(lab)

	return nil
}

func (s ServerHandlers) CreateLaboratory(response *apifirst.Response, request *http.Request) error {
	pLab := models.Laboratory{}
	if err := json.NewDecoder(request.Body).Decode(&pLab); err != nil {
		response.SetStatus(http.StatusBadRequest)
		log.Println(err)
		return err
	}

	err := services.CreateLaboratory(&pLab)

	if err != nil {
		if e, ok := err.(services.ErrorExternalRessourceExist); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusConflict)
			return err
		}
	}

	response.SetStatus(http.StatusCreated)

	return nil
}

func (s ServerHandlers) UpdateLaboratory(response *apifirst.Response, request *http.Request) error {
	//Path params
	params := mux.Vars(request)
	laboratoryId := params["laboratoryid"]

	//Body param
	pLab := models.LaboratoryUpdate{}
	if err := json.NewDecoder(request.Body).Decode(&pLab); err != nil {
		response.SetStatus(http.StatusBadRequest)
		log.Println(err)
		return err
	}

	err := services.UpdateLaboratory(laboratoryId, &pLab)

	if err != nil {
		if e, ok := err.(services.ErrorExternalServerError); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if e, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusNotFound)
			return err
		}
	}

	response.SetStatus(http.StatusOK)

	return nil
}

func (s ServerHandlers) AddLaboratoryUsers(response *apifirst.Response, request *http.Request) error {
	//Path param
	params := mux.Vars(request)
	laboratoryId := params["laboratoryid"]

	//Body param
	var usernameList []string
	if err := json.NewDecoder(request.Body).Decode(&usernameList); err != nil {
		response.SetStatus(http.StatusBadRequest)
		log.Println(err)
		return err
	}

	err := services.AddLaboratoryUsers(laboratoryId, usernameList)

	if err != nil {
		if e, ok := err.(services.ErrorExternalServerError); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if e, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusNotFound)
			return err
		}
	}

	response.SetStatus(http.StatusCreated)

	return nil

}

func (s ServerHandlers) RemoveLaboratoryUsers(response *apifirst.Response, request *http.Request) error {
	//Path param
	params := mux.Vars(request)
	laboratoryId := params["laboratoryid"]

	//Body param
	var usernameList []string
	if err := json.NewDecoder(request.Body).Decode(&usernameList); err != nil {
		response.SetStatus(http.StatusBadRequest)
		log.Println(err)
		return err
	}

	err := services.AddLaboratoryUsers(laboratoryId, usernameList)

	if err != nil {
		if e, ok := err.(services.ErrorExternalServerError); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if e, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusNotFound)
			return err
		}
	}

	response.SetStatus(http.StatusOK)

	return nil
}

func (s ServerHandlers) AttachOpenshiftProjectToLaboratory(response *apifirst.Response, request *http.Request) error {
	//Path param
	params := mux.Vars(request)
	laboratoryId := params["laboratoryid"]
	projectId := params["projectid"]

	err := services.AttachOpenshiftProjectToLaboratory(laboratoryId, projectId)

	if err != nil {
		if e, ok := err.(services.ErrorExternalServerError); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if e, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusNotFound)
			return err
		}
	}

	response.SetStatus(http.StatusOK)

	return nil
}

func (s ServerHandlers) DetachOpenshiftProjectFromLaboratory(response *apifirst.Response, request *http.Request) error {
	//Path param
	params := mux.Vars(request)
	laboratoryId := params["laboratoryid"]
	projectId := params["projectid"]

	err := services.DetachOpenshiftProjectFromLaboratory(laboratoryId, projectId)

	if err != nil {
		if e, ok := err.(services.ErrorExternalServerError); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if e, ok := err.(services.ErrorExternalRessourceNotFound); ok {
			log.Println(e.Error())
			response.SetStatus(http.StatusNotFound)
			return err
		}
	}

	response.SetStatus(http.StatusOK)

	return nil
}
