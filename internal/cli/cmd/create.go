package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Créer une ressource",
	Long:  `Créer une ressource du CEAI: user, lab, project, account`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		//Validate root persistent Flags here

		//Input Format
		//TODO

		return nil
	},
}

func init() {
	createCmd.PersistentFlags().BoolP("interactive", "i", false, "Création interactive (experimental)")
	createCmd.PersistentFlags().String("json", "", "Objet JSON")
	createCmd.PersistentFlags().String("yaml", "", "Objet JSON")
	createCmd.PersistentFlags().String("jsonfile", "", "Fichier JSON")
	createCmd.PersistentFlags().String("yamlfile", "", "Fichier YAML")
	rootCmd.AddCommand(createCmd)

}

//Accept Slice interface
func HandleInput(v interface{}) error {
	//Interactive input
	if interactive, _ := createCmd.Flags().GetBool("interactive"); interactive {
		//Extract Struct type from slice pointer
		vPtrValue := reflect.ValueOf(v)
		vPtrType := vPtrValue.Type()

		//struct type =       slice  struct
		vElemType := vPtrType.Elem().Elem()

		sliceValue := vPtrValue.Elem()

		for {
			newValue, err := CreateInstanceFromPrompt(vElemType)

			if err != nil {
				return err
			}

			sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(newValue)))

			//Ask to enter another value
			prompt := promptui.Prompt{
				Label: "Voulez-vous créer un autre élément? (O)ui|(N)on",
			}

			result, err := prompt.Run()

			if err != nil {
				return err
			}

			if strings.ToLower(result) != "o" || strings.ToLower(result) != "oui" {
				break
			}
		}

		return nil
	}

	//Inline inputs
	inlineJson, err := createCmd.Flags().GetString("json")

	if inlineJson != "" && err == nil {
		return json.Unmarshal([]byte(inlineJson), &v)
	}

	inlineYaml, err := createCmd.Flags().GetString("yaml")

	if inlineYaml != "" && err == nil {
		return yaml.Unmarshal([]byte(inlineYaml), &v)
	}

	//File inputs
	jsonFilePath, err := createCmd.Flags().GetString("jsonfile")

	if jsonFilePath != "" && err == nil {
		//read whole file
		jsonContent, err := os.ReadFile(jsonFilePath)

		if err != nil {
			return err
		}

		return json.Unmarshal([]byte(jsonContent), &v)
	}

	yamlFilePath, err := createCmd.Flags().GetString("yamlfile")

	if yamlFilePath != "" && err == nil {
		//read whole file
		yamlContent, err := os.ReadFile(yamlFilePath)

		if err != nil {
			return err
		}

		return yaml.Unmarshal([]byte(yamlContent), &v)
	}

	//No flag specified?
	return fmt.Errorf("Aucun format d'entrée spécifié")
}

func CreateInstanceFromPrompt(t reflect.Type) (interface{}, error) {
	ptr := reflect.New(t)
	value := ptr.Elem()

	for i := 0; i < value.NumField(); i++ {
		switch t.Field(i).Type.Kind() {
		case reflect.String:
			//TODO fetch input info from OpenAPI Spec
			prompt := promptui.Prompt{
				Label: t.Field(i).Name,
			}

			result, err := prompt.Run()

			if err != nil {
				return nil, err
			}

			value.Field(i).SetString(result)

		//Handle "inheritance (recursive)"
		case reflect.Struct:
			structInterface, err := CreateInstanceFromPrompt(t.Field(i).Type)

			if err != nil {
				return nil, err
			}

			value.Field(i).Set(reflect.ValueOf(structInterface))
		}

		/*
			//TODO Handle pointer values
			if t.Field(i).Type.Kind() == reflect.Ptr {
				test := false
				value.Field(i).Set(reflect.ValueOf(&test))
			}
		*/
	}

	return value.Interface(), nil
}
