#!/bin/bash
project_root_dir=/root/workspace/openmcp/resource-controller
docker_id="atyx300"
controller_name="openmcp-service-controller"
resource_name="openmcpservice"

org_dir=`pwd`
cd $project_root_dir

export GO111MODULE=on
go mod vendor

go build -o `pwd`/controllers/$resource_name/build/_output/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor resource-controller/controllers/$resource_name/pkg/main && \
docker build -t $docker_id/$controller_name:v0.0.1 controllers/$resource_name/build && \
docker push $docker_id/$controller_name:v0.0.1

cd $org_dir
