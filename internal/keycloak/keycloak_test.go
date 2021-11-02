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
