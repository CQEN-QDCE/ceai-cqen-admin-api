package client

import "github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"

func GetLaboratoryFromId(laboratoryid string) (*models.LaboratoryWithResources, error) {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return nil, err
	}

	pathParam := map[string]string{
		"laboratoryid": laboratoryid,
	}

	resp, err := client.Request("GetLaboratoryFromId", &pathParam, nil)

	if err != nil {
		return nil, err
	}

	var lab models.LaboratoryWithResources
	err = resp.UnmarshalBody(&lab)

	if err != nil {
		return nil, err
	}

	return &lab, nil
}

func GetLaboratories() (*[]models.Laboratory, error) {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return nil, err
	}

	resp, err := client.Request("GetLaboratories", nil, nil)

	if err != nil {
		return nil, err
	}

	var labs []models.Laboratory
	err = resp.UnmarshalBody(&labs)

	if err != nil {
		return nil, err
	}

	return &labs, nil
}
