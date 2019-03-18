#!/usr/bin/env bash

pid=`ps -ef |grep console-website.jar|grep -v auto|grep -v 'grep'| awk  '{print $2}'`
kill -9 $pid
echo "Shutting down console(PID: $pid)..."