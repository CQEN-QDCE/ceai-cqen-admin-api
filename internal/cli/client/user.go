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

func CreateUser(user *models.User) error {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return err
	}

	_, err = client.Request("CreateUser", nil, user)

	return err
}

func UpdateUser(username string, user *models.UserUpdate) error {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return err
	}

	pathParam := map[string]string{
		"username": username,
	}

	_, err = client.Request("UpdateUser", &pathParam, user)

	return err
}

func DeleteUser(username string) error {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return err
	}

	pathParam := map[string]string{
		"username": username,
	}

	_, err = client.Request("DeleteUser", &pathParam, nil)

	return err
}

func ResetUserCredential(username string, credentialType string) error {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return err
	}

	pathParam := map[string]string{
		"username":       username,
		"credentialType": credentialType,
	}

	_, err = client.Request("ResetUserCredential", &pathParam, nil)

	return err
}

func SendRequiredActionEmail(username string) error {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return err
	}

	pathParam := map[string]string{
		"username": username,
	}

	_, err = client.Request("SendRequiredActionEmail", &pathParam, nil)

	return err
}
