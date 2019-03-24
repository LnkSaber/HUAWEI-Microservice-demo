package schema

import (
	"fmt"
	"net/http"

	"github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-mesh/openlogging"
)


type Service struct{}

func (*Service) Hello(ctx *restful.Context) {
	openlogging.GetLogger().Info("access provider hello")
	name := ctx.ReadPathParameter("name")
	ctx.WriteJSON(fmt.Sprintf("hello , %s", name), "application/json")
}

// URLPatterns
func (*Service) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/provider/v0/hello/{name}", ResourceFuncName: "Hello"},
	}
}
