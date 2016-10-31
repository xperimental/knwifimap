#!/bin/bash -xe

TAG=$(git rev-parse --short HEAD); readonly TAG

go build --ldflags '-w -extldflags "-static"' -o knwifimap .
docker build -t xperimental/knwifimap:$TAG .
docker tag xperimental/knwifimap:$TAG xperimental/knwifimap:latest
