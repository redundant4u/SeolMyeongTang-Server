#!/bin/sh

if [ $# -ne 1 ]; then
	echo 'No argument'
	exit -1
fi

### Pull Docker Image From ECR
ECR_URL=$(aws sts get-caller-identity --query Account --output text).dkr.ecr.ap-northeast-2.amazonaws.com

docker pull $ECR_URL/smt-server:$1
docker tag $ECR_URL/smt-server:$1 smt-server:prod

### Deploy
docker-compose -p smt-server up -d

docker image prune -a -f
