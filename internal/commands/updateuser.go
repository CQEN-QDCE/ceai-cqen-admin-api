package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/handlers"
	"github.com/spf13/cobra"
)

var UpdateUserCmd = &cobra.Command{

	Use:   "updateuser",
	Short: "Update User",
	Long:  `This command update user from the ceai api`,
	Run: func(cmd *cobra.Command, args []string) {
		Email, _ := cmd.Flags().GetString("email")
		Firstname, _ := cmd.Flags().GetString("firstname")
		Lastname, _ := cmd.Flags().GetString("lastname")
		Organisation, _ := cmd.Flags().GetString("organisation")
		Infrarole, _ := cmd.Flags().GetString("role")
		UpdateUser(Email, Firstname, Lastname, Organisation, Infrarole)
	},
}

func UpdateFlags() {
	UpdateUserCmd.PersistentFlags().StringP("email", "e", "", "The email")
	UpdateUserCmd.PersistentFlags().StringP("firstname", "f", "", "The First name")
	UpdateUserCmd.PersistentFlags().StringP("lastname", "l", "", "The last name")
	UpdateUserCmd.PersistentFlags().StringP("organisation", "o", "", "The name of organisation")
	UpdateUserCmd.PersistentFlags().StringP("role", "r", "", "The infra role Developer/Admin")
}

func init() {
	rootCmd.AddCommand(UpdateUserCmd)
	UpdateFlags()
}

func UpdateUser(email string, firstname string, lastname string, organisation string, infrarole string) {
	if email == "" {
		fmt.Println("Veuillez saisir le courriel")
	} else {
		body := &handlers.UserUpdate{
			Firstname:    &firstname,
			Lastname:     &lastname,
			Organisation: &organisation,
			Infrarole:    &infrarole,
		}

		// Create an HTTP request
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(body)
		url := os.Getenv("SERVER_URL")
		req, _ := http.NewRequest("PUT", url+"/user/"+email, buf)

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
			fmt.Println("l'usager", email, "a été mis à jour")
		} else {
			fmt.Println("L'exécution du traitement a échoué")
		}

		// Print the body to the stdout
		io.Copy(os.Stdout, res.Body)
	}

}
