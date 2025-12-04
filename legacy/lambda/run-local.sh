#!/bin/sh

### Run DynamoDB container
docker compose -p dynamodb-local up -d

### Enable golang build automation
air &

### Run serverless
cd local
serverless offline start &

### Wait until the processes are finished
trap "kill $!" INT
wait $!

### Clean up
docker compose -p dynamodb-local down
docker ps -a --filter "status=exited" --filter "ancestor=lambci/lambda:go1.x" -q | xargs -r docker rm
