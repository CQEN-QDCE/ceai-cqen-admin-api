package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/spf13/cobra"
)

var GetUserCmd = &cobra.Command{
	Use:   "getuser",
	Short: "Get User",
	Long:  `This command fetches a user from the ceai api`,
	Run: func(cmd *cobra.Command, args []string) {
		Email, _ := cmd.Flags().GetString("email")
		Format, _ := cmd.Flags().GetString("out")
		GetUser(Email, Format)
	},
}

func init() {
	rootCmd.AddCommand(GetUserCmd)
	GetUserCmd.PersistentFlags().StringP("out", "o", "none", "Ouputs result in specified format [none, csv, json, jsonpretty]")
	GetUserCmd.PersistentFlags().StringP("email", "e", "", "The email")
}

func GetUser(Email string, format string) {

	if Email == "" {
		fmt.Println("Veuillez saisir le courriel")
	} else {
		// Create an HTTP request
		url := os.Getenv("SERVER_URL")
		req, _ := http.NewRequest("GET", url+"/user/"+Email, nil)

		// Create an HTTP client
		client := &http.Client{}

		// Send the request
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		// Make sure to close after reading
		defer res.Body.Close()

		if res.StatusCode == 404 {
			fmt.Println("The user", Email, "does not exist")
		} else if res.StatusCode != 200 {
			fmt.Println("The execution has failed")
		} else {
			// read json http response and turn the JSON array into a Go array
			var jsonDataUser models.User
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
					fmt.Printf("%v,%v\n", jsonDataUser.Email, jsonDataUser.Infrarole)
				} else {
					var TABULATION = 55
					fmt.Printf("Email: %v,%v Role: %v \n",
						jsonDataUser.Email,
						strings.Repeat(" ", TABULATION-(utf8.RuneCountInString("Email: ,"+jsonDataUser.Email))),
						jsonDataUser.Infrarole)

				}
			}
		}
	}

}
