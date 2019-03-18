#!/usr/bin/env bash

pid=`ps -ef |grep local-cse-config-center.jar|grep -v auto|grep -v 'grep'| awk  '{print $2}'`
kill -9 $pid
echo "Shutting down config center(PID: $pid)..."