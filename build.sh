#!/bin/bash -e

TAG=fffbot/fffbot:1.0

echo Building app
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fffbot .

echo Building docker image
docker build -t $TAG -f Dockerfile.scratch .

echo Saving docker image
docker save -o ./fffbot-1.0.tar.gz $TAG
