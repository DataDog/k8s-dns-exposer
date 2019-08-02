#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v1.15.1/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/
curl -Lo kind https://github.com/kubernetes-sigs/kind/releases/download/v0.4.0/kind-linux-amd64 && chmod +x kind && sudo mv kind /usr/local/bin/

kind create cluster