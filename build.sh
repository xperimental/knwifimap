#!/bin/bash -xe

GITHASH=$(git rev-parse --short HEAD); readonly GITHASH
DBHASH=$(shasum <database.sqlite | awk '{print $1}'); readonly DBHASH
TAG="${GITHASH}-${DBHASH}"

go build --ldflags '-w -extldflags "-static"' -o knwifimap .
docker build -t xperimental/knwifimap:$TAG .
docker tag xperimental/knwifimap:$TAG xperimental/knwifimap:latest
