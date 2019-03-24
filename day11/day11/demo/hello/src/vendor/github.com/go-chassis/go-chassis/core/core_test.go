package core_test

import (
	_ "github.com/go-chassis/go-chassis/initiator"

	"context"
	"testing"

	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	_ "github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/examples/schemas/helloworld"
	"github.com/go-chassis/go-chassis/pkg/util/tags"

	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	"github.com/stretchr/testify/assert"
)

func initenv() {
	config.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	//handler.Init()
}

func TestRPCInvoker_InvokeFailinChainInit(t *testing.T) {
	initenv()
	config.GlobalDefinition = &model.GlobalCfg{}
	invoker := core.NewRPCInvoker(core.ChainName(""))
	replyOne := &helloworld.HelloReply{}
	ctx := context.WithValue(context.Background(), common.ContextHeaderKey{}, map[string]string{
		"X-User": "tianxiaoliang",
	})

	config.GlobalDefinition.Cse.References = make(map[string]model.ReferencesStruct)
	version := model.ReferencesStruct{Version: ""}
	config.GlobalDefinition.Cse.References["Server"] = version
	err := invoker.Invoke(ctx, "Server", "HelloServer", "SayHello", &helloworld.HelloRequest{Name: "Peter"}, replyOne,
		core.WithMetadata(nil), core.WithStrategy(""), core.StreamingRequest())
	assert.Error(t, err)
}
func TestRestInvoker_ContextDo(t *testing.T) {
	initenv()
	restinvoker := core.NewRestInvoker()
	req, _ := rest.NewRequest("GET", "http://Server/sayhello/myidtest", nil)
	httputil.SetContentType(req, "application/json")
	//use the invoker like http client.
	_, err := restinvoker.ContextDo(context.TODO(), req, core.WithEndpoint("0.0.0.0"), core.WithProtocol("rest"), core.WithFilters(""))
	assert.Error(t, err)
}

func TestOptions(t *testing.T) {
	opt := core.InvokeOptions{}
	option := core.DefaultCallOptions(opt)
	assert.NotEmpty(t, option)

	inv := core.StreamingRequest()
	assert.NotEmpty(t, inv)

	inv = core.WithEndpoint("0.0.0.0")
	assert.NotEmpty(t, inv)

	inv = core.WithProtocol("0.0")
	assert.NotEmpty(t, inv)

	inv = core.WithFilters("")
	assert.NotEmpty(t, inv)

	inv = core.WithStrategy("")
	assert.NotEmpty(t, inv)

	inv = core.WithMetadata(nil)
	assert.NotEmpty(t, inv)

	t.Log("with router tag")
	testKey := "testKey"
	testValue := "testValue"
	m := map[string]string{
		testKey: testValue,
	}
	op := core.WithRouteTags(m)
	assert.NotNil(t, op)
	op(&opt)
	assert.Equal(t, testValue, opt.RouteTags.KV[testKey])
	assert.Empty(t, opt.RouteTags.KV[common.BuildinTagApp])
	assert.Equal(t, utiltags.LabelOfTags(opt.RouteTags.KV), opt.RouteTags.Label)
}
