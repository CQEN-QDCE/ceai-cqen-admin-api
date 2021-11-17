package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type Lab struct {
	Id          string `json:"id"`
	Displayname string `json:"displayname"`
	Description string `json:"description"`
	Gitrepo     string `json:"gitrepo"`
}

// getLabsCmd represents the getLabs command
var getLabsCmd = &cobra.Command{
	Use:   "getlabs",
	Short: "Get Labs",
	Long:  `This command fetches laboratories from the ceai api`,
	Run: func(cmd *cobra.Command, args []string) {
		Format, _ := cmd.Flags().GetString("out")
		GetLabs(Format)
	},
}

func init() {
	rootCmd.AddCommand(getLabsCmd)
	getLabsCmd.PersistentFlags().StringP("out", "o", "none", "Ouputs result in specified format [none, csv, json, jsonpretty]")
}

func GetLabs(format string) {

	// Create an HTTP request
	url := os.Getenv("SERVER_URL")
	res, err := http.Get(url + "/laboratory")

	if err != nil {
		panic(err)
	}

	// Make sure to close after reading
	defer res.Body.Close()

	// read json http response and turn the JSON array into a Go array
	var jsonDataLabs []Lab
	jsonDataFromHttp, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	if format == "json" {

		fmt.Println(string(jsonDataFromHttp))

	} else if format == "jsonpretty" {

		var jsonPretty bytes.Buffer
		err := json.Indent(&jsonPretty, jsonDataFromHttp, "", "\t")

		if err != nil {
			panic(err)
		}

		fmt.Println(jsonPretty.String())

	} else {

		err = json.Unmarshal([]byte(jsonDataFromHttp), &jsonDataLabs)

		if err != nil {
			panic(err)
		}

		// Loop over array and print the data of labs
		if format == "csv" {
			fmt.Printf("id,displayname,description,gitrepo\n")
			for _, e := range jsonDataLabs {
				if len(e.Gitrepo) == 0 {
					e.Gitrepo = "none"
				}
				fmt.Printf("%v,%v,%v,%v\n", e.Id, e.Displayname, e.Description, e.Gitrepo)
			}
		} else {
			for _, e := range jsonDataLabs {
				if len(e.Gitrepo) == 0 {
					e.Gitrepo = "none"
				}
				fmt.Printf("Displayname: %v, Gitrepo: %v,\nDescription: %v\n\n",
					e.Displayname,
					e.Gitrepo,
					e.Description)
			}
		}
	}

}
