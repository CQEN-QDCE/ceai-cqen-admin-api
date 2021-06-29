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
	Email     string `json:"email,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Username  string `json:"username,omitempty"`
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

	url := os.Getenv("SERVER_URL")
	resp, err := http.Get(url + "/" + "user")

	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	// read json http response
	jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var jsonData []User

	err = json.Unmarshal([]byte(jsonDataFromHttp), &jsonData)

	if err != nil {
		panic(err)
	}

	// test struct data
	fmt.Println(jsonData)

}
