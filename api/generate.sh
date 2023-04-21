#! /bin/bash

# Run the openapi generator in a docker container
docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate \
    -i /local/api/openapi.yaml \
    -g go-server \
    -o /local/api-backend \
    --additional-properties outputAsLibrary=true,serverPort=3000

