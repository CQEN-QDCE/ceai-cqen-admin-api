package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/spf13/cobra"
)

// getLabsCmd represents the getLabs command
var getLabsCmd = &cobra.Command{
	Use:   "getlabs",
	Short: "Retourne Info Labs",
	Long:  `Cette commande retourne tout les labs avec l'API du CEAI`,
	Run: func(cmd *cobra.Command, args []string) {
		Format, _ := cmd.Flags().GetString("out")
		GetLabs(Format)
	},
}

func init() {
	rootCmd.AddCommand(getLabsCmd)
	getLabsCmd.PersistentFlags().StringP("out", "o", "none", "Retourne le résultat de la requête selon un format [none, csv, json, jsonpretty]")
}

func GetLabs(format string) {

	// Create an HTTP request
	url := os.Getenv("SERVER_URL")
	res, err := http.Get(url + "/laboratory")

	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		fmt.Println("L'exécution du traitement a échoué")
	} else {

		// Make sure to close after reading
		defer res.Body.Close()

		// read json http response and turn the JSON array into a Go array
		var jsonDataLabs []models.LaboratoryWithResources
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

			err = json.Unmarshal([]byte(jsonDataFromHttp), &jsonDataLabs)

			if err != nil {
				panic(err)
			}

			// Loop over array and print the data of labs
			if format == "csv" {
				fmt.Printf("id,displayname,description,gitrepo\n")
				for _, e := range jsonDataLabs {
					if len(*e.Gitrepo) == 0 {
						*e.Gitrepo = "none"
					}
					fmt.Printf("%v,%v,%v,%v\n", e.Id, e.Displayname, e.Description, *e.Gitrepo)
				}
			} else {
				for _, e := range jsonDataLabs {
					if e.Gitrepo == nil || len(*e.Gitrepo) == 0 {
						e.Gitrepo = new(string)
						*e.Gitrepo = "none"
					}
					if e.Users == nil {
						users := make([]string, 1)
						users[0] = "none"
						e.Users = &users
					}
					fmt.Printf("ID: %v\nDisplayname: %v\nGitrepo: %v\nDescription: %v\nType: %v\nUsers: %v\n\n",
						e.Id,
						e.Displayname,
						*e.Gitrepo,
						e.Description,
						e.Type,
						strings.Join(*e.Users, ", "))
				}
				fmt.Printf("Resultat abrégé, pour avoir toute les informations des labs essayez [-o json] ou [-o jsonpretty]\n")
			}
		}
	}
}
