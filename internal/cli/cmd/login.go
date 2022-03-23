package cmd

import (
	"fmt"
	"time"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var LoginCmd = &cobra.Command{
	Use:   "login [server]",
	Short: "Authentification",
	Long:  `Crée une session au serveur d'API`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		User, _ := cmd.Flags().GetString("user")

		if User == "" {
			prompt := promptui.Prompt{
				Label: "Nom d'usager",
			}

			User, err = prompt.Run()

			if err != nil {
				fmt.Println(err)
				return
			}
		}

		Password, _ := cmd.Flags().GetString("password")

		if Password == "" {
			prompt := promptui.Prompt{
				Label: "Mot de passe",
				Mask:  '*',
			}

			Password, err = prompt.Run()

			if err != nil {
				fmt.Println(err)
				return
			}
		}

		Totp, _ := cmd.Flags().GetString("totp")

		if Totp == "" {
			prompt := promptui.Prompt{
				Label: "Jeton OTP",
			}

			Totp, err = prompt.Run()

			if err != nil {
				fmt.Println(err)
				return
			}
		}

		Login(args[0], User, Password, Totp)
	},
}

func init() {
	LoginCmd.Flags().StringP("user", "u", "", "Nom d'usager")

	LoginCmd.Flags().StringP("password", "p", "", "Mot de passe")

	LoginCmd.Flags().StringP("totp", "t", "", "Jeton OTP")

	rootCmd.AddCommand(LoginCmd)
}

func Login(server string, username string, password string, totp string) {
	//Request time
	requestTime := time.Now().Unix()

	token, err := client.GetKeycloakAccessToken(server, username, password, totp)

	if err != nil {
		fmt.Printf("Error: %v \n", err)
		return
	}

	//TODO check audience
	//Need to decode the token...

	session := client.NewSession(server, requestTime, token)

	err = client.StoreSession(session)

	if err != nil {
		fmt.Printf("Error: %v \n", err)
		return
	}

	fmt.Printf("Logged succesfully as %v on %v \n", username, server)
}
