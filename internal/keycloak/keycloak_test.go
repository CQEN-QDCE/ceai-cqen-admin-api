package keycloak

import (
	"testing"

	"github.com/Nerzal/gocloak/v8"
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

	for _, user := range users {
		println(*user.Username)
	}
}

func TestGetUser(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	user, err := GetUser("test@example.com")

	if err != nil {
		t.Fatal(err.Error())
	}

	println(*user.ID)
}

func TestGetGroup(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	group, err := GetGroup("Laboratories")

	if err != nil {
		t.Fatal(err.Error())
	}

	println(*group.ID)
}

func TestGetGroups(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	groups, err := GetGroups(gocloak.StringP("Laboratories"))

	if err != nil {
		t.Fatal(err.Error())
	}

	if groups != nil {
		for _, group := range *groups {
			t.Log(group.Name)
		}
	} else {
		t.Log("No group found")
	}

}

func TestAddUserToGroup(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	group, err := GetGroup("dev2")

	if err != nil {
		t.Fatal(err.Error())
	}

	user, err := GetUser("test@example.com")

	if err != nil {
		t.Fatal(err.Error())
	}

	err = AddUserToGroup(user, group)

	if err != nil {
		t.Fatal(err.Error())
	}
}
