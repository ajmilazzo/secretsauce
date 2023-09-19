#!/bin/bash

echo "post-create start" >> ~/status

# this runs in background after UI is available

# (optional) upgrade packages
sudo apt-get update
sudo apt-get upgrade -y
sudo apt-get autoremove -y
sudo apt-get clean -y

# Install gosec
curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.17.0

# Installing kubebuidler
wget https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH) -P /tmp/
sudo mv /tmp/amd64 /usr/local/bin/kubebuilder 
sudo chmod +x /usr/local/bin/kubebuilder 

# Installing controller-tools
go get sigs.k8s.io/controller-tools/cmd/controller-gen 

echo "post-create complete" >> ~/status