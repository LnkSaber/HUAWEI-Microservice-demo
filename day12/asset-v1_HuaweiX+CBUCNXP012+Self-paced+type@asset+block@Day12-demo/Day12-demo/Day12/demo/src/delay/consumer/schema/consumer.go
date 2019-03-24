package consumer

import (
	"context"
	"net/http"

	"fmt"
	"strconv"

	"time"

	"sync"

	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	"github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-mesh/openlogging"
)

// default call provider_v1 times
const DefaultTimes = time.Second * 10
const DefaultCount int = 10
const DefaultDelayTime string = "30"

type Consumer struct {
}
type result struct {
	Reply string `json:"reply"`
	Error string `json:"error"`
	Num   int    `json:"num"`
}
type replys struct {
	TestType string    `json:"test_type"`
	Replys   *[]result `json:"replys"`
}

var lock sync.Mutex
var restInvoker = core.NewRestInvoker()

func (*Consumer) Delay(ctx *restful.Context) {
	// get delay time
	delay := ctx.ReadQueryParameter("delay")
	if delay == "" {
		delay = DefaultDelayTime
	}

	times := ctx.ReadQueryParameter("t")
	duration, err := time.ParseDuration(times)
	if err != nil {
		duration = DefaultTimes
	}
	c := ctx.ReadQueryParameter("c")
	cs, err := strconv.Atoi(c)
	if err != nil {
		cs = DefaultCount
	}

	results := &[]result{}

	fault_limit_circuit_router(cs, delay, duration, results)

	ctx.WriteJSON(replys{TestType: "router", Replys: results}, "application/json")
}
func fault_limit_circuit_router(c int, delay string, duration time.Duration, result *[]result) {

	cancels := make([]context.CancelFunc, 0)
	for i := 0; i < c; i++ {
		openlogging.GetLogger().Info("launched one http benchmark thread")
		ctx, cancel := context.WithCancel(context.Background())
		cancels = append(cancels, cancel)
		go callHTTP(ctx, restInvoker, http.MethodGet, fmt.Sprintf("http://provider/provider/v0/delay/%s", delay), result)
	}
	time.Sleep(duration)
	for _, cancel := range cancels {
		cancel()
	}
}

func callHTTP(ctx context.Context, restInvoker *core.RestInvoker, method string, url string, r *[]result) {
	for {

		req, err := rest.NewRequest(method, url, nil)
		if err != nil {
			panic(err)
		}
		//req.Header.Set("Foo", "bar")

		resp, err := restInvoker.ContextDo(ctx, req)
		if err != nil {
			lock.Lock()
			*r = append(*r, result{Num: len(*r) + 1, Error: err.Error()})
			lock.Unlock()
			if !endOfCtx(ctx) {
				break
			}
			continue
		}

		lock.Lock()
		*r = append(*r, result{Num: len(*r) + 1, Reply: fmt.Sprintf(string(httputil.ReadBody(resp))+"time :%d", time.Now().UnixNano()/1e6)})
		lock.Unlock()
		resp.Body.Close()
		if !endOfCtx(ctx) {
			break
		}
	}
}
func endOfCtx(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		return true
	}
}

// URLPatterns
func (*Consumer) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/consumer/v0/delay", ResourceFuncName: "Delay"},
	}
}
