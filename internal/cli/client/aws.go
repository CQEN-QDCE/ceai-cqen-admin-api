package client

import "github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"

func GetAwsAccounts() (*[]models.AWSAccount, error) {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return nil, err
	}

	resp, err := client.Request("GetAwsAccounts", nil, nil)

	if err != nil {
		return nil, err
	}

	var accounts []models.AWSAccount
	err = resp.UnmarshalBody(&accounts)

	if err != nil {
		return nil, err
	}

	return &accounts, nil
}

func GetAwsAccount(accountid string) (*models.AWSAccount, error) {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return nil, err
	}

	pathParam := map[string]string{
		"accountid": accountid,
	}

	resp, err := client.Request("GetAwsAccount", &pathParam, nil)

	if err != nil {
		return nil, err
	}

	var account models.AWSAccount
	err = resp.UnmarshalBody(&account)

	if err != nil {
		return nil, err
	}

	return &account, nil
}
