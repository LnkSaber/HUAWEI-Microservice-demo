package main

import (
	"context"

	"github.com/go-chassis/go-chassis"
	_ "github.com/go-chassis/go-chassis/bootstrap"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	"time"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/client/
func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	// run is only to enable metric exporter
	go chassis.Run()
	ctx := context.WithValue(context.TODO(), common.ContextHeaderKey{}, map[string]string{
		"user": "peter",
	})
	invoker := core.NewRestInvoker()

	for i := 0; i < 500; i++ {
		req, err := rest.NewRequest("GET", "http://ErrServer/lock", nil)
		if err != nil {
			lager.Logger.Error("new request failed.")
			return
		}
		resp, err := invoker.ContextDo(ctx, req)
		if err != nil {
			lager.Logger.Errorf("deadlock request failed. %s", err.Error())
		} else {
			lager.Logger.Info("REST Server [GET]: " + string(httputil.ReadBody(resp)))
		}

	}

	for {
		req, err := rest.NewRequest("GET", "http://ErrServer/sayhimessage", nil)
		if err != nil {
			lager.Logger.Error("new request failed.")
			return
		}
		resp, err := invoker.ContextDo(ctx, req)
		if err != nil {
			lager.Logger.Errorf("normal request failed. %s", err.Error())
		} else {
			lager.Logger.Info("REST Server [GET]: " + string(httputil.ReadBody(resp)))
		}
		time.Sleep(30 * time.Second)
	}
}
