#!/bin/bash

helm upgrade --install gitlab gitlab/gitlab -f ./kubernetes/gitlab/values.dev.yaml --namespace gitlab --create-namespace --timeout 600s