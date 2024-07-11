package aws

import (
	"testing"

	scim "github.com/CQEN-QDCE/aws-sso-scim-goclient"
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
		println(user.Username)
		println(user.DisplayName)
		println(user.Name.FamilyName)
		println(user.Name.GivenName)
	}
}

func TestGetUser(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	user, err := GetUser("user.test@example.com")

	if err != nil {
		t.Fatal(err.Error())
	}

	println(user.ID)
	println(user.Username)
	println(user.DisplayName)
	println(user.Name.FamilyName)
	println(user.Name.GivenName)
}

func TestCreateUser(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	user := scim.NewUser(
		"User",
		"Test",
		"user.test@example.com",
		true,
	)

	newuser, err := CreateUser(user)

	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(newuser.ID)
}

func TestCreateGroups(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	grp, err := CreateGroup(scim.NewGroup("Developer"))

	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(grp.DisplayName)
}

func TestUpdateGroup(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	grp, err := GetGroup("Lab_dev2")

	grp.DisplayName = "Lab_dev2b"

	err = UpdateGroup(grp)

	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(grp.DisplayName)
}

func TestAddUserToGroup(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	user, err := GetUser("user.test@example.com")

	if err != nil {
		t.Fatal(err.Error())
	}

	grp, err := GetGroup("Admin")

	if err != nil {
		t.Fatal(err.Error())
	}

	AddUserToGroup(user, grp)

	if err != nil {
		t.Fatal(err.Error())
	}
}
