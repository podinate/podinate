#! /bin/bash

[ -d ./tmp/ ] || mkdir ./tmp/
echo "Creating tar file..."
tar -czf ./tmp/controller.tar.gz ./controller/ go.mod go.sum
# Store the pod name in a variable
POD_NAME=$(kubectl --kubeconfig ~/.kube/config --namespace podinate get pods -l app=podinate-controller -o custom-columns=NAME:metadata.name --no-headers)
kubectl -n podinate cp ./tmp/controller.tar.gz "$POD_NAME":/tmp/controller.tar.gz
echo "Copied to $POD_NAME"
# Extract the tar file in the kubernetes pod
kubectl -n podinate exec "$POD_NAME" -- tar -xvf /tmp/controller.tar.gz -C /go/src/github.com/johncave/podinate/
# Remove the tar file from the kubernetes pod
kubectl -n podinate exec "$POD_NAME" -- rm /tmp/controller.tar.gz