#!/bin/bash

DOCKER_HUB_USERNAME="ery4z"
IMAGE_NAME="link-interceptor"
IMAGE_VERSION="latest" # Use 'latest' if version is not provided

# Building the Docker image
docker build -t ${DOCKER_HUB_USERNAME}/${IMAGE_NAME}:${IMAGE_VERSION} .

# Check if the build was successful
if [ $? -ne 0 ]; then
    echo "Docker build failed. Exiting."
    exit 1
fi

echo "Successfully built ${DOCKER_HUB_USERNAME}/${IMAGE_NAME}:${IMAGE_VERSION}"

# Pushing the Docker image to Docker Hub
docker push ${DOCKER_HUB_USERNAME}/${IMAGE_NAME}:${IMAGE_VERSION}

# Check if the push was successful
if [ $? -ne 0 ]; then
    echo "Failed to push the image to Docker Hub. Exiting."
    exit 1
fi

echo "Successfully pushed ${DOCKER_HUB_USERNAME}/${IMAGE_NAME}:${IMAGE_VERSION} to Docker Hub"
