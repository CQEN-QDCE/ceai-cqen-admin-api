package cmd

import (
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Mettre à jour les propriétés d'une ressource",
	Long:  `Mettre à jour les propriétés d'une ressource du CEAI: user, lab, project, account`,
}

func init() {

	rootCmd.AddCommand(updateCmd)

}

func GenerateUpdateFlags(updateStruct interface{}, structUpdateCmd *cobra.Command) error {
	ptrValue := reflect.ValueOf(updateStruct)
	value := ptrValue.Elem()

	for i := 0; i < value.NumField(); i++ {
		//TODO fetch flag info from OpenAPI spec
		switch value.Type().Field(i).Type.Elem().Kind() {
		case reflect.String:
			structUpdateCmd.Flags().String(value.Type().Field(i).Name, "", value.Type().Field(i).Name)
		case reflect.Bool:
			structUpdateCmd.Flags().String(value.Type().Field(i).Name, "", value.Type().Field(i).Name+" 'true' ou 'false'")
		}
	}

	return nil
}

func GetUpdateFlagsValues(updateStruct interface{}, structUpdateCmd *cobra.Command) error {
	ptrValue := reflect.ValueOf(updateStruct)
	value := ptrValue.Elem()

	flagFound := false

	for i := 0; i < value.NumField(); i++ {
		flagVal, err := structUpdateCmd.Flags().GetString(value.Type().Field(i).Name)

		if flagVal != "" && err == nil {
			flagFound = true

			value.Field(i).Set(reflect.New(value.Field(i).Type().Elem()))

			switch value.Type().Field(i).Type.Elem().Kind() {
			case reflect.String:
				value.Field(i).Elem().SetString(flagVal)
			case reflect.Bool:
				var boolval bool
				if flagVal == "true" {
					boolval = true
				} else if flagVal == "false" {
					boolval = false
				} else {
					return fmt.Errorf("les valeurs possibles pour la propritété %s sont 'true' ou 'false'", value.Type().Field(i).Name)
				}

				value.Field(i).Elem().SetBool(boolval)
			}
		}
	}

	if !flagFound {
		return fmt.Errorf("aucune valeur à mettre à jour")
	}

	return nil
}
