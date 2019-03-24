package main


import (
	"demo/service/schema"
	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-mesh/openlogging"
)

func main() {
	chassis.RegisterSchema("rest", &schema.Service{}, server.WithSchemaID("server-demo"))
	if err := chassis.Init(); err != nil {
		openlogging.GetLogger().Error("Init failed." + err.Error())
		return
	}

	chassis.Run()
}
