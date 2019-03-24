package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core/client"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	"net/http"
)

// Test is the function to call provider health check api and check the response
func Test(ctx context.Context, protocol, endpoint string, expected Reply) (err error) {
	switch protocol {
	case common.ProtocolRest:
		err = restTest(ctx, endpoint, expected)
	case common.ProtocolHighway:
		err = highwayTest(ctx, endpoint, expected)
	default:
		err = fmt.Errorf("Unsupport protocol %s", protocol)
	}
	return
}

func restTest(ctx context.Context, endpoint string, expected Reply) (err error) {
	c, err := client.GetClient(common.ProtocolRest, expected.ServiceName, "")
	if err != nil {
		return
	}

	arg, _ := rest.NewRequest(http.MethodGet, "http://"+expected.ServiceName+"/healthz", nil)
	req := &invocation.Invocation{Args: arg}
	rsp := rest.NewResponse()
	err = c.Call(ctx, endpoint, req, rsp)
	if rsp.Body != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return
	}
	if rsp.StatusCode != http.StatusOK {
		return nil
	}
	var actual Reply
	err = json.Unmarshal(httputil.ReadBody(rsp), &actual)
	if err != nil {
		return
	}
	if actual != expected {
		return fmt.Errorf("endpoint is belong to %s:%s:%s",
			actual.ServiceName, actual.Version, actual.AppId)
	}
	return
}

func highwayTest(ctx context.Context, endpoint string, expected Reply) (err error) {
	c, err := client.GetClient(common.ProtocolHighway, expected.ServiceName, "")
	if err != nil {
		return
	}
	req := &invocation.Invocation{
		MicroServiceName: expected.ServiceName,
		SchemaID:         "_chassis_highway_healthz",
		OperationID:      "HighwayCheck",
		Args:             &Request{},
	}
	var actual Reply
	err = c.Call(ctx, endpoint, req, &actual)
	if err != nil {
		return
	}
	if actual != expected {
		return fmt.Errorf("Endpoint is belong to %s:%s:%s",
			actual.ServiceName, actual.Version, actual.AppId)
	}
	return
}
