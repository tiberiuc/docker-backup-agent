#!/bin/bash

set -exuo pipefail

IMAGE=tiberiuc/backup_agent
TAG=${1:-latest}

docker build -t ${IMAGE}:${TAG} .
docker login
docker push ${IMAGE}:${TAG}
