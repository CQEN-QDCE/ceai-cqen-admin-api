package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var _idsRemoveFromLab []string

var RemoveUsersFromLabCmd = &cobra.Command{
	Use:   "rmfromlab",
	Short: "Retire utilisateur(s) d'un lab",
	Long:  `Cette commande retire un ou des utilisateurs d'un lab`,
	Run: func(cmd *cobra.Command, args []string) {
		Id, _ := cmd.Flags().GetString("id")
		AddUsersToLab(Id, _idsRemoveFromLab)
	},
}

func RemoveUsersFromLabFlags() {
	RemoveUsersFromLabCmd.PersistentFlags().StringP("id", "i", "", "L'id du lab")
	// this makes the user to enter multiple values for a flag
	// ref: https://github.com/spf13/cobra/issues/661
	RemoveUsersFromLabCmd.Flags().StringSliceVarP(&_idsRemoveFromLab, "ids", "s", []string{},
		"L'id de l'utilisateur à retirer (peut être répété: [-s <id> -s <id> -s ...] ou [-s <id>,<id>,...])")
}

func init() {
	rootCmd.AddCommand(RemoveUsersFromLabCmd)
	RemoveUsersFromLabFlags()
}

func RemoveUsersFromLab(Id string, Ids []string) {
	if Id == "" {
		fmt.Println("Veuillez spécifier l'id du lab à modifier [-i]")
	} else if len(Ids) == 0 {
		fmt.Println("Veuillez spécifier l'id du ou des membres à retirer [-i <id>]")
	} else {
		// Create an HTTP request
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(Ids)
		url := os.Getenv("SERVER_URL")
		req, _ := http.NewRequest("PUT", url+"/laboratory/"+Id+"/user", buf)

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

		if res.StatusCode == 201 {
			fmt.Println("Le lab", Id, "a été mis à jour")
		} else {
			fmt.Println("L'exécution du traitement a échoué")
		}

		// Print the body to the stdout
		io.Copy(os.Stdout, res.Body)
	}

}
