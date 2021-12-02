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

var _ids []string

var AddUsersToLabCmd = &cobra.Command{
	Use:   "addtolab",
	Short: "Add Users to Lab",
	Long:  `This command updates a laboratory from the ceai api to add users to it`,
	Run: func(cmd *cobra.Command, args []string) {
		Id, _ := cmd.Flags().GetString("id")
		AddUsersToLab(Id, _ids)
	},
}

func AddUsersToLabFlags() {
	AddUsersToLabCmd.PersistentFlags().StringP("id", "i", "", "The lab's id")
	// this makes the user to enter multiple values for a flag
	// ref: https://github.com/spf13/cobra/issues/661
	AddUsersToLabCmd.Flags().StringSliceVarP(&_ids, "ids", "s", []string{},
		"The id(s) to add (can be repeated: -s <id> -s <id> -s ... or -s <id>,<id>,...)")
}

func init() {
	rootCmd.AddCommand(AddUsersToLabCmd)
	AddUsersToLabFlags()
}

func AddUsersToLab(Id string, Ids []string) {
	if Id == "" {
		fmt.Println("Please specify the lab's id to modify (-i)")
	} else if len(Ids) == 0 {
		fmt.Println("Please specify the Id of the user to add with flag [-i <id>]")
	} else {
		// Create an HTTP request
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(Ids)
		url := os.Getenv("SERVER_URL")
		req, _ := http.NewRequest("PUT", url+"/laboratory/"+Id+"user", buf)

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
