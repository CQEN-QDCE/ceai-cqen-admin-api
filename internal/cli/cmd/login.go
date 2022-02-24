package cmd

import (
	"fmt"
	"time"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login",
	Long:  `Crée une session au serveur d'API`,
	Run: func(cmd *cobra.Command, args []string) {
		ServerUrl, _ := cmd.Flags().GetString("server")
		User, _ := cmd.Flags().GetString("user")
		Password, _ := cmd.Flags().GetString("password")
		Totp, _ := cmd.Flags().GetString("totp")

		Login(ServerUrl, User, Password, Totp)
	},
}

func init() {
	rootCmd.AddCommand(LoginCmd)
	LoginCmd.Flags().StringP("server", "s", "", "Url du serveur d'API")
	LoginCmd.MarkFlagRequired("server")

	LoginCmd.Flags().StringP("user", "u", "", "Nom d'usager")
	LoginCmd.MarkFlagRequired("user")

	LoginCmd.Flags().StringP("password", "p", "", "Mot de passe")
	LoginCmd.MarkFlagRequired("password")

	LoginCmd.Flags().StringP("totp", "t", "", "Jeton OTP")
	LoginCmd.MarkFlagRequired("totp")
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
