#!/bin/bash
set -e 
set -x 

workspace=$(cd $(dirname $0);pwd)
cd $workspace

IMAGE=gosdk-istio-client
TAG=latest

docker build -t $IMAGE:$TAG .
