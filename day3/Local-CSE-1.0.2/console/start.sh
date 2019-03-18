#!/usr/bin/env bash

export APP_ROOT=`pwd`
java -jar console-website.jar > ./console-website.log 2>&1 &
echo "console started, listening on 0.0.0.0:30106..."