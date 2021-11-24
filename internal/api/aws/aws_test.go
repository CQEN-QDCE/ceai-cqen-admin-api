package aws

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestDescribePermissionSet(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	instanceArn := os.Getenv("AWS_SSO_INSTANCE_ARN")
	permissionSetArn := os.Getenv("AWS_DEV_PERMISSION_SET_ARN")

	ps, err := DescribePermissionSet(instanceArn, permissionSetArn)

	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println(*ps.Name)
}
