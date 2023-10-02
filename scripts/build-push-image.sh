#!/bin/bash 

# set dir to the directory of the script

cd "$(dirname "$0")"

docker buildx build ./server/Dockerfile -t matesousa/go-k8s-metrics:latest

docker push matesousa/go-k8s-metrics:latest


