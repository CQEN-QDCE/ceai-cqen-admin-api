package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/handlers"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
	"github.com/joho/godotenv"
	"github.com/rakyll/statik/fs"
	"gopkg.in/yaml.v2"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"

	_ "github.com/CQEN-QDCE/ceai-cqen-admin-api/api/swaggerui"
)

var Handlers handlers.ServerHandlers

var OpenAPIDoc *openapi3.T

var Config models.Config

func main() {
	log.Println("API d'administration du CEAI v2.0.0")

	err := godotenv.Load()
	if err != nil {
		log.Println("fichier .env non trouvé")
	} else {
		log.Println("fichier .env trouvé, chargement des variables d'environnement")
	}

	log.Println("chargement du fichier de configuration")

	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("la lecture du fichier de configuration a échoué: %v", err.Error())

	}

	log.Printf("chargement du %v complété", Config.Description)

	//Chargement de la spécification OpenAPI
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}

	OpenAPIDoc, err := loader.LoadFromFile(*Config.OpenAPIPath)
	if err != nil {
		log.Fatalf("erreur dans le chargement de la spécification OpenAPI: %v", err.Error())
	}

	if err = OpenAPIDoc.Validate(ctx); err != nil {
		log.Fatalf("fichier de spécification OpenAPI invalide: %v", err.Error())
	}

	//Application de la spécification OpenAPI de sécurité
	var fnAuth openapi3filter.AuthenticationFunc = Authenticate
	fnCallLog := CustomCallLogFunction

	options := &apifirst.RouterOptions{
		AuthenticationFunc: &fnAuth,
		CustomCallLogFunc:  &fnCallLog,
	}

	Handlers.Config = config
	r := apifirst.NewRouter(OpenAPIDoc, Handlers, options)

	//Servir le Swagger UI si désiré
	if os.Getenv("SWAGGER_UI_PATH") != "" {
		swaggerPath := "/" + *Config.SwaggerUIPath

		statikFS, err := fs.New()
		if err != nil {
			panic(err)
		}
		staticServer := http.FileServer(statikFS)

		sh := http.StripPrefix(swaggerPath, staticServer)
		r.Router.PathPrefix(swaggerPath).Handler(sh)

		log.Printf("SwaggerUI disponible au %v/ ", swaggerPath)
	}

	//Healthcheck
	r.Router.HandleFunc("/healthcheck", func(response http.ResponseWriter, request *http.Request) {
		//TODO Effectuer un test de tous les services de la config

		response.WriteHeader(http.StatusOK)
	})

	port := *Config.Port

	log.Fatal(r.Serve(port))
}

func LoadConfig() (*models.Config, error) {

	//Vérifier si un path est spécifié en variable d'environnement
	configFilePath := os.Getenv("CONFIG_FILE")

	//TODO Chercher le fichier si le path est nil

	var config *models.Config

	yamlContent, err := os.ReadFile(configFilePath)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlContent, config)

	if err != nil {
		return nil, err
	}

	//TODO Valider les types de services

	return config, nil
}

func Authenticate(ctx context.Context, authenticationInput *openapi3filter.AuthenticationInput) error {
	headerName := authenticationInput.SecurityScheme.Name
	headerValue := authenticationInput.RequestValidationInput.Request.Header.Get(headerName)

	if headerValue == "" {
		return errors.New(authenticationInput.SecuritySchemeName + " n'est pas spécifié par le gateway")
	}

	switch authenticationInput.SecuritySchemeName {
	case "GatewaySecret":
		//Le GatewaySecret envoyé par le Gateway doit correspondre à celui spécifié dans la configuration
		gatewaySecret := *Config.GatewaySecret

		if gatewaySecret != headerValue {
			return errors.New("le gateway secret spécifié n'est pas valide")
		}
	case "UserRoles":
		//Le UserRoles doit correspondre à la méthode HTTP de la requête: api-read = GET, api-write = *
		switch headerValue {
		case "api-read":
			if authenticationInput.RequestValidationInput.Request.Method != "GET" {
				return errors.New("droits insuffisants pour effectuer cette action")
			}
		case "api-write":
			//Droit à tout pour l'instant
		default:
			return errors.New("role '" + headerValue + "' non reconnu")
		}
	default:
		return errors.New("schéma de sécurité '" + authenticationInput.SecuritySchemeName + "' non reconnu")
	}

	return nil
}

func CustomCallLogFunction(request *http.Request, response *apifirst.ResponseWriter, err error) error {
	output := fmt.Sprintf("%v | %v %v %v", request.Header.Get("X-CEAI-Username"), request.Method, request.RequestURI, response.Status)

	if err != nil {
		output = fmt.Sprintf("%v %v", output, err.Error())

		if e, ok := err.(*openapi3filter.SecurityRequirementsError); ok {
			output = fmt.Sprintf("%v %v", output, e.Errors)
		}
	}

	log.Print(output)

	return nil
}
