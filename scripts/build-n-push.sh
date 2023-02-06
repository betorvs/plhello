#!/usr/bin/env bash

REPOSITORY="${1:-betorvs/plhello}"
TAG="${2:-latest}"
VERSION=${3:-develop}

docker build -t ${REPOSITORY}:${TAG} \
    -f Dockerfile \
    --build-arg BUILD_REF=${VERSION} \
    .
docker push ${REPOSITORY}:${TAG}
