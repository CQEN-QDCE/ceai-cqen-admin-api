package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var DeleteLabCmd = &cobra.Command{
	Use:   "deletelab",
	Short: "Supprimer Lab",
	Long:  `Cette commande supprime un lab avc l'API du CEAI`,
	Run: func(cmd *cobra.Command, args []string) {
		Id, _ := cmd.Flags().GetString("id")
		DeleteLab(Id)
	},
}

func init() {
	rootCmd.AddCommand(DeleteLabCmd)
	DeleteLabCmd.PersistentFlags().StringP("id", "i", "", "L'ID du lab")
}

func DeleteLab(Id string) {

	if Id == "" {
		fmt.Println("Veuillez spécifier l'id du lab (-i)")
	} else {
		// Create an HTTP request
		url := os.Getenv("SERVER_URL")
		req, _ := http.NewRequest("DELETE", url+"/laboratory/"+Id, nil)

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
		if res.StatusCode == 200 {
			fmt.Println("Vous avez bien supprimé le lab", Id)
		} else {
			fmt.Println("L'exécution du traitement a échoué")
		}

		// Print the body to the stdout
		io.Copy(os.Stdout, res.Body)
	}

}
