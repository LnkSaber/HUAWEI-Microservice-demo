package main

import (
	"demo/client/schema"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-mesh/openlogging"
)

func main() {
	chassis.RegisterSchema("rest", &schema.Client{}, server.WithSchemaID("client-demo"))
	if err := chassis.Init(); err != nil {
		openlogging.GetLogger().Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
