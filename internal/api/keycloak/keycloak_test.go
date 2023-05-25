package keycloak

import (
	"errors"
	"testing"

	"github.com/Nerzal/gocloak/v11"
	"github.com/joho/godotenv"
)

func TestGetUsers(t *testing.T) {
	err := godotenv.Load("../../../.env")
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
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	user, err := GetUser("test@example.com")

	if err != nil {
		t.Fatal(err.Error())
	}

	//User should have groups
	if len(*user.Groups) < 1 {
		t.Fatal(errors.New("User has no groups"))
	}

	//User should have RealmRoles
	if len(*user.RealmRoles) < 1 {
		t.Fatal(errors.New("User has no realmRoles"))
	}

	println(*user.ID)
}

func TestGetUserById(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	user, err := GetUserById("user-id-xxxx-xxxx-xxxxxxxxxx")

	if err != nil {
		t.Fatal(err.Error())
	}

	//User should have groups
	if len(*user.Groups) < 1 {
		t.Fatal(errors.New("User has no groups"))
	}

	//User should have RealmRoles
	if len(*user.RealmRoles) < 1 {
		t.Fatal(errors.New("User has no realmRoles"))
	}

	println(*user.ID)
}

func TestGetUserLastLoginEvent(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	user, err := GetUserById("user-id-xxxx-xxxx-xxxxxxxxxx")

	if err != nil {
		t.Fatal(err.Error())
	}

	lastLogin, err := GetUserLastLoginEvent(user)

	if err != nil {
		t.Fatal(err.Error())
	}

	println(lastLogin.Time)
}

func TestGetGroup(t *testing.T) {
	err := godotenv.Load("../../../.env")
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
	err := godotenv.Load("../../../.env")
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
	err := godotenv.Load("../../../.env")
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
