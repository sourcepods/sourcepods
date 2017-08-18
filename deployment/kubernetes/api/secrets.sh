#!/bin/bash

kubectl --namespace=gitpods-try create secret generic secrets \
    --from-literal=secret=supersecretsecret
