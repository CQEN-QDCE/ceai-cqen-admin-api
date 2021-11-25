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
	Long:  `This command updates a laboratory from the ceai api`,
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
	UpdateLabCmd.PersistentFlags().StringP("id", "i", "", "The id")
	UpdateLabCmd.PersistentFlags().StringP("description", "d", "", "The lab description")
	UpdateLabCmd.PersistentFlags().StringP("displayname", "n", "", "The lab displayed name")
	UpdateLabCmd.PersistentFlags().StringP("type", "t", "", "The type of lab")
	UpdateLabCmd.PersistentFlags().StringP("gitrepo", "g", "", "The lab's gitrepo url (optional)")
}

func init() {
	rootCmd.AddCommand(UpdateLabCmd)
	UpdateLabFlags()
}

func UpdateLab(Id string, Description string, Displayname string, Type string, Gitrepo string) {
	if Id == "" {
		fmt.Println("Please specify the Id of the lab to modify with flag [-i <id>]")
	} else if Description == "" && Displayname == "" && Type == "" && Gitrepo == "" {
		fmt.Println("Please specify at least one attribute to update about the lab (see --help for options)")
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
			fmt.Println("the lab", Id, "has been updated")
		} else {
			fmt.Println("the execution has failed")
		}

		// Print the body to the stdout
		io.Copy(os.Stdout, res.Body)
	}

}
