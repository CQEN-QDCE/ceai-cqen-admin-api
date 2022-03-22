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

func CreateLaboratory(laboratory *models.Laboratory) error {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return err
	}

	_, err = client.Request("CreateLaboratory", nil, laboratory)

	return err
}

func UpdateLaboratory(laboratoryid string, laboratory *models.LaboratoryUpdate) error {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return err
	}

	pathParam := map[string]string{
		"laboratoryid": laboratoryid,
	}

	_, err = client.Request("UpdateLaboratory", &pathParam, laboratory)

	return err
}
