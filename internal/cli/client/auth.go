package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
)

const SESSION_ENV_VAR = "CEAI_API_SESSION"
const SESSION_FILE_PATH = ".ceai"
const SESSION_FILE_NAME = "session"

type Session struct {
	Token                  *models.KeycloakAccessToken
	ServerUrl              string
	RefreshTokenExpireTime int64
}

func NewSession(serverUrl string, requestTime int64, accessToken *models.KeycloakAccessToken) *Session {
	session := Session{
		Token:                  accessToken,
		ServerUrl:              serverUrl,
		RefreshTokenExpireTime: requestTime + int64(accessToken.RefreshExpiresIn),
	}

	return &session
}

func GetSessionFilePath() (*string, error) {
	homedir, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	sessionPath := homedir + string(os.PathSeparator) + SESSION_FILE_PATH

	if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
		err := os.MkdirAll(sessionPath, 0700)

		if err != nil {
			return nil, err
		}
	}

	sessionPath = sessionPath + string(os.PathSeparator) + SESSION_FILE_NAME

	return &sessionPath, nil
}

func GetKeycloakAccessToken(serverUrl string, username string, password, totp string) (*models.KeycloakAccessToken, error) {
	client, err := GetClientToUrl(serverUrl)

	if err != nil {
		return nil, err
	}

	resp, err := client.Request("GetKeycloakAccessToken", nil, models.KeycloakCredentials{
		Username: username,
		Password: password,
		Totp:     totp,
	})

	if err != nil {
		return nil, err
	}

	var token models.KeycloakAccessToken
	err = resp.UnmarshalBody(&token)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func RefreshKeycloakAccessToken(serverUrl string, refreshToken string) (*models.KeycloakAccessToken, error) {
	client, err := GetClientToUrl(serverUrl)

	if err != nil {
		return nil, err
	}

	resp, err := client.Request("RefreshKeycloakAccessToken", nil, refreshToken)

	if err != nil {
		return nil, err
	}

	var token models.KeycloakAccessToken
	err = resp.UnmarshalBody(&token)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func Whoami() (*models.AuthenticatedUser, error) {
	client, err := GetAuthenticatedClient()

	if err != nil {
		return nil, err
	}

	resp, err := client.Request("GetCurrentUserInfo", nil, nil)

	if err != nil {
		return nil, err
	}

	var authUser models.AuthenticatedUser
	err = resp.UnmarshalBody(&authUser)

	if err != nil {
		return nil, err
	}

	return &authUser, nil
}

func StoreSession(session *Session) error {
	jsonSession, err := json.Marshal(session)

	if err != nil {
		return err
	}

	filePath, err := GetSessionFilePath()

	if err != nil {
		return err
	}

	file, err := os.OpenFile(*filePath, os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		return err
	}

	_, err = file.WriteString(base64.StdEncoding.EncodeToString(jsonSession))

	if err != nil {
		return err
	}

	return nil
}

func ReadSession() (*Session, error) {
	filePath, err := GetSessionFilePath()

	if err != nil {
		return nil, err
	}

	envSession, err := os.ReadFile(*filePath)

	if err != nil {
		return nil, err
	}

	jsonSession, err := base64.StdEncoding.DecodeString(string(envSession))

	if err != nil {
		return nil, err
	}

	var session Session
	err = json.Unmarshal([]byte(jsonSession), &session)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func DeleteSession() error {
	filePath, err := GetSessionFilePath()

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	return os.Remove(*filePath)
}

func GetSession() (*Session, error) {
	session, _ := ReadSession()

	if session != nil {
		if session.RefreshTokenExpireTime > time.Now().Unix() {
			requestTime := time.Now().Unix()

			newToken, err := RefreshKeycloakAccessToken(session.ServerUrl, session.Token.RefreshToken)

			if err != nil {
				return nil, err
			}

			newSession := NewSession(session.ServerUrl, requestTime, newToken)

			err = StoreSession(newSession)

			if err != nil {
				return nil, err
			}

			return newSession, nil
		}

		//Delete expired session
		DeleteSession()

		return nil, fmt.Errorf("session expirée")
	}

	return nil, fmt.Errorf("aucune session valide trouvée")
}
