package aws

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	organizationstypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	ssoadmintypes "github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"
)

const AWS_REGION = "ca-central-1"
const AWS_CLIENT_TOKEN_TTL = 60

var ssoClient *ssoadmin.Client
var ssoClientTime int64

var orgClient *organizations.Client
var orgClientTime int64

func GetClientConfig() (*aws.Config, error) {
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecret := os.Getenv("AWS_SECRET")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     awsAccessKey,
				SecretAccessKey: awsSecret,
			},
		}))

	cfg.Region = AWS_REGION

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func GetSsoClient() (*ssoadmin.Client, error) {
	if ssoClient != nil && (time.Now().Unix()-ssoClientTime < AWS_CLIENT_TOKEN_TTL) {
		return ssoClient, nil
	}

	cfg, err := GetClientConfig()

	if err != nil {
		return nil, err
	}

	ssoClient = ssoadmin.NewFromConfig(*cfg)
	ssoClientTime = time.Now().Unix()

	return ssoClient, nil
}

func GetOrganizationsClient() (*organizations.Client, error) {
	if orgClient != nil && (time.Now().Unix()-orgClientTime < AWS_CLIENT_TOKEN_TTL) {
		return orgClient, nil
	}

	cfg, err := GetClientConfig()

	if err != nil {
		return nil, err
	}

	orgClient = organizations.NewFromConfig(*cfg)
	orgClientTime = time.Now().Unix()

	return orgClient, nil
}

func DescribePermissionSet(instanceArn string, permissionSetArn string) (*ssoadmintypes.PermissionSet, error) {
	c, err := GetSsoClient()
	if err != nil {
		return nil, err
	}

	permDescription, err := c.DescribePermissionSet(context.TODO(), &ssoadmin.DescribePermissionSetInput{
		InstanceArn:      &instanceArn,
		PermissionSetArn: &permissionSetArn,
	})
	if err != nil {
		return nil, err
	}

	return permDescription.PermissionSet, nil
}

func DescribeAccount(accountId string) (*organizationstypes.Account, error) {
	c, err := GetOrganizationsClient()
	if err != nil {
		return nil, err
	}

	accountOutput, err := c.DescribeAccount(context.TODO(), &organizations.DescribeAccountInput{
		AccountId: &accountId,
	})

	if err != nil {
		return nil, err
	}

	return accountOutput.Account, nil
}
