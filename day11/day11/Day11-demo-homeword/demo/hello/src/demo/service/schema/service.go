package schema

import (
	"fmt"
	"net/http"

	"time"

	"github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-mesh/openlogging"
)

const Male string = "MALE"
const Female string = "FEMALE"

type GreetingResponse struct {
	Msg       string        `json:"msg"`
	Timestamp time.Duration `json:"timestamp"`
}
type Person struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

type Service struct{}

func (*Service) Hello(ctx *restful.Context) {
	openlogging.GetLogger().Info("access provider hello")
	name := ctx.ReadPathParameter("name")
	ctx.WriteJSON(fmt.Sprintf("hello , %s", name), "application/json")
}
func (*Service) Greeting(ctx *restful.Context) {
	openlogging.GetLogger().Info("access provider greeting")
	param := Person{}
	ctx.ReadEntity(&param)
	r := GreetingResponse{}
	if param.Gender == Male {
		r.Msg = fmt.Sprintf("Hello , Mr.%s", param.Name)
	} else if param.Gender == Female {
		r.Msg = fmt.Sprintf("Hello , Ms.%s", param.Name)
	}
	r.Timestamp = time.Duration(time.Now().UnixNano() / 1e6)
	ctx.WriteJSON(r, "application/json")
}

// URLPatterns
func (*Service) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/provider/v0/hello/{name}", ResourceFuncName: "Hello"},
		{Method: http.MethodPost, Path: "/provider/v0/greeting", ResourceFuncName: "Greeting"},
	}
}
