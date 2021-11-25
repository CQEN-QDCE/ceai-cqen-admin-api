package services

import (
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/aws"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	organizationstypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
)

const AWS_LAB_GROUP_PREFIX = "Lab_"

func MapAccount(account *organizationstypes.Account) *models.AWSAccount {
	var awsAccount models.AWSAccount

	awsAccount.Id = account.Id
	awsAccount.Email = account.Email
	awsAccount.Name = account.Name

	return &awsAccount
}

func GetAwsAccounts() (*[]models.AWSAccount, error) {
	accountsList, err := aws.ListAccounts()

	if err != nil {
		return nil, NewErrorExternalServerError(err, ERROR_SERVER_AWS)
	}

	var accounts []models.AWSAccount

	for _, account := range *accountsList {
		accounts = append(accounts, *MapAccount(&account))
	}

	return &accounts, nil
}

func GetAwsAccount(accountId string) (*models.AWSAccount, error) {
	accountInfo, err := aws.DescribeAccount(accountId)

	if err != nil {
		return nil, NewErrorExternalRessourceNotFound(err, ERROR_SERVER_AWS)
	}

	return MapAccount(accountInfo), nil
}

func GetAwsLabGroupName(laboratoryId string) string {
	return AWS_LAB_GROUP_PREFIX + laboratoryId
}
