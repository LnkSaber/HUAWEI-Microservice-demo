package handler_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-chassis/go-archaius/core"
	"github.com/go-chassis/go-archaius/core/cast"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/config"
	chassisModel "github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/core/registry"
	mk "github.com/go-chassis/go-chassis/core/registry/mock"
	_ "github.com/go-chassis/go-chassis/core/registry/servicecenter"
	"github.com/go-chassis/go-chassis/examples/schemas/helloworld"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
)

const CallTimes = 15

var callTimes = 0

type handler1 struct {
}

func (th *handler1) Name() string {
	return "loadbalancer"
}

func (th *handler1) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	callTimes++
	cb(&invocation.Response{})
}

type handler2 struct {
}

func (h *handler2) Name() string {
	return "handler2"
}

func (h *handler2) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	callTimes++
	r := &invocation.Response{
		Err: fmt.Errorf("A fake error from handler2"),
	}
	if callTimes < CallTimes {
		cb(r)
		return
	}
	cb(&invocation.Response{})
}

/*======================================================================================================================
       Mocking ConfigurationFactory interface with a dummy struct
=======================================================================================================================*/
type MockConfigurationFactory struct {
	mock.Mock
}

func (m *MockConfigurationFactory) Init() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConfigurationFactory) GetConfigurations() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}

func (m *MockConfigurationFactory) GetConfigurationsByDimensionInfo(dimensionInfo string) map[string]interface{} {
	args := m.Called(dimensionInfo)
	return args.Get(0).(map[string]interface{})
}

func (m *MockConfigurationFactory) GetConfigurationByKey(key string) interface{} {
	args := m.Called(key)
	return args.Get(0)
}
func (m *MockConfigurationFactory) IsKeyExist(key string) bool {
	args := m.Called(key)
	return args.Bool(0)
}
func (m *MockConfigurationFactory) Unmarshal(structure interface{}) error {
	args := m.Called(structure)
	return args.Error(0)
}
func (m *MockConfigurationFactory) AddSource(cfg core.ConfigSource) error {
	args := m.Called(cfg)
	return args.Error(0)
}
func (m *MockConfigurationFactory) RegisterListener(listenerObj core.EventListener, key ...string) error {
	args := m.Called(listenerObj, key)
	return args.Error(0)
}
func (m *MockConfigurationFactory) UnRegisterListener(listenerObj core.EventListener, key ...string) error {
	args := m.Called(listenerObj, key)
	return args.Error(0)
}
func (m *MockConfigurationFactory) DeInit() error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockConfigurationFactory) GetValue(key string) cast.Value {
	args := m.Called(key)
	return args.Get(0).(cast.Value)
}

func (m *MockConfigurationFactory) GetConfigurationByKeyAndDimensionInfo(dimensionInfo, key string) interface{} {
	args := m.Called(key)
	return args.Get(0)
}

func (m *MockConfigurationFactory) GetValueByDI(dimensionInfo, key string) cast.Value {
	args := m.Called(key)
	return args.Get(0).(cast.Value)
}

func (m *MockConfigurationFactory) AddByDimensionInfo(dimensionInfo string) (map[string]string, error) {
	args := m.Called(dimensionInfo)
	return args.Get(0).(map[string]string), nil
}

/*======================================================================================================================
			                              END
 =======================================================================================================================*/

func TestLBHandlerWithRetry(t *testing.T) {
	microContent := `---
#微服务的私有属性
service_description:
  name: Client
  level: FRONT
  version: 0.1`

	os.Setenv("CHASSIS_HOME", "/tmp")
	defer os.Unsetenv("CHASSIS_HOME")
	chassisConf := filepath.Join("/tmp/", "conf")
	logConf := filepath.Join("/tmp/", "log")
	err := os.MkdirAll(chassisConf, 0700)
	assert.NoError(t, err)
	err = os.MkdirAll(logConf, 0700)
	assert.NoError(t, err)
	chassisyaml := filepath.Join(chassisConf, "chassis.yaml")
	microserviceyaml := filepath.Join(chassisConf, "microservice.yaml")
	f1, err := os.Create(chassisyaml)
	assert.NoError(t, err)
	f2, err := os.Create(microserviceyaml)
	assert.NoError(t, err)
	_, err = io.WriteString(f1, yamlContent)
	assert.NoError(t, err)
	_, err = io.WriteString(f2, microContent)

	t.Log("testing load balance handler with retry")
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	runtime.ServiceID = "selfServiceID"
	err = config.Init()
	assert.NoError(t, err)
	err = control.Init()
	assert.NoError(t, err)
	testConfigFactoryObj := new(MockConfigurationFactory)
	key1 := fmt.Sprint("cse.loadbalance.retryEnabled")
	key2 := fmt.Sprint("cse.loadbalance.service1.retryEnabled")
	key3 := fmt.Sprint("cse.loadbalance.retryOnSame")
	key4 := fmt.Sprint("cse.loadbalance.service1.retryOnSame")
	key5 := fmt.Sprint("cse.loadbalance.retryOnNext")
	key6 := fmt.Sprint("cse.loadbalance.service1.retryOnNext")
	key8 := fmt.Sprint("cse.loadbalance.backoff.kind")
	key9 := fmt.Sprint("cse.loadbalance.service1.backoff.kind")
	key10 := fmt.Sprint("cse.loadbalance.backoff.minMs")
	key11 := fmt.Sprint("cse.loadbalance.service1.backoff.minMs")
	key12 := fmt.Sprint("cse.loadbalance.backoff.maxMs")
	key13 := fmt.Sprint("cse.loadbalance.service1.backoff.maxMs")
	key14 := fmt.Sprint("cse.loadbalance.strategy.name")
	key15 := fmt.Sprint("cse.loadbalance.SessionStickinessRule.sessionTimeoutInSeconds")
	key16 := fmt.Sprint("cse.loadbalance.SessionStickinessRule.successiveFailedTimes")
	key17 := fmt.Sprint("cse.loadbalance.serverListFilters")
	val1 := cast.NewValue(true, nil)
	val2 := cast.NewValue(1, nil)
	val4 := cast.NewValue("jittered", nil)
	val5 := cast.NewValue("Random", nil)
	val6 := cast.NewValue(10, nil)
	val7 := cast.NewValue(2, nil)
	val8 := cast.NewValue(loadbalancer.ZoneAware, nil)
	testConfigFactoryObj.On("GetValue", key1).Return(val1)
	testConfigFactoryObj.On("GetValue", key2).Return(val1)
	testConfigFactoryObj.On("GetValue", key3).Return(val2)
	testConfigFactoryObj.On("GetValue", key4).Return(val2)
	testConfigFactoryObj.On("GetValue", key5).Return(val2)
	testConfigFactoryObj.On("GetValue", key6).Return(val2)
	testConfigFactoryObj.On("GetValue", key8).Return(val4)
	testConfigFactoryObj.On("GetValue", key9).Return(val4)
	testConfigFactoryObj.On("GetValue", key10).Return(val2)
	testConfigFactoryObj.On("GetValue", key11).Return(val2)
	testConfigFactoryObj.On("GetValue", key12).Return(val2)
	testConfigFactoryObj.On("GetValue", key13).Return(val2)
	testConfigFactoryObj.On("GetConfigurationByKey", key14).Return(val5)
	testConfigFactoryObj.On("GetValue", key14).Return(val5)
	testConfigFactoryObj.On("GetValue", key15).Return(val6)
	testConfigFactoryObj.On("GetValue", key16).Return(val7)
	testConfigFactoryObj.On("GetValue", key17).Return(val8)

	c := handler.Chain{}
	c.AddHandler(&handler.LBHandler{})

	var mss []*registry.MicroServiceInstance
	var ms1 = &registry.MicroServiceInstance{
		InstanceID:   "instanceID",
		EndpointsMap: map[string]string{"rest": "127.0.0.1"},
	}
	var ms2 = new(registry.MicroServiceInstance)
	ms2.EndpointsMap = map[string]string{"rest": "127.0.0.1"}
	ms2.InstanceID = "ins2"
	mss = append(mss, ms1)
	mss = append(mss, ms2)

	testRegistryObj := new(mk.DiscoveryMock)
	registry.DefaultServiceDiscoveryService = testRegistryObj
	testRegistryObj.On("FindMicroServiceInstances", "selfServiceID", "appID", "service1", "1.0", "").Return(mss, nil)

	config.GlobalDefinition = &chassisModel.GlobalCfg{}
	config.GetLoadBalancing().Strategy = make(map[string]string)
	loadbalancer.Enable()
	req, _ := rest.NewRequest("GET", "127.0.0.1", nil)
	req.Header.Set("Set-Cookie", "sessionid=100")
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             req,
		Strategy:         loadbalancer.StrategyRoundRobin,
		RouteTags:        utiltags.NewDefaultTag("1.0", "appID"),
		SourceServiceID:  runtime.ServiceID,
	}
	t.Log(i.SourceServiceID)
	c.Next(i, func(r *invocation.Response) error {
		assert.NoError(t, r.Err)
		return r.Err
	})

	var lbh = new(handler.LBHandler)
	str := lbh.Name()
	assert.Equal(t, "loadbalancer", str)
	t.Log(i.Protocol)
	t.Log(i.Endpoint)
}
func TestLBHandlerWithNoRetry(t *testing.T) {
	t.Log("testing load balance handler with No retry")
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "client"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	err := config.Init()
	assert.NoError(t, err)
	//config.GlobalDefinition = &chassisModel.GlobalCfg{}

	err = control.Init()
	assert.NoError(t, err)
	c := handler.Chain{}
	c.AddHandler(&handler.LBHandler{})

	var mss []*registry.MicroServiceInstance
	var ms1 = new(registry.MicroServiceInstance)
	var ms2 = new(registry.MicroServiceInstance)
	var mp = make(map[string]string)
	mp["any"] = "127.0.0.1"
	ms1.EndpointsMap = mp
	ms1.InstanceID = "ins1"
	ms2.EndpointsMap = mp
	ms2.InstanceID = "ins2"

	mss = append(mss, ms1)
	mss = append(mss, ms2)

	testRegistryObj := new(mk.DiscoveryMock)
	registry.DefaultServiceDiscoveryService = testRegistryObj
	testRegistryObj.On("FindMicroServiceInstances", "selfServiceID", "appID", "service1", "1.0", "").Return(mss, nil)
	config.GlobalDefinition = &chassisModel.GlobalCfg{}
	config.GetLoadBalancing().Strategy = make(map[string]string)
	loadbalancer.Enable()
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
		Strategy:         loadbalancer.StrategyRoundRobin,
		SourceServiceID:  "selfServiceID",
		RouteTags:        utiltags.NewDefaultTag("1.0", "appID"),
	}
	c.Next(i, func(r *invocation.Response) error {
		assert.NoError(t, r.Err)
		return r.Err
	})

	var lbh = new(handler.LBHandler)
	str := lbh.Name()
	assert.Equal(t, "loadbalancer", str)
	t.Log(i.Protocol)
	t.Log(i.Endpoint)
}

func BenchmarkLBHandler_Handle(b *testing.B) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "client"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	control.Init()
	registry.Enable()
	registry.DoRegister()
	loadbalancer.Enable()
	testData1 := []*registry.MicroService{
		{
			ServiceName: "test2",
			AppID:       "default",
			Version:     "1.0",
			Status:      "UP",
		},
	}
	testData2 := []*registry.MicroServiceInstance{
		{
			HostName:     "test1",
			Status:       "UP",
			EndpointsMap: map[string]string{"highway": "10.0.0.4:1234"},
		},
		{
			HostName:     "test2",
			Status:       "UP",
			EndpointsMap: map[string]string{"highway": "10.0.0.3:1234"},
		},
	}
	sid, _, _ := registry.DefaultRegistrator.RegisterServiceAndInstance(testData1[0], testData2[0])
	_, _, _ = registry.DefaultRegistrator.RegisterServiceAndInstance(testData1[0], testData2[1])
	c := handler.Chain{}
	c.AddHandler(&handler.LBHandler{})
	c.AddHandler(&handler1{})
	runtime.ServiceID = sid
	iv := &invocation.Invocation{
		MicroServiceName: "test2",
		Protocol:         "highway",
		Strategy:         loadbalancer.StrategyRoundRobin,
		SourceServiceID:  runtime.ServiceID,
		RouteTags:        utiltags.NewDefaultTag("1.0", "default"),
	}

	b.Log(runtime.ServiceID)
	time.Sleep(1 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Next(iv, func(r *invocation.Response) error {
			return r.Err
		})
		c.Reset()
	}

}
