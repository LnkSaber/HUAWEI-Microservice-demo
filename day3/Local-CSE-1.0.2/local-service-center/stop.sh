#!/bin/sh

CURRENT_DIR=$(cd $(dirname $0); pwd)

ps -ef | grep ${CURRENT_DIR}/bin/frontend | grep -v grep | awk -F ' ' '{print $2}'|while read line
do
  echo "Shutting down frontend(PID: $line)..."
  eval "kill -9 $line"
done

ps -ef | grep ${CURRENT_DIR}/bin/servicecenter | grep -v grep | awk -F ' ' '{print $2}'|while read line
do
  echo "Shutting down service center(PID: $line)..."
  eval "kill -9 $line"
done

echo "Done!"
