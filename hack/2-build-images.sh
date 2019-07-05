#!/usr/bin/env bash

PACKAGE_NAME=sf-playground

RED='\033[0;31m'
GREEN='\033[0;32m'
INVERTED='\033[7m'
NC='\033[0m' # No Color

set -e #terminate script immediately in case of errors

IMAGE_NAME=${PACKAGE_NAME}-quote:1.0.0

GOOS=linux GOARCH=amd64 go build -o deployment.bin ./cmd/quote

echo -e "${GREEN}Build docker image ${IMAGE_NAME}${NC}"
docker build -t ${IMAGE_NAME} -f Dockerfile-to-deploy .

echo -e "${GREEN}Tag image${NC}"
docker tag ${IMAGE_NAME} mszostok/${IMAGE_NAME}

IMAGE_NAME=${PACKAGE_NAME}-meme:1.0.0

GOOS=linux GOARCH=amd64 go build -o deployment.bin ./cmd/meme

echo -e "${GREEN}Build docker image ${IMAGE_NAME}${NC}"
docker build -t ${IMAGE_NAME} -f deployments/meme/Dockerfile-to-deploy .

echo -e "${GREEN}Tag image${NC}"
docker tag ${IMAGE_NAME} repository.hybris.com:5003/gopher/${IMAGE_NAME}
