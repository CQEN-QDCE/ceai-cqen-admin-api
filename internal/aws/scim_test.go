package aws

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

	for _, user := range users {
		println(user.Username)
	}
}
