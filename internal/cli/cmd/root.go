package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/structprinter"
	"github.com/spf13/cobra"
)

//Output format
var outputFormat string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ceai",
	Short:   "Console d'administration CLI du CEAI",
	Long:    `Console d'administration CLI du CEAI`,
	Version: "0.2.0",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		//Validate root persistent Flags here

		//Output Format
		if outputFormat != "text" && outputFormat != "json" && outputFormat != "yaml" {
			return fmt.Errorf("format de sortie invalide: %q. Veuillez choisir 'text', 'json' ou 'yaml'", outputFormat)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "Format de sortie des données ('text', 'json' ou 'yaml')")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	//Search for a session token and a server URL set by the login command

	//TODO persist a list of used server?
}

func PrintOutput(value interface{}) error {
	switch outputFormat {
	case "json":
		return structprinter.PrintJson(value, true)
	case "yaml":
		return structprinter.PrintYaml(value)
	case "text":
		return structprinter.PrintTable(value)
	}

	return fmt.Errorf("format de sortie non supporté: %v", outputFormat)
}
