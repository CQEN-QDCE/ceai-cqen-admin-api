package handlers

import (
	"net/http"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
)

type ServerHandlers struct {
	UserHandlersInterface
	LaboratoryHandlersInterface
	OpenshiftHandlersInterface
	AwsHandlersInterface
	KeycloakHandlersInterface
}

func (s ServerHandlers) GetCurrentUserInfo(response *apifirst.ResponseWriter, request *http.Request) error {

	username := request.Header.Get("X-CEAI-Username")
	roles := request.Header.Get("X-CEAI-UserRoles")

	authUser := models.AuthenticatedUser{
		Username: &username,
		Roles:    &roles,
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(authUser)

	return nil
}
