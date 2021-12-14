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

var UpdateLabCmd = &cobra.Command{
	Use:   "updatelab",
	Short: "Update Lab",
	Long:  `Cette commande met à jour un lab à l'aide de l'API du CEAI`,
	Run: func(cmd *cobra.Command, args []string) {
		Id, _ := cmd.Flags().GetString("id")
		Description, _ := cmd.Flags().GetString("description")
		Displayname, _ := cmd.Flags().GetString("displayname")
		Type, _ := cmd.Flags().GetString("type")
		Gitrepo, _ := cmd.Flags().GetString("gitrepo")
		UpdateLab(Id, Description, Displayname, Type, Gitrepo)
	},
}

func UpdateLabFlags() {
	UpdateLabCmd.PersistentFlags().StringP("id", "i", "", "L'id")
	UpdateLabCmd.PersistentFlags().StringP("description", "d", "", "Desciption du lab")
	UpdateLabCmd.PersistentFlags().StringP("displayname", "n", "", "DisplayName du lab")
	UpdateLabCmd.PersistentFlags().StringP("type", "t", "", "Type du lab")
	UpdateLabCmd.PersistentFlags().StringP("gitrepo", "g", "", "L'url du dépôt github (optionnel)")
}

func init() {
	rootCmd.AddCommand(UpdateLabCmd)
	UpdateLabFlags()
}

func UpdateLab(Id string, Description string, Displayname string, Type string, Gitrepo string) {
	if Id == "" {
		fmt.Println("Veuillez spécifier l'ID du lab à modifier avec le flag [-i <id>]")
	} else if Description == "" && Displayname == "" && Type == "" && Gitrepo == "" {
		fmt.Println("Veuillez spécifier au moins un attribut à modifier du lab (--help pour voir options)")
	} else {
		body := &models.LaboratoryUpdate{
			Description: &Description,
			Displayname: &Displayname,
			Type:        &Type,
			Gitrepo:     &Gitrepo,
		}

		// Create an HTTP request
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(body)
		url := os.Getenv("SERVER_URL")
		req, _ := http.NewRequest("PUT", url+"/laboratory/"+Id, buf)

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
			fmt.Println("Le lab", Id, "a été mis à jour")
		} else {
			fmt.Println("L'exécution du traitement a échoué")
		}

		// Print the body to the stdout
		io.Copy(os.Stdout, res.Body)
	}

}
