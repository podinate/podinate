#! /bin/bash

# Run the openapi generator in a docker container
docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate \
    -i /local/api/openapi.yaml \
    -g go-server \
    -o /local/api-backend \
    --additional-properties outputAsLibrary=true,serverPort=3000

docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate \
    -i /local/api/openapi.yaml \
    -g go \
    -o /local/lib/api_client \
    --additional-properties packageName=api_client

echo "Using sudo to fix user permissions on generated files."

sudo chown $USER:$(id -g) -R ./lib/api_client ./api-backend/go 

# Wipe out the default go.mod and go.sum files
rm ./lib/api_client/go.mod
rm ./lib/api_client/go.sum

echo "Running go fmt on generated code."

# Run fmt on the code to fix errors that the generator creates
go fmt ./lib/api_client/...
go fmt ./api-backend/go/...

# echo "Generating Terraform SDK."

# Generate the Terraform SDK 
# Killing Terraform for now, will build our own package manager
# speakeasy generate sdk --lang terraform -o ./speakeasy/ -s ./api/openapi.yaml
# cd ./speakeasy
# go build -o terraform-provider-podinate