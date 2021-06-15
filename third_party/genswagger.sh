#!/bin/bash

go get github.com/rakyll/statik

cp ../api/openapi-v1.yaml ./swaggerui/swagger.yaml

statik -src=./swaggerui -p=swaggerui -f

mv ./swaggerui/statik.go ../api/swaggerui.go
