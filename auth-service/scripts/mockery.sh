#!/bin/bash

set -e
export VERSION=v2.42.2


while IFS="=" read -r package interfaces || [ -n "$package" ]
do
  package=$(echo "$package" | sed 's/\.\///g; s/\//_/g')
  for i in $(echo "$interfaces" | sed 's/[][]//g; s/ //g; s/,/ /g')
  do
    docker run --rm -v "$PWD":/src -w /src vektra/mockery:${VERSION} \
      --output=./internal/mocks \
      --name="${i}" \
      --filename="${package}_${i}".go \
      --structname="${i}_${package}" \
      --dir="${package//_/\/}"
  done
done <<< "$(cat .mockery)"
sudo chown -R $(id -u):$(id -g) ./internal/mocks