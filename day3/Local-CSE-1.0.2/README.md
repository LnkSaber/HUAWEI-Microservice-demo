# Local Cloud-Service-Engine 

Local CSE is a GUI-based tool, used for local deploymentof lightweight service center and configuration center.

## Features
 - Integrates Local Service Center.
 - Integrates Local Config Center.
 - Provides a Lightweight Console.

## Quick Start

### Getting Local Cloud-Service-Engine

Download the latest releases of [Local Cloud-Service-Engine]

### Running Local Cloud-Service-Engine

When you get these package, you can execute the start script to run Local Cloud-Service-Engine.

Windows:
```
start.bat
```

Linux:
```sh
./start.sh
```

Ensure that the following ports are not occupied before starting：
```
service-center: 30100 
internal-usage: 30103
console-website: 30106
config-center: 30113 & 30114
```
Services EntryPoints:
- config-center:
```
http://${HOST_IP}:30113
```
- service-center:
```
http://${HOST_IP}:30100
```
- console website:
```
http://${HOST_IP}:30106
```
Tip: HOST_IP is the ip of your machine, you can also use the [localhost].

[localhost]: 127.0.0.1

Note:So far, it's not allowed to change these listened ports.

### Access the Local Service Center
- Edit the **microservices.yaml** in the microservice
```
cse:
  service:
    registry:
      address: http://${HOST_IP}:30100
      instance:
        watch: false
  config:
    client:
      serverUri: http://${HOST_IP}:30113
      refreshMode: 1                     # 0:config-center push configs to service 1:service pull configs from config-center, 0 by default 
      refresh_interval: 5000                   
      #refreshPort: 30114                #optional, the websocket port when refresh mode is 0
```
- Open console website in browser: **http://${HOST_IP}:30106**


### Configure/Update Local Service Center

#### The Releases of Local Cloud-Service-Engine integrate [Local Service Center], packaged at directory [local-service-center]
- You can configure the integrated Service Center via [README.md] in directory [local-service-center]
- Update the integrated Service Center:
    - unzip the releases of [Local Service Center] and find ./bin (directory)
    - Windows: replace [local-service-center/bin]/servicecenter.exe with ./bin/servicecenter.exe
    - Linux: replace [local-service-center/bin]/servicecenter with ./bin/servicecenter


[local-service-center]: ./local-service-center
[local-service-center/bin]: ./local-service-center/bin
[README.md]: ./local-service-center/README.md

## Get The Latest Release

[Local Service Center]

[Local Service Center]: https://console.huaweicloud.com/servicestage/#/cse/tools
[Local Cloud-Service-Engine]: https://console.huaweicloud.com/servicestage/#/cse/tools