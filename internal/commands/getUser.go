package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var GetUserCmd = &cobra.Command{
	Use:   "getuser",
	Short: "Get User",
	Long:  `This command fetches a user from the ceai api`,
	Run: func(cmd *cobra.Command, args []string) {
		Email, _ := cmd.Flags().GetString("email")
		GetUser(Email)
	},
}

func init() {
	rootCmd.AddCommand(GetUserCmd)
	GetUserCmd.PersistentFlags().StringP("email", "e", "", "The email")
}

func GetUser(Email string) {

	if Email == "" {
		fmt.Println("Veuillez saisir le courriel")
	} else {
		// Create an HTTP request
		url := os.Getenv("SERVER_URL")
		req, _ := http.NewRequest("GET", url+"/user/"+Email, nil)

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
			fmt.Println("L'usager", Email, "n'existe pas")
		} else if res.StatusCode != 200 {
			fmt.Println("L'exécution du traitement a échoué")
		}

		// Print the body to the stdout
		io.Copy(os.Stdout, res.Body)
	}

}
