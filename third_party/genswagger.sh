#!/bin/bash
SCRIPT=$(readlink -f $0)
SCRIPTPATH=`dirname $SCRIPT`

go get github.com/rakyll/statik

cp ${SCRIPTPATH}/../api/openapi-v1.yaml ${SCRIPTPATH}/swaggerui/swagger.yaml

statik -src=${SCRIPTPATH}/swaggerui -p=swaggerui -f -dest=${SCRIPTPATH}

mv ${SCRIPTPATH}/swaggerui/statik.go ${SCRIPTPATH}/../api/swaggerui.go
