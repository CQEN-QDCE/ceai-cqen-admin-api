package aws

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"
)

const SSO_CLIENT_TOKEN_TTL = 60

var ssoClient *ssoadmin.Client
var ssoClientTime int64

func GetSsoClient() (*ssoadmin.Client, error) {
	if ssoClient != nil && (time.Now().Unix()-ssoClientTime < SSO_CLIENT_TOKEN_TTL) {
		return ssoClient, nil
	}

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecret := os.Getenv("AWS_SECRET")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     awsAccessKey,
				SecretAccessKey: awsSecret,
			},
		}))

	if err != nil {
		return nil, err
	}

	ssoClient = ssoadmin.NewFromConfig(cfg)
	ssoClientTime = time.Now().Unix()

	return ssoClient, nil
}

func DescribePermissionSet(instanceArn string, permissionSetArn string) (*types.PermissionSet, error) {
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
