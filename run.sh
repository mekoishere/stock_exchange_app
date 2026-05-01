#!/bin/bash
if [ -z "$1" ]; then
    echo "Usage: ./run.sh <port>"
    exit 1
fi

export APP_EXTERNAL_PORT=$1
docker-compose build --no-cache
docker-compose up
