#!/bin/bash

#This script will build and start a container that use a RedHat based image
#Run `podman login registry.redhat.io` command prior running this script the first time 

#Verify if container is built
if [ ! $( docker images -q ceai_admin_api_server ) ];
then
    docker build -t ceai_admin_api_server .
fi


if [ ! "$(docker ps -q -f name=ceai_admin_api_server)" ]; 
then
    if [ "$(docker ps -aq -f status=exited -f name=ceai_admin_api_server)" ]; 
    then
        docker start ceai_admin_api_server
    else
        docker run -d --name ceai_admin_api_server --env-file .env --env OPENAPI_PATH=/go/bin/openapi-v1.yaml --env PORT=8080 -p 8080:8080 ceai_admin_api_server
    fi
else
    docker stop ceai_admin_api_server
    docker start ceai_admin_api_server
fi