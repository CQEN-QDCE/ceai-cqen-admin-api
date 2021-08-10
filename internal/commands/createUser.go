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

var CreateUserCmd = &cobra.Command{
	Use:   "createuser",
	Short: "Create User",
	Long:  `This command create user from the ceai api`,
	Run: func(cmd *cobra.Command, args []string) {
		Email, _ := cmd.Flags().GetString("email")
		Firstname, _ := cmd.Flags().GetString("firstname")
		Lastname, _ := cmd.Flags().GetString("lastname")
		Organisation, _ := cmd.Flags().GetString("organisation")
		Infrarole, _ := cmd.Flags().GetString("role")
		CreateUsers(Email, Firstname, Lastname, Organisation, Infrarole)
	},
}

func CreateFlags() {
	CreateUserCmd.PersistentFlags().StringP("email", "e", "", "The email")
	CreateUserCmd.PersistentFlags().StringP("firstname", "f", "", "The First name")
	CreateUserCmd.PersistentFlags().StringP("lastname", "l", "", "The last name")
	CreateUserCmd.PersistentFlags().StringP("organisation", "o", "", "The name of organisation")
	CreateUserCmd.PersistentFlags().StringP("role", "r", "", "The infra role Developer/Admin")
}

func init() {
	rootCmd.AddCommand(CreateUserCmd)
	CreateFlags()
}

func CreateUsers(Email string, Firstname string, Lastname string, Organisation string, Infrarole string) {

	if Email == "" {
		fmt.Println("Veuillez saisir le courriel")
	} else if Firstname == "" {
		fmt.Println("Veuillez saisir le prénom")
	} else if Lastname == "" {
		fmt.Println("Veuillez saisir le nom")
	} else if Organisation == "" {
		fmt.Println("Veuillez saisir l'organistation")
	} else if Infrarole == "" {
		fmt.Println("Veuillez saisir le role")
	} else {
		body := &handlers.User{
			Email:        Email,
			Firstname:    Firstname,
			Lastname:     Lastname,
			Organisation: Organisation,
			Infrarole:    Infrarole,
		}

		// Create an HTTP request
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(body)
		url := os.Getenv("SERVER_URL")
		req, _ := http.NewRequest("POST", url+"/user", buf)

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
			fmt.Println("Vous avez bien créé l'usager", Email)
		} else if res.StatusCode == 409 {
			fmt.Println("L'utilisateur existe déjà")
		} else {
			fmt.Println("L'exécution du traitement a échoué")
		}
		// Print the body to the stdout
		io.Copy(os.Stdout, res.Body)
	}

}
