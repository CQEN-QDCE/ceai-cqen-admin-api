package aws

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestDescribePermissionSet(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	err = godotenv.Load("../../test.env")
	if err != nil {
		t.Fatal("Error loading test.env file (This test requires test values): " + err.Error())
	}

	instanceArn := os.Getenv("AWS_INSTANCE_ARN")
	permissionSetArn := os.Getenv("AWS_PERMISSION_SET_ARN")

	ps, err := DescribePermissionSet(instanceArn, permissionSetArn)

	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println(*ps.Name)
}
