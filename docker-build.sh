#!/bin/sh
# docker-build.sh — build the Docker image and run the container.
#
# Usage:
#   ./docker-build.sh            # build image + start container
#   ./docker-build.sh stop       # stop and remove the container
#   ./docker-build.sh clean      # stop container + remove image

IMAGE_NAME="ascii-art-web-docker"
CONTAINER_NAME="dockerize"
PORT="8080"

stop() {
    echo "Stopping container '$CONTAINER_NAME'..."
    docker container stop "$CONTAINER_NAME" 2>/dev/null || true
    docker container rm   "$CONTAINER_NAME" 2>/dev/null || true
    echo "Done."
}

clean() {
    stop
    echo "Removing image '$IMAGE_NAME'..."
    docker image rm "$IMAGE_NAME" 2>/dev/null || true
    echo "Done."
}

build() {
    echo "Building image '$IMAGE_NAME'..."
    docker image build -f Dockerfile -t "$IMAGE_NAME" .
}

run() {
    docker container rm "$CONTAINER_NAME" 2>/dev/null || true
    echo "Starting container '$CONTAINER_NAME'..."
    docker container run \
        --publish "$PORT:$PORT" \
        --detach \
        --name "$CONTAINER_NAME" \
        "$IMAGE_NAME"
    echo "Container running. Visit: http://localhost:$PORT"
}

case "${1:-}" in
    stop)  stop  ;;
    clean) clean ;;
    *)
        build
        run
        ;;
esac
