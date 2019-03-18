#!/bin/sh

set -e

CURRENT_DIR=$(cd $(dirname $0); pwd)
cd ${CURRENT_DIR}

if [ "$1"x == "-image"x ]; then
    ${CURRENT_DIR}/bin/frontend &
    ${CURRENT_DIR}/bin/servicecenter
    exit 0
else
    sed -i s@"^logfile.*=.*$"@"logfile = ./service-center.log"@g ./conf/app.conf
    ${CURRENT_DIR}/bin/frontend > ./frontend.log 2>&1 &
    ${CURRENT_DIR}/bin/servicecenter > /dev/null 2>&1 &
fi

echo "Service center has been launched."
echo "you can check ${CURRENT_DIR}/service-center.log for more details."
