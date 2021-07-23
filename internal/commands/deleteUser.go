package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var DeleteUserCmd = &cobra.Command{
	Use:   "deleteuser",
	Short: "Delete User",
	Long:  `This command delete user from the ceai api`,
	Run: func(cmd *cobra.Command, args []string) {
		Email, _ := cmd.Flags().GetString("email")
		DeleteUser(Email)
	},
}

func init() {
	rootCmd.AddCommand(DeleteUserCmd)
	DeleteUserCmd.PersistentFlags().StringP("email", "e", "", "The email")
}

func DeleteUser(Email string) {

	if Email == "" {
		fmt.Println("Veuillez saisir le courriel")
	} else {
		// Create an HTTP request
		url := os.Getenv("SERVER_URL")
		req, _ := http.NewRequest("DELETE", url+"/user/"+Email, nil)

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
		if res.Status == "200 OK" {
			fmt.Println("Vous avez bien supprimé l'usager", Email)
		} else {
			fmt.Println("L'exécution du traitement a échoué")
		}

		// Print the body to the stdout
		io.Copy(os.Stdout, res.Body)
	}

}
