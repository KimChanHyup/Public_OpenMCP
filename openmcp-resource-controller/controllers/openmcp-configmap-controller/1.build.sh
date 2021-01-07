#!/bin/bash
docker_registry_ip="10.0.3.40:5005"
docker_id="openmcp"
controller_name="openmcp-configmap-controller"

export GO111MODULE=on
go mod vendor

go build -o `pwd`/build/_output/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor openmcp/openmcp/openmcp-resource-controller/controllers/openmcp-configmap-controller/pkg/main && \

docker build -t $docker_registry_ip/$docker_id/$controller_name:v0.0.1 build && \
docker push $docker_registry_ip/$docker_id/$controller_name:v0.0.1

