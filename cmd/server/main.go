package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/handlers"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
	"github.com/joho/godotenv"
	"github.com/rakyll/statik/fs"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
)

var Handlers handlers.ServerHandlers

var OpenAPIDoc *openapi3.T

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	} else {
		log.Println(".env file found. Loading environment variables.")
	}

	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}

	OpenAPIDoc, err := loader.LoadFromFile(os.Getenv("OPENAPI_PATH"))
	if err != nil {
		log.Fatal("Error loading OpenAPI Spec file: " + err.Error())
	}

	if err = OpenAPIDoc.Validate(ctx); err != nil {
		log.Fatal("Invalid OpenAPI Spec file: " + err.Error())
	}

	//API Security validation to support OpenAPI security scheme
	var fnAuth openapi3filter.AuthenticationFunc
	fnAuth = Authenticate

	options := &apifirst.RouterOptions{
		AuthenticationFunc: &fnAuth,
	}

	r := apifirst.NewRouter(OpenAPIDoc, Handlers, options)

	//Serve Swagger UI if wanted
	if os.Getenv("SWAGGER_UI_PATH") != "" {
		swaggerPath := "/" + os.Getenv("SWAGGER_UI_PATH")

		statikFS, err := fs.New()
		if err != nil {
			panic(err)
		}
		staticServer := http.FileServer(statikFS)

		sh := http.StripPrefix(swaggerPath, staticServer)
		r.Router.PathPrefix(swaggerPath).Handler(sh)

		log.Printf("SwaggerUI available at %v/ endpoint", swaggerPath)
	}

	port := os.Getenv("PORT")

	log.Fatal(r.Serve(port))
}

func Authenticate(ctx context.Context, authenticationInput *openapi3filter.AuthenticationInput) error {
	//Assume .env loaded in main or exported
	switch authenticationInput.SecuritySchemeName {
	case "GatewaySecret":
		gatewaySecret := os.Getenv("GATEWAY_SECRET")
		gatewaySecretHeaderName := authenticationInput.SecurityScheme.Name
		gatewaySecretHeaderValue := authenticationInput.RequestValidationInput.Request.Header.Get(gatewaySecretHeaderName)

		if gatewaySecret != gatewaySecretHeaderValue {
			return errors.New("Gateway secrets does not match.")
		}
	case "Username", "UserRoles":
		gatewaySecretHeaderName := authenticationInput.SecurityScheme.Name
		gatewaySecretHeaderValue := authenticationInput.RequestValidationInput.Request.Header.Get(gatewaySecretHeaderName)

		if gatewaySecretHeaderValue == "" {
			return errors.New(authenticationInput.SecuritySchemeName + " not supplied by Gateway.")
		}
	default:
		return errors.New("Unimplemented security scheme.")
	}

	return nil
}
