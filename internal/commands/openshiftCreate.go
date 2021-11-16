package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/spf13/cobra"
)

var CreateOpenshiftCmd = &cobra.Command{
	Use:   "createshift",
	Short: "Create Openshift Project",
	Long:  `Cette commande crée un projet openshift avec l'API du ceai`,
	Run: func(cmd *cobra.Command, args []string) {
		Id, _ := cmd.Flags().GetString("id")
		Description, _ := cmd.Flags().GetString("description")
		Displayname, _ := cmd.Flags().GetString("displayname")
		IdLab, _ := cmd.Flags().GetString("idlab")
		CreateOpenshift(Id, Description, Displayname, IdLab)
	},
}

func CreateOpenshiftFlags() {
	CreateOpenshiftCmd.PersistentFlags().StringP("id", "i", "", "L'id du projet")
	CreateOpenshiftCmd.PersistentFlags().StringP("idlab", "l", "", "Le lab associé au projet")
	CreateOpenshiftCmd.PersistentFlags().StringP("description", "d", "", "La description du projet")
	CreateOpenshiftCmd.PersistentFlags().StringP("displayname", "n", "", "Le displayname du projet")
}

func init() {
	rootCmd.AddCommand(CreateOpenshiftCmd)
	CreateOpenshiftFlags()
}

func CreateOpenshift(Id string, Description string, Displayname string, IdLab string) {

	if Id == "" {
		fmt.Println("Veuillez saisir l'Id (-i)")
	} else {
		body := &models.OpenshiftProjectWithLab{
			IdLab: IdLab,
			OpenshiftProject: &models.OpenshiftProject{
				Id:          Id,
				Description: Description,
				Displayname: Displayname,
			},
		}

		// Create an HTTP request
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(body)
		url := os.Getenv("SERVER_URL")
		req, _ := http.NewRequest("POST", url+"/openshift/project", buf)

		// Add any defined headers
		req.Header.Set("content-type", "application/json")

		// Create an HTTP client
		client := &http.Client{}

		// Send the request
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		// Make sure to close after reading
		defer res.Body.Close()

		// Display an error or success message
		if res.StatusCode == 201 {
			fmt.Println("Vous avez bien créé le projet openshift", Id)
		} else if res.StatusCode == 409 {
			fmt.Println("Le projet openshift existe déjà")
		} else {
			fmt.Println("L'exécution du traitement a échoué")
		}
		// Print the body to the stdout
		io.Copy(os.Stdout, res.Body)
	}

}
