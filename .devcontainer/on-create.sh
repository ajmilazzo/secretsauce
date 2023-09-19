#!/bin/bash

echo "on-create start" >> ~/status

# create local registry
docker network create k3d
k3d registry create registry.localhost --port 5500
docker network connect k3d k3d-registry.localhost

k3d cluster create --registry-use k3d-registry.localhost:5500 --config .devcontainer/k3d.yaml



echo "on-create complete" >> ~/status
