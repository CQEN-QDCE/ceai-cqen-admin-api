#!/bin/bash

go get github.com/rakyll/statik

cp ./openapi-v1.yaml ./swaggerui/swagger.yaml

statik -src=./swaggerui -p=swaggerui -f

mv ./swaggerui/statik.go ./swaggerui.go
