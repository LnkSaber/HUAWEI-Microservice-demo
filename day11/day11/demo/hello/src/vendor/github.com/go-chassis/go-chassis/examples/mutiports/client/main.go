package main

import (
	"context"

	"github.com/go-chassis/go-chassis"
	_ "github.com/go-chassis/go-chassis/bootstrap"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/client/
func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}

	req, err := rest.NewRequest("GET", "http://RESTServer/hello", nil)
	if err != nil {
		lager.Logger.Error("new request failed." + err.Error())
		return
	}
	defer req.Body.Close()

	resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req)
	if err != nil {
		lager.Logger.Error("do request failed." + err.Error())
		return
	}
	defer resp.Body.Close()
	lager.Logger.Info("REST Server sayhello[GET]: " + string(httputil.ReadBody(resp)))

	req, err = rest.NewRequest("GET", "http://RESTServer:legacy/legacy", nil)
	if err != nil {
		lager.Logger.Error("new request failed." + err.Error())
		return
	}
	defer req.Body.Close()

	resp, err = core.NewRestInvoker().ContextDo(context.TODO(), req)
	if err != nil {
		lager.Logger.Error("do request failed." + err.Error())
		return
	}
	defer resp.Body.Close()
	lager.Logger.Info("REST Server sayhello[GET]: " + string(httputil.ReadBody(resp)))
}
