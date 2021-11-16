package cmd

import "github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"

// Replaces empty fields of an OpenshiftProjectWithMeta object with "none"
func ReplaceEmptyFieldsShifts(shift *models.OpenshiftProjectWithMeta) {
	if shift.Displayname == "" {
		shift.Displayname = "none"
	}
	if shift.Description == "" {
		shift.Description = "none"
	}
	if shift.IdLab == "" {
		shift.Displayname = "none"
	}
	if shift.Requester == nil || len(*shift.Requester) == 0 {
		shift.Requester = new(string)
		*shift.Requester = "none"
	}
}
