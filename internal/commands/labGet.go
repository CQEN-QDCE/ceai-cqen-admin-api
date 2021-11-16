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

var GetLabCmd = &cobra.Command{
	Use:   "getlab",
	Short: "Retourne info Lab",
	Long:  `Cette commande retourne l'info d'un lab avec l'API du CEAI`,
	Run: func(cmd *cobra.Command, args []string) {
		Id, _ := cmd.Flags().GetString("id")
		Format, _ := cmd.Flags().GetString("out")
		GetLab(Id, Format)
	},
}

func init() {
	rootCmd.AddCommand(GetLabCmd)
	GetLabCmd.PersistentFlags().StringP("id", "i", "", "L'id du lab")
	GetLabCmd.PersistentFlags().StringP("out", "o", "none", "Retourne le résultat de la requête selon un format [none, csv, json, jsonpretty]")
}

func GetLab(Id string, Format string) {

	if Id == "" {
		fmt.Println("Veuillez spécifier l'ID du lab")
	} else {
		// Create an HTTP request
		url := os.Getenv("SERVER_URL")
		req, _ := http.NewRequest("GET", url+"/laboratory/"+Id, nil)

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
			fmt.Println("Le lab", Id, "n'existe pas")
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
				var jsonDataLab models.Laboratory
				err = json.Unmarshal([]byte(body), &jsonDataLab)

				if err != nil {
					panic(err)
				}

				// Loop over array and print the data of labs
				if Format == "csv" {
					fmt.Printf("id,displayname,description,gitrepo\n")
					if jsonDataLab.Gitrepo == nil || len(*jsonDataLab.Gitrepo) == 0 {
						jsonDataLab.Gitrepo = new(string)
						*jsonDataLab.Gitrepo = "none"
					}
					fmt.Printf("%v,%v,%v,%v\n", jsonDataLab.Id, jsonDataLab.Displayname, jsonDataLab.Description, *jsonDataLab.Gitrepo)

				} else {
					if len(*jsonDataLab.Gitrepo) == 0 {
						*jsonDataLab.Gitrepo = "none"
					}
					fmt.Printf("ID: %v\nDisplayname: %v\nGitrepo: %v\nDescription: %v\n",
						jsonDataLab.Id,
						jsonDataLab.Displayname,
						*jsonDataLab.Gitrepo,
						jsonDataLab.Description)
				}
			}
		}

	}

}
