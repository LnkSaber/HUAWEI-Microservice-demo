package main

import (
	"delay/provider_v2/schema"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-mesh/openlogging"
)

func main() {
	chassis.RegisterSchema("rest", &schema.Provider{}, server.WithSchemaID("Hello-v2"))
	if err := chassis.Init(); err != nil {
		openlogging.GetLogger().Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
