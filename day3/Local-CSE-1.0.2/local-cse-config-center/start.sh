#!/usr/bin/env bash

export STORAGE_MODE=file

java -jar ./local-cse-config-center.jar > ./config-center.log 2>&1 &

echo "config center started, listening on 0.0.0.0:30113..."