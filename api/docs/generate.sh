#! /bin/bash

# Run the openapi generator in a docker container
docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate \
    -i /local/api/openapi.yaml \
    -g html2 \
    -o /local/api/docs \
    --additional-properties packageName=api,serverPort=3000

