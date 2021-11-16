package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/spf13/cobra"
)

var GetOpenshiftCmd = &cobra.Command{
	Use:   "getshift",
	Short: "Retourne info Projet Openshift",
	Long:  `Cette commande retourne l'info d'un projet openshift avec l'API du CEAI`,
	Run: func(cmd *cobra.Command, args []string) {
		Id, _ := cmd.Flags().GetString("id")
		Format, _ := cmd.Flags().GetString("out")
		GetOpenshift(Id, Format)
	},
}

func init() {
	rootCmd.AddCommand(GetOpenshiftCmd)
	GetOpenshiftCmd.PersistentFlags().StringP("id", "i", "", "L'id du projet openshift")
	GetOpenshiftCmd.PersistentFlags().StringP("out", "o", "none", "Retourne le résultat de la requête selon un format [none, csv, json, jsonpretty]")
}

func GetOpenshift(Id string, Format string) {

	if Id == "" {
		fmt.Println("Veuillez spécifier l'ID du projet openshift")
	} else {
		// Create an HTTP request
		url := os.Getenv("SERVER_URL")
		req, _ := http.NewRequest("GET", url+"/openshift/project/"+Id, nil)

		// Create an HTTP client
		client := &http.Client{}

		// Send the request
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		// Make sure to close after reading
		defer res.Body.Close()
		// Check for error
		if res.StatusCode == 404 {
			fmt.Println("Le projet openshift", Id, "n'existe pas")
		} else if res.StatusCode != 200 {
			fmt.Println("L'exécution du traitement a échoué")
		} else {
			// OUTPUTTING FORMATS
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}

			if Format == "json" {
				print(string(body))

			} else if Format == "jsonpretty" {
				var jsonPretty bytes.Buffer
				err := json.Indent(&jsonPretty, body, "", "\t")

				if err != nil {
					panic(err)
				}

				fmt.Println(jsonPretty.String())
			} else {
				var jsonDataShift models.OpenshiftProjectWithMeta
				err = json.Unmarshal([]byte(body), &jsonDataShift)

				if err != nil {
					panic(err)
				}

				// Loop over array and print the data of labs
				if Format == "csv" {
					fmt.Printf("id,displayname,description,creationDate,idlab,requester\n")
					ReplaceEmptyFieldsShifts(&jsonDataShift)
					fmt.Printf("%v,%v,%v,%v,%v,%v\n", jsonDataShift.Id, jsonDataShift.Displayname, jsonDataShift.Description,
						jsonDataShift.CreationDate, jsonDataShift.IdLab, *jsonDataShift.Requester)

				} else {
					ReplaceEmptyFieldsShifts(&jsonDataShift)
					fmt.Printf("ID: %v\nDisplayname: %v\nCreation Date: %v\nId Lab: %v\nRequester: %v\nDescription: %v\n",
						jsonDataShift.Id, jsonDataShift.Displayname,
						jsonDataShift.CreationDate, jsonDataShift.IdLab, *jsonDataShift.Requester,
						jsonDataShift.Description)
				}
			}
		}

	}

}
