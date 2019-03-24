package main

import (
	"delay/provider_v1/schema"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-mesh/openlogging"
)

func main() {
	chassis.RegisterSchema("rest", &schema.Provider{}, server.WithSchemaID("Hello"))
	if err := chassis.Init(); err != nil {
		openlogging.GetLogger().Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
