package main

import (
	"context"
	"log"
	"os"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/handlers"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
	"github.com/joho/godotenv"

	"github.com/getkin/kin-openapi/openapi3"
)

var Handlers handlers.ServerHandlers

var OpenAPIDoc *openapi3.T

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}

	OpenAPIDoc, err := loader.LoadFromFile("./api/openapi-v1.yaml")
	if err != nil {
		log.Fatal("Error loading OpenAPI Spec file: " + err.Error())
	}

	if err = OpenAPIDoc.Validate(ctx); err != nil {
		log.Fatal("Invalid OpenAPI Spec file: " + err.Error())
	}

	r := apifirst.NewRouter(OpenAPIDoc, Handlers)

	port := os.Getenv("PORT")

	log.Fatal(r.Serve(port))
}
