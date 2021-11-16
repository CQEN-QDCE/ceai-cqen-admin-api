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

	err = AddUserInGroup("test@example.com", "Developer")

	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestUpdateGroup(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	group, err := GetGroup("Lab_dev2")

	if err != nil {
		t.Fatal(err.Error())
	}

	group.SetName("Lab_dev2b")

	_, err = UpdateGroup(group)

	if err != nil {
		t.Fatal(err.Error())
	}
}
