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

var UpdateOpenshiftCmd = &cobra.Command{
	Use:   "updateshift",
	Short: "Update Openshift Project",
	Long:  `Cette commande met à jour un projet openshift avec l'API du ceai`,
	Run: func(cmd *cobra.Command, args []string) {
		Id, _ := cmd.Flags().GetString("id")
		Description, _ := cmd.Flags().GetString("description")
		Displayname, _ := cmd.Flags().GetString("displayname")
		IdLab, _ := cmd.Flags().GetString("idlab")
		UpdateOpenshift(Id, Description, Displayname, IdLab)
	},
}

func UpdateOpenshiftFlags() {
	UpdateOpenshiftCmd.PersistentFlags().StringP("id", "i", "", "L'id du projet")
	UpdateOpenshiftCmd.PersistentFlags().StringP("idlab", "l", "", "Le lab associé au projet")
	UpdateOpenshiftCmd.PersistentFlags().StringP("description", "d", "", "La description du projet")
	UpdateOpenshiftCmd.PersistentFlags().StringP("displayname", "n", "", "Le displayname du projet")
}

func init() {
	rootCmd.AddCommand(UpdateOpenshiftCmd)
	UpdateOpenshiftFlags()
}

func UpdateOpenshift(Id string, Description string, Displayname string, IdLab string) {
	if Id == "" {
		fmt.Println("Veuillez spécifier l'ID du lab à modifier avec le flag [-i <id>]")
	} else if Description == "" && Displayname == "" && IdLab == "" {
		fmt.Println("Veuillez spécifier au moins un attribut à modifier du lab (--help pour voir options)")
	} else {
		body := &models.OpenshiftProjectUpdate{
			Description: &Description,
			Displayname: &Displayname,
			IdLab:       &IdLab,
		}

		// Create an HTTP request
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(body)
		url := os.Getenv("SERVER_URL")
		req, _ := http.NewRequest("PUT", url+"/openshift/project/"+Id, buf)

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

		if res.StatusCode == 200 {
			fmt.Println("Le projet", Id, "a été mis à jour")
		} else {
			fmt.Println("L'exécution du traitement a échoué")
		}

		// Print the body to the stdout
		io.Copy(os.Stdout, res.Body)
	}

}
