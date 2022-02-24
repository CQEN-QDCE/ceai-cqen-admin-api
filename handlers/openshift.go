package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/services"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
	"github.com/gorilla/mux"
)

type OpenshiftHandlersInterface interface {
	// (GET /openshift/project)
	GetOpenshiftProjects(response *apifirst.ResponseWriter, request *http.Request) error

	// (POST /openshift/project)
	CreateOpenshiftProject(response *apifirst.ResponseWriter, request *http.Request) error

	// (GET /openshift/project/{projectid})
	GetOpenshiftProjectFromId(response *apifirst.ResponseWriter, request *http.Request) error

	// (PUT /openshift/project/{projectid})
	UpdateOpenshiftProject(response *apifirst.ResponseWriter, request *http.Request) error
}

func (s ServerHandlers) GetOpenshiftProjects(response *apifirst.ResponseWriter, request *http.Request) error {
	projectList, err := services.GetOpenshiftProjects()

	if err != nil {
		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(projectList)

	return nil
}

func (s ServerHandlers) GetOpenshiftProjectFromId(response *apifirst.ResponseWriter, request *http.Request) error {
	//Path param
	params := mux.Vars(request)
	projectId := params["projectid"]

	project, err := services.GetOpenshiftProjectFromId(projectId)

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
	response.SetBody(project)

	return nil
}

func (s ServerHandlers) CreateOpenshiftProject(response *apifirst.ResponseWriter, request *http.Request) error {
	createParam := models.OpenshiftProjectWithLab{}
	if err := json.NewDecoder(request.Body).Decode(&createParam); err != nil {
		response.SetStatus(http.StatusBadRequest)
		return err
	}

	err := services.CreateOpenshiftProject(&createParam)

	if err != nil {
		if _, ok := err.(services.ErrorExternalServerError); ok {
			response.SetStatus(http.StatusInternalServerError)
			return err
		}

		if _, ok := err.(services.ErrorExternalRessourceExist); ok {
			response.SetStatus(http.StatusConflict)
			return err
		}
	}

	response.SetStatus(http.StatusCreated)

	return nil
}

func (s ServerHandlers) UpdateOpenshiftProject(response *apifirst.ResponseWriter, request *http.Request) error {
	params := mux.Vars(request)
	projectId := params["projectid"]

	updateParam := models.OpenshiftProjectUpdate{}
	if err := json.NewDecoder(request.Body).Decode(&updateParam); err != nil {
		response.SetStatus(http.StatusBadRequest)
		return err
	}

	err := services.UpdateOpenshiftProject(projectId, &updateParam)

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
