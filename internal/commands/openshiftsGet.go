package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/spf13/cobra"
)

// getLabsCmd represents the getLabs command
var getOpenshiftsCmd = &cobra.Command{
	Use:   "getshifts",
	Short: "Retourne Info sur projets Openshift",
	Long:  `Cette commande retourne tout les projets openshifts avec l'API du CEAI`,
	Run: func(cmd *cobra.Command, args []string) {
		Format, _ := cmd.Flags().GetString("out")
		GetShifts(Format)
	},
}

func init() {
	rootCmd.AddCommand(getOpenshiftsCmd)
	getOpenshiftsCmd.PersistentFlags().StringP("out", "o", "none", "Retourne le résultat de la requête selon un format [none, csv, json, jsonpretty]")
}

func GetShifts(format string) {

	// Create an HTTP request
	url := os.Getenv("SERVER_URL")
	res, err := http.Get(url + "/openshift/project")

	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		fmt.Println("L'exécution du traitement a échoué")
	} else {

		// Make sure to close after reading
		defer res.Body.Close()

		// read json http response and turn the JSON array into a Go array
		var jsonDataLabs []models.OpenshiftProjectWithMeta
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
				fmt.Printf("id,displayname,description,creationDate,idlab,requester\n")
				for _, e := range jsonDataLabs {
					fmt.Printf("%v,%v,%v,%v,%v,%v\n", e.Id, e.Displayname, e.Description, *e.CreationDate, e.IdLab, *e.Requester)
				}
			} else {
				for _, e := range jsonDataLabs {
					fmt.Printf("ID: %v\nDisplayname: %v\nDescription: %v\nCreationDate: %v\nIdLab: %v\nRequester: %v\n\n",
						e.Id,
						e.Displayname,
						e.Description,
						*e.CreationDate,
						e.IdLab,
						*e.Requester,
					)
				}
			}
		}
	}
}
