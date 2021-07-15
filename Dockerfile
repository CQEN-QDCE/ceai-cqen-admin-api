
############################
# STEP 1 build executable binary
############################
FROM registry.redhat.io/rhel8/go-toolset:1.15 AS builder

#Copy source files
ENV GOPATH=/go
WORKDIR $GOPATH/src/github.com/CQEN-QDCE/ceai-cqen-admin-api
COPY . .

# Use go mod
ENV GO111MODULE=on

USER root

#Build Swagger Statik package
RUN go get github.com/rakyll/statik
RUN go install github.com/rakyll/statik

RUN cp ./api/openapi-v1.yaml ./third_party/swaggerui/swagger.yaml
RUN /go/bin/statik -src=./third_party/swaggerui -p=swaggerui -f -dest=./third_party
RUN mv ./third_party/swaggerui/statik.go ./api/swaggerui.go

#Build and install server
RUN go build -o /go/bin/server ./cmd/server

#Copy Openapi spec to bin path
RUN cp ./api/openapi-v1.yaml /go/bin/openapi-v1.yaml

#Copy Openshift service account kubeconfig file to bin path
RUN cp ./openshift/kubeconfig /go/bin/kubeconfig


############################
# STEP 2 build a small image
############################
FROM registry.access.redhat.com/ubi8/ubi:latest

# Copy our static executable.
COPY --from=builder /go/bin /go/bin

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/server

# Define fixed env variables
ENV PORT=8080
ENV OPENAPI_PATH="/go/bin/openapi-v1.yaml"
ENV KUBECONFIG_PATH="/go/bin/kubeconfig"

# Document that the service listens on port 8080.
EXPOSE 8080

