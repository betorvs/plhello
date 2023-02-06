#!/usr/bin/env bash

REPOSITORY="${1:-betorvs/plhello}"
TAG="${2:-latest}"
VERSION=${3:-develop}
TARGET=${4:-local}

docker build -t ${REPOSITORY}:${TAG} \
    -f Dockerfile \
    --build-arg BUILD_REF=${VERSION} .

if [ "${TARGET}" = "local" ]; then 
    # to avoid error when pushing images to local registry in gh workflow
    echo "pushing using localhost"
    docker tag ${REPOSITORY}:${TAG} localhost:5050/plhello:${TAG}
    docker push localhost:5050/plhello:${TAG}
else 
    echo "pushing using to ${REPOSITORY}"
   docker push ${REPOSITORY}:${TAG}
fi

