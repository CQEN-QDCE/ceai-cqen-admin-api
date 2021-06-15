#!/bin/bash

go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen

oapi-codegen --config=config.yaml ../../api/openapi-v1.yaml