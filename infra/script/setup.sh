#!/bin/bash
set -euo  pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

QUARKUS_JSON="$SCRIPT_DIR/../keycloak/quarkus-realm.json"

sleep 120

POD_NAME=$(kubectl get pods -n default -l app=keycloak --no-headers -o custom-columns=":metadata.name" | head -1) && 
cat "$QUARKUS_JSON" | kubectl exec -n default -i $POD_NAME -- sh -c "cat > /tmp/quarkus-realm.json" && kubectl exec -n default $POD_NAME -- /opt/keycloak/bin/kc.sh import --file /tmp/quarkus-realm.json

