package openshift

import (
	"testing"

	"github.com/joho/godotenv"
	userv1 "github.com/openshift/api/user/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	for _, user := range *users {
		println(user.Name)
	}
}

func TestCreateUser(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	_, err = CreateUser(&userv1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test@example.com",
		},
		FullName: "Bobby Test",
	})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = AddUserInGroup("test@example.com", "Developer")

	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestAddUserInGroup(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	err = AddUserInGroup("test@example.com", "Developer")

	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestUpdateGroup(t *testing.T) {
	err := godotenv.Load("../../../.env")
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
