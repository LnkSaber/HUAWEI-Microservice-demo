package schema

import (
	"context"
	"net/http"

	"fmt"

	"encoding/json"

	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	"github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-mesh/openlogging"
)

type Client struct {
}

func (*Client) Hello(ctx *restful.Context) {
	openlogging.GetLogger().Info("access consumer Hello")
	//可以使用 cse:// 和 http://作为前缀
	//req, err := rest.NewRequest(http.MethodGet, fmt.Sprintf("http://provider/provider/v0/hello/%s", ctx.ReadQueryParameter("name")), nil)
	req, err := rest.NewRequest(http.MethodGet, fmt.Sprintf("cse://provider/provider/v0/hello/%s", ctx.ReadQueryParameter("name")), nil)
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req)
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err)
		return
	}
	var reply string
	b := httputil.ReadBody(resp)
	json.Unmarshal(b, &reply)
	ctx.WriteJSON(reply, "application/json")
}

// URLPatterns
func (*Client) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/consumer/v0/hello", ResourceFuncName: "Hello"},
	}
}
