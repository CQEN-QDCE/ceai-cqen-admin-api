package client

import "github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"

func GetUsers() (*[]models.User, error) {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return nil, err
	}

	resp, err := client.Request("GetUsers", nil, nil)

	if err != nil {
		return nil, err
	}

	var users []models.User
	err = resp.UnmarshalBody(&users)

	if err != nil {
		return nil, err
	}

	return &users, nil
}

func GetUserFromUsername(username string) (*models.UserWithLabs, error) {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return nil, err
	}

	pathParam := map[string]string{
		"username": username,
	}

	resp, err := client.Request("GetUserFromUsername", &pathParam, nil)

	if err != nil {
		return nil, err
	}

	var user models.UserWithLabs
	err = resp.UnmarshalBody(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
