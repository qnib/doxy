#!/bin/bash

COMMIT_TAG=$(git describe --tags)
TEMP_TAG=$(base64 /dev/urandom | tr -d '/+' | dd bs=32 count=1 2>/dev/null)
mkdir -p bin/

echo ">> Build temp-image 'qnib/$(basename $(pwd)):${TEMP_TAG}'"
docker build -t qnib/$(basename $(pwd)):${TEMP_TAG} -f Dockerfile.ubuntu .
echo ">> Start image as ${TEMP_TAG} to copy binary"
ID=$(docker run -d --name ${TEMP_TAG} qnib/$(basename $(pwd)):${TEMP_TAG} tail -f /dev/null)
echo ">> CONTAINER_ID=${ID}"
docker cp ${ID}:/usr/local/bin/doxy bin/doxy_x86_${COMMIT_TAG}
echo ">> Remove container and image"
docker rm -f ${ID}
docker rmi qnib/$(basename $(pwd)):${TEMP_TAG}
