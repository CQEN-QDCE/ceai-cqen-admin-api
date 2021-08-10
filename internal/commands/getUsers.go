package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type User struct {
	Disabled     *bool  `json:"disabled,omitempty"`
	Email        string `json:"email"`
	Firstname    string `json:"firstname"`
	Infrarole    string `json:"infrarole"`
	Lastname     string `json:"lastname"`
	Organisation string `json:"organisation"`
}

// getUsersCmd represents the getUsers command
var getUsersCmd = &cobra.Command{
	Use:   "getusers",
	Short: "Get Users",
	Long:  `This command fetches users from the ceai api`,
	Run: func(cmd *cobra.Command, args []string) {
		GetUsers()
	},
}

func init() {
	rootCmd.AddCommand(getUsersCmd)
}

func GetUsers() {

	// Create an HTTP request
	url := os.Getenv("SERVER_URL")
	res, err := http.Get(url + "/user")

	if err != nil {
		panic(err)
	}

	// Make sure to close after reading
	defer res.Body.Close()

	// read json http response and turn the JSON array into a Go array
	var jsonDataUser []User
	jsonDataFromHttp, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(jsonDataFromHttp), &jsonDataUser)

	if err != nil {
		panic(err)
	}

	// Loop over array and print the data of users
	for _, e := range jsonDataUser {
		fmt.Printf("Email: %v, Role: %v \n", e.Email, e.Infrarole)
	}

}
