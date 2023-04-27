#! /bin/bash

[ -d ./tmp/ ] || mkdir ./tmp/
echo "Creating tar file..."
tar -czf ./tmp/api-backend.tar.gz ./api-backend/ go.mod go.sum
# Store the pod name in a variable
POD_NAME=$(kubectl --context k3d-podinate-dev --kubeconfig ~/.kube/config --namespace api get pods -l app=api-backend -o custom-columns=NAME:metadata.name --no-headers)
kubectl -n api cp ./tmp/api-backend.tar.gz "$POD_NAME":/tmp/api-backend.tar.gz
# Extract the tar file in the kubernetes pod
kubectl -n api exec "$POD_NAME" -- tar -xvf /tmp/api-backend.tar.gz -C /go/src/github.com/johncave/podinate/
# Remove the tar file from the kubernetes pod
kubectl -n api exec "$POD_NAME" -- rm /tmp/api-backend.tar.gz