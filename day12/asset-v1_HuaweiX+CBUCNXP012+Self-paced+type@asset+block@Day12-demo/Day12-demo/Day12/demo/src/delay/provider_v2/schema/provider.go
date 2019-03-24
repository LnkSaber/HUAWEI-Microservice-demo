package schema

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-mesh/openlogging"
)

type Provider struct {
}

func (s *Provider) Delay(ctx *restful.Context) {

	t := ctx.ReadPathParameter("time")
	if t != "" {
		n, err := strconv.Atoi(t)
		if err != nil {
			openlogging.GetLogger().Error("string to int error")
			ctx.WriteError(http.StatusInternalServerError, err)
			return
		}
		time.Sleep(time.Duration(n) * time.Millisecond)
	}
	ctx.Write([]byte(fmt.Sprintf("%v,version v2:delay method for provider_v1 sleep %s ms", runtime.InstanceID, t)))
}

// URLPatterns
func (*Provider) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/provider/v0/delay/{time}", ResourceFuncName: "Delay"},
	}
}
