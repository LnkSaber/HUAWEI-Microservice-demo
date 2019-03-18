#!/usr/bin/env bash

cd ./local-cse-config-center
chmod +x *.sh
sh start.sh
cd ../local-service-center
chmod +x *.sh
bash start.sh
cd ../console
chmod +x *.sh
sh start.sh