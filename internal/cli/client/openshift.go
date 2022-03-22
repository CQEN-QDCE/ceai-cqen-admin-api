package client

import "github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"

func GetOpenshiftProjects() (*[]models.OpenshiftProjectWithMeta, error) {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return nil, err
	}

	resp, err := client.Request("GetOpenshiftProjects", nil, nil)

	if err != nil {
		return nil, err
	}

	var projects []models.OpenshiftProjectWithMeta
	err = resp.UnmarshalBody(&projects)

	if err != nil {
		return nil, err
	}

	return &projects, nil
}

func GetOpenshiftProjectFromId(projectid string) (*models.OpenshiftProjectWithMeta, error) {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return nil, err
	}

	pathParam := map[string]string{
		"projectid": projectid,
	}

	resp, err := client.Request("GetOpenshiftProjectFromId", &pathParam, nil)

	if err != nil {
		return nil, err
	}

	var project models.OpenshiftProjectWithMeta
	err = resp.UnmarshalBody(&project)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func CreateOpenshiftProject(project *models.OpenshiftProjectWithLab) error {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return err
	}

	_, err = client.Request("CreateOpenshiftProject", nil, project)

	return err
}

func UpdateOpenshiftProject(projectid string, project *models.OpenshiftProjectUpdate) error {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return err
	}

	pathParam := map[string]string{
		"projectid": projectid,
	}

	_, err = client.Request("UpdateOpenshiftProject", &pathParam, project)

	return err
}
