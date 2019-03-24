### Config Center client for Go-Chassis
[![Build Status](https://travis-ci.org/go-chassis/go-cc-client.svg?branch=master)](https://travis-ci.org/go-chassis/go-cc-client)  

This is cc client which interacts with config server to get the configurations for 
particular micro-services. It can create a web socket connection with the config server
and receive all the change events in any of the configuration for the particular micro-service.
It can also use rest http connections to pull the data from config server at regular intervals.
