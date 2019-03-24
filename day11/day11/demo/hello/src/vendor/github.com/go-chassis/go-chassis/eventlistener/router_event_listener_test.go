package eventlistener_test

import (
	"testing"
	"time"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-archaius/core"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/router"
	"github.com/go-chassis/go-chassis/core/router/cse"
	"github.com/go-chassis/go-chassis/eventlistener"

	"github.com/stretchr/testify/assert"
)

const (
	svcDarkLaunch       = "svcDarkLaunch"
	svcDarkLaunchConfig = `{"policyType":"RATE","ruleItems":[{"groupName":"s0"},{"groupName":"s1"}]}`
)

func TestDarkLaunchEventListenerEvent(t *testing.T) {
	lager.Initialize("", "DEBUG", "", "size", true, 1, 10, 7)

	err := router.Init()
	assert.NoError(t, err)

	e := &core.Event{
		EventSource: cse.RouteDarkLaunchGovernSourceName,
		EventType:   core.Create,
		Key:         svcDarkLaunch,
		Value:       svcDarkLaunchConfig,
	}

	t.Log("Before event, there should be no router config")
	assert.Nil(t, router.DefaultRouter.FetchRouteRuleByServiceName(svcDarkLaunch))

	t.Log("After event, there should exists router config")
	archaius.AddKeyValue(eventlistener.DarkLaunchPrefix+svcDarkLaunch, svcDarkLaunchConfig)
	l := &eventlistener.DarkLaunchEventListener{}
	l.Event(e)
	time.Sleep(100 * time.Millisecond)
	r := router.DefaultRouter.FetchRouteRuleByServiceName(svcDarkLaunch)
	assert.NotNil(t, r)
	if r == nil {
		t.FailNow()
	}
	assert.Equal(t, 2, len(r[0].Routes))
}
