#!/usr/bin/env bash
IMAGE_NAME=quote:1.0.0

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

readonly ROOT_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

pushd ${ROOT_PATH}/../.. > /dev/null

# Exit handler. This function is called anytime an EXIT signal is received.
# This function should never be explicitly called.
function _trap_exit () {
    popd > /dev/null
}
trap _trap_exit EXIT

set -e #terminate script immediately in case of errors

echo -e "${GREEN}Check minikube status:${NC}"
minikubeStat=$(minikube status --format "{{.Kubelet}}_{{.ApiServer}}_{{.Host}}")
if [[ ${minikubeStat} != Running_Running_Running* ]] ;
then
    echo -e "${RED}Got wrong status: ${minikubeStat} ${NC}"
    exit 1
fi

echo -e "${GREEN}Use minikube docker deamon${NC}"
eval $(minikube docker-env --shell bash)

echo -e "${GREEN}Build docker image ${IMAGE_NAME}${NC}"

docker build -t ${IMAGE_NAME} -f Dockerfile-quote .

echo -e "${GREEN}Create deployment and service ${NC}"
kubectl create -f deployments/quote/quote.yaml

echo -e "${GREEN}Get service url ${NC}"
url=$(minikube service quote --url)
full_url="${url}/quote"
echo ${full_url}
