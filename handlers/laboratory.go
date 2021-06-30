package handlers

import (
	"net/http"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
)

type LaboratoryHandlersInterface interface {

	// (GET /laboratory)
	GetLaboratories(response *apifirst.Response, r *http.Request) error
}
