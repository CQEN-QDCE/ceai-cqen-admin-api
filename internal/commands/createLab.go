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

var CreateLabCmd = &cobra.Command{
	Use:   "createlab",
	Short: "Create Lab",
	Long:  `This command create a laboratory from the ceai api`,
	Run: func(cmd *cobra.Command, args []string) {
		Id, _ := cmd.Flags().GetString("id")
		Description, _ := cmd.Flags().GetString("description")
		Displayname, _ := cmd.Flags().GetString("displayname")
		Type, _ := cmd.Flags().GetString("type")
		Gitrepo, _ := cmd.Flags().GetString("gitrepo")
		CreateLabs(Id, Description, Displayname, Type, &Gitrepo)
	},
}

func CreateLabFlags() {
	CreateLabCmd.PersistentFlags().StringP("id", "i", "", "The id")
	CreateLabCmd.PersistentFlags().StringP("description", "d", "", "The lab description")
	CreateLabCmd.PersistentFlags().StringP("displayname", "n", "", "The lab displayed name")
	CreateLabCmd.PersistentFlags().StringP("type", "t", "", "The type of lab")
	CreateLabCmd.PersistentFlags().StringP("gitrepo", "g", "", "The lab's gitrepo url (optional)")
}

func init() {
	rootCmd.AddCommand(CreateLabCmd)
	CreateLabFlags()
}

func CreateLabs(Id string, Description string, Displayname string, Type string, Gitrepo *string) {

	if Id == "" {
		fmt.Println("Veuillez saisir l'Id")
	} else if Description == "" {
		fmt.Println("Veuillez saisir la Description")
	} else if Displayname == "" {
		fmt.Println("Veuillez saisir le Displayname")
	} else if Type == "" {
		fmt.Println("Veuillez saisir le type du lab")
	} else {
		body := &models.Laboratory{
			Id:          Id,
			Description: Description,
			Displayname: Displayname,
			Type:        Type,
			Gitrepo:     Gitrepo,
		}

		// Create an HTTP request
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(body)
		url := os.Getenv("SERVER_URL")
		req, _ := http.NewRequest("POST", url+"/laboratory", buf)

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
			fmt.Println("Vous avez bien créé l'usager", Id)
		} else if res.StatusCode == 409 {
			fmt.Println("Le lab existe déjà")
		} else {
			fmt.Println("L'exécution du traitement a échoué")
		}
		// Print the body to the stdout
		io.Copy(os.Stdout, res.Body)
	}

}
