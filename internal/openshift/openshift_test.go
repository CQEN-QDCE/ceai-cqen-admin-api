package openshift

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestGetUsers(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	users, err := GetUsers()

	if err != nil {
		t.Fatal(err.Error())
	}

	for _, user := range *users {
		println(user.Name)
	}
}

func TestAddUserInGroup(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	err = AddUserInGroup("francis.gagne@sct.gouv.qc.ca", "Developer")

	if err != nil {
		t.Fatal(err.Error())
	}
}
