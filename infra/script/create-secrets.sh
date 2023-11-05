#!/bin/bash
set -euo  pipefail

namespace="default"

read -r -p "Enter mongo username: " mongo_username
mongo_username=${mongo_username:-"backend-admin"}

read -r -p "Enter mongo password: " mongo_password
mongo_password=${mongo_password:-"djVcLqP3pNzrE"}

read -r -p "Enter postgress username: " postgres_username
postgres_username=${postgres_username:-"backend-admin"}

read -r -p "Enter postgress password: " postgres_password
postgres_password=${postgres_password:-"5r6xso6v"}

read -r -p "Enter keycloak username: " keycloak_username
keycloak_username=${keycloak_username:-"backend-admin"}

read -r -p "Enter keycloak password: " keycloak_password
keycloak_password=${keycloak_password:-"4vcqcbog"}

read -r -p "Enter redis password: " redis_password
redis_password=${redis_password:-"M!4exph3r02p"}

read -r -p "Enter rabbitmq username: " rabbitmq_username
rabbitmq_username=${rabbitmq_username:-"user"}

read -r -p "Enter aeskey: " aes_key_config
aes_key_config=${aes_key_config:-"RS5D@C24yWH&3wk"}

read -r -p "Enter auth-service-keycloak-username: " auth_service_keycloak_username
auth_service_keycloak_username=${auth_service_keycloak_username:-"auth-service"}

read -r -p "Enter auth-service-keycloak-password: " auth_service_keycloak_password
auth_service_keycloak_password=${auth_service_keycloak_password:-"&&u1PG9OE*c8"}

PASSWORD=$(kubectl get secret --namespace default rabbitmq -o jsonpath="{.data.rabbitmq-password}" | base64 -d)

echo "Creating mongo secrets"

if kubectl get secret mongo-username -n "$namespace" &> /dev/null; then
  
  echo "The secret mongo-username exists, so you can delete it"
else
  kubectl -n default create secret generic mongo-username --from-literal=mongo-username="$mongo_username"
fi

if kubectl get secret mongo-password -n "$namespace" &> /dev/null; then
  echo "The secret mongo-password exists, so you can delete it"
else
  kubectl -n default create secret generic mongo-password --from-literal=mongo-password="$mongo_password"
fi

echo "Creating aes key secrets"

if kubectl get secret aeskey -n "$namespace" &> /dev/null; then
  echo "The secret aeskey exists, so you can delete it"
else
  kubectl -n default create secret generic aeskey --from-literal=aes-key-config="$aes_key_config"
fi

echo "Creating router-service secrets"

if kubectl get secret router-service-rabbitmq-username -n "$namespace" &> /dev/null ; then
  echo "The secret router-service-rabbitmq-username exists, so you can delete it"
else
  kubectl -n default create secret generic  router-service-rabbitmq-username --from-literal=rabbitmq-username="$rabbitmq_username"
fi

if kubectl get secret router-service-rabbitmq-password -n "$namespace" &> /dev/null; then
  echo "The router-service-rabbitmq-password exists, so you can delete it"
else
  kubectl -n default create secret generic router-service-rabbitmq-password  --from-literal=rabbitmq-password="${PASSWORD}"
fi

echo "Creating postgres secrets"

if kubectl get secret postgres-username -n "$namespace" &> /dev/null; then
  echo "The secret postgres-username exists, so you can delete it"
else
  kubectl -n default create secret generic postgres-username --from-literal=postgres-username="$postgres_username"
fi

if kubectl get secret postgres-password -n "$namespace" &> /dev/null; then
  echo "The secret postgres-password exists, so you can delete it"
else
  kubectl -n default create secret generic postgres-password --from-literal=postgres-password="$postgres_password"
fi

echo "Creating keycloak secrets"

if kubectl get secret keycloak-username -n "$namespace" &> /dev/null; then
  echo "The secret keycloak-username exists, so you can delete it"
else
  kubectl -n default create secret generic keycloak-username --from-literal=keycloak-username="$keycloak_username"
fi

if kubectl get secret keycloak-password -n "$namespace" &> /dev/null; then
  echo "The secret keycloak-password exists, so you can delete it"
else
  kubectl -n default create secret generic keycloak-password --from-literal=keycloak-password="$keycloak_password"
fi

echo "Creating auth-service secrets"

if kubectl get secret auth-service-keycloak-username -n "$namespace" &> /dev/null; then
  echo "The secret auth-service-keycloak-username exists, so you can delete it"
else
  kubectl -n default create secret generic auth-service-keycloak-username --from-literal=keycloak-username="$auth_service_keycloak_username"
fi

if kubectl get secret auth-service-keycloak-password -n  "$namespace" &> /dev/null; then
  echo "The secret auth-service-keycloak-password exists, so you can delete it"
else
  kubectl -n default create secret generic auth-service-keycloak-password --from-literal=keycloak-password="$auth_service_keycloak_password"
fi

echo "Creating redis secrets"

if kubectl get secret redis -n "$namespace" &> /dev/null; then
  echo "The secret redis exists, so you can delete it"
else
  kubectl -n default create secret generic redis --from-literal=REDIS_PASSWORD="$redis_password"
fi

echo "Creating business-service secrets"

if  kubectl get secret business-service-redis-password -n  "$namespace" &> /dev/null; then
  echo "The secret business-service-redis-password exists, so you can delete it"
else
  kubectl -n default create secret generic business-service-redis-password --from-literal=REDIS_PASSWORD="$redis_password"
fi

if kubectl get secret business-service-rabbitmq-username -n "$namespace" &> /dev/null; then
  echo "The secret business-service-rabbitmq-username exists, so you can delete it"
else
  kubectl -n default create secret generic business-service-rabbitmq-username --from-literal=rabbitmq-username="$rabbitmq_username"
fi

if kubectl get secret business-service-rabbitmq-password -n "$namespace" &> /dev/null; then
  echo "The secret business-service-redis-password exists, so you can delete it"
else
 kubectl -n default create secret generic business-service-rabbitmq-password  --from-literal=rabbitmq-password="${PASSWORD}"
fi

if kubectl get secret business-service-aes-key -n "$namespace" &> /dev/null; then
  echo "The secret business-service-aes-key exists, so you can delete it"
else
  kubectl -n default create secret generic business-service-aes-key --from-literal=aes-key-config="$aes_key_config"
fi

