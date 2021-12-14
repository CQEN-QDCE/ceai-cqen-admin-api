package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/spf13/cobra"
)

// getUsersCmd represents the getUsers command
var GetUsersCmd = &cobra.Command{
	Use:   "getusers",
	Short: "Get Users",
	Long:  `Cette commande retourne tous les utilisateurs avec l'API du CEAI`,
	Run: func(cmd *cobra.Command, args []string) {
		Format, _ := cmd.Flags().GetString("out")
		GetUsers(Format)
	},
}

func init() {
	rootCmd.AddCommand(GetUsersCmd)
	GetUsersCmd.PersistentFlags().StringP("out", "o", "none", "Retourne le résultat de la requête selon un format [none, csv, json, jsonpretty]")
}

func GetUsers(format string) {

	// Create an HTTP request
	url := os.Getenv("SERVER_URL")
	res, err := http.Get(url + "/user")

	if err != nil {
		panic(err)
	}

	// Make sure to close after reading
	defer res.Body.Close()

	// read json http response and turn the JSON array into a Go array
	var jsonDataUser []models.User
	jsonDataFromHttp, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	if format == "json" {

		fmt.Println(string(jsonDataFromHttp))

	} else if format == "jsonpretty" {

		var jsonPretty bytes.Buffer
		err := json.Indent(&jsonPretty, jsonDataFromHttp, "", "\t")

		if err != nil {
			panic(err)
		}

		fmt.Println(jsonPretty.String())

	} else {

		err = json.Unmarshal([]byte(jsonDataFromHttp), &jsonDataUser)

		if err != nil {
			panic(err)
		}

		// Loop over array and print the data of users
		if format == "csv" {
			fmt.Printf("email,role\n")
			for _, e := range jsonDataUser {
				fmt.Printf("%v,%v\n", e.Email, e.Infrarole)
			}
		} else {
			var TABULATION = 55
			for _, e := range jsonDataUser {
				fmt.Printf("Email: %v,%v Role: %v \n",
					e.Email,
					strings.Repeat(" ", TABULATION-(utf8.RuneCountInString("Email: ,"+e.Email))),
					e.Infrarole)
			}
		}
	}

}
