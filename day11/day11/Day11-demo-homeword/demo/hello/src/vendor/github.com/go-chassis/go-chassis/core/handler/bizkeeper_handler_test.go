package handler_test

import (
	"log"
	"os"
	"testing"

	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/examples/schemas/helloworld"
	"github.com/stretchr/testify/assert"
)

func TestCBInit(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/server/")

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	err := control.Init()
	assert.NoError(t, err)
}

func TestBizKeeperConsumerHandler_Handle(t *testing.T) {
	t.Log("testing bizkeeper consumer handler")

	c := handler.Chain{}
	c.AddHandler(&handler.BizKeeperConsumerHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = make(map[string]string)
	config.GlobalDefinition.Cse.Handler.Chain.Consumer["bizkeeperconsumerdefault"] = "bizkeeper-consumer"
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
	}

	c.Next(i, func(r *invocation.Response) error {
		assert.NoError(t, r.Err)
		log.Println(r.Result)
		return r.Err
	})
}
func TestBizKeeperProviderHandler_Handle(t *testing.T) {
	t.Log("testing bizkeeper provider handler")

	c := handler.Chain{}
	c.AddHandler(&handler.BizKeeperProviderHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Handler.Chain.Provider = make(map[string]string)
	config.GlobalDefinition.Cse.Handler.Chain.Provider["bizkeeperproviderdefault"] = "bizkeeper-provider"
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
	}

	c.Next(i, func(r *invocation.Response) error {
		assert.NoError(t, r.Err)
		log.Println(r.Result)
		return r.Err
	})
}

func TestBizKeeperHandler_Names(t *testing.T) {
	bizPro := &handler.BizKeeperProviderHandler{}
	proName := bizPro.Name()
	assert.Equal(t, "bizkeeper-provider", proName)

	bizCon := &handler.BizKeeperConsumerHandler{}
	conName := bizCon.Name()
	assert.Equal(t, "bizkeeper-consumer", conName)

}

func BenchmarkBizKeepConsumerHandler_Handler(b *testing.B) {
	b.Log("benchmark for bizkeeper consumer handler")
	c := handler.Chain{}
	c.AddHandler(&handler.BizKeeperConsumerHandler{})
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/client/")

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	control.Init()
	inv := &invocation.Invocation{
		MicroServiceName: "fakeService",
		SchemaID:         "schema",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Next(inv, func(r *invocation.Response) error {
			return r.Err
		})
		c.Reset()
	}
}
