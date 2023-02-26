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

#!/bin/sh

if [ $# -ne 1 ]; then
	echo 'No argument'
	exit -1
fi

### Load Next Build From S3
aws s3 cp s3://seolmyeongtang-cicd/client/$1.tar.gz $1.tar.gz

rm -rf build/*

### Unzip Next Build
tar -zxf $1.tar.gz -C build --strip 1

### Cleanup
rm $1.tar.gz

