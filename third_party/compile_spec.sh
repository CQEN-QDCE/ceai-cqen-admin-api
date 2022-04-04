#!/bin/bash
SCRIPT=$(readlink -f $0)
SCRIPTPATH=`dirname $SCRIPT`

go get github.com/rakyll/statik

#SwaggerUI
cp ${SCRIPTPATH}/../api/openapi-v1.yaml ${SCRIPTPATH}/swaggerui/swagger.yaml
statik -src=${SCRIPTPATH}/swaggerui -p=swaggerui -f -dest=${SCRIPTPATH}
mv ${SCRIPTPATH}/swaggerui/statik.go ${SCRIPTPATH}/../api/swaggerui/swaggerui.go

#OAPISpecs
mkdir ${SCRIPTPATH}/oapispecs
cp ${SCRIPTPATH}/../api/openapi-v1.yaml ${SCRIPTPATH}/oapispecs/openapi-v1.yaml
statik -src=${SCRIPTPATH}/oapispecs -p=oapispecs -f -dest=${SCRIPTPATH}
mv ${SCRIPTPATH}/oapispecs/statik.go ${SCRIPTPATH}/../api/oapispecs/openapiv1.go
rm -R ${SCRIPTPATH}/oapispecs