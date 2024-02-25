#!/bin/bash
set -euo  pipefail

namespace="default"

read -r -p "Enter the name of the secret to delete: " secret_name

if kubectl get secret "$secret_name" -n "$namespace" &> /dev/null; then

    kubectl delete secret "$secret_name" -n "$namespace"
    echo "The secret $secret_name has been deleted."
else
  echo "The secret $secret_name was not found in namespace $namespace."
fi
