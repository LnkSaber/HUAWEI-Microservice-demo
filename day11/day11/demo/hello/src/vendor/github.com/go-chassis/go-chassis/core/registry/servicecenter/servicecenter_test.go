package servicecenter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/registry"
	_ "github.com/go-chassis/go-chassis/core/registry/servicecenter"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/scclient"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	_ "github.com/go-chassis/go-chassis/security/plugins/plain"
	"github.com/stretchr/testify/assert"
)

func TestServicecenter_RegisterServiceAndInstance(t *testing.T) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server"))
	t.Log("Test servercenter.go")
	config.Init()
	runtime.Init()
	t.Log(os.Getenv("CHASSIS_HOME"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	registry.Enable()
	registry.DoRegister()

	testRegisterServiceAndInstance(t, registry.DefaultRegistrator, registry.DefaultServiceDiscoveryService)
	sid := testGetMicroServiceID(t, "CSE", "DSFtestAppThree", "2.0.3", registry.DefaultServiceDiscoveryService)
	t.Log("获取依赖的实例")
	tags := utiltags.NewDefaultTag("2.0.3", "CSE")
	instances, err := registry.DefaultServiceDiscoveryService.FindMicroServiceInstances(sid, "DSFtestAppThree", tags)
	assert.NoError(t, err)
	assert.NotZero(t, len(instances))

	err = registry.DefaultRegistrator.AddSchemas(sid, "dsfapp.HelloHuawei", "Testschemainfo")
	assert.NoError(t, err)

	microservices, err := registry.DefaultServiceDiscoveryService.GetAllMicroServices()
	assert.NoError(t, err)
	assert.NotZero(t, len(microservices))
}

func testRegisterServiceAndInstance(t *testing.T, scc registry.Registrator, sd registry.ServiceDiscovery) {
	microservice := &registry.MicroService{
		AppID:       "CSE",
		ServiceName: "DSFtestAppThree",
		Version:     "2.0.3",
		Status:      client.MicorserviceUp,
		Level:       "FRONT",
		Schemas:     []string{"dsfapp.HelloHuawei"},
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]string{"rest": "10.146.207.197:8080"},
		HostName:     "default",
		Status:       client.MSInstanceUP,
	}
	sid, insID, err := scc.RegisterServiceAndInstance(microservice, microServiceInstance)
	assert.NoError(t, err)
	t.Log("test update")
	scc.UpdateMicroServiceProperties(sid, map[string]string{"test": "test"})
	microservice2, err := sd.GetMicroService(sid)
	assert.NoError(t, err)
	assert.Equal(t, "test", microservice2.Metadata["test"])

	success, err := scc.Heartbeat(sid, insID)
	assert.Equal(t, success, true)
	assert.NoError(t, err)

	_, err = scc.Heartbeat("jdfhbh", insID)
	assert.Error(t, err)

	ins, err := sd.GetMicroServiceInstances(sid, sid)
	assert.NotZero(t, len(ins))
	assert.NoError(t, err)

	err = scc.UpdateMicroServiceInstanceStatus(sid, insID, "UP")
	assert.NoError(t, err)

	err = scc.UpdateMicroServiceInstanceProperties(sid, insID, map[string]string{"test": "test"})
	assert.NoError(t, err)

	msdep := &registry.MicroServiceDependency{
		Consumer:  &registry.MicroService{AppID: "CSE", ServiceName: "DSFtestAppThree", Version: "2.0.3"},
		Providers: []*registry.MicroService{{AppID: "CSE", ServiceName: "DSFtestAppThree", Version: "2.0.3"}},
	}
	scc.AddDependencies(msdep)

	heartBeatSvc := registry.HeartbeatService{}
	heartBeatSvc.RefreshTask(sid, insID)
	heartBeatSvc.RemoveTask(sid, insID)
	heartBeatSvc.Stop()
	scc.Close()
}

func testGetMicroServiceID(t *testing.T, appID, microServiceName, version string, sd registry.ServiceDiscovery) string {
	sid, err := sd.GetMicroServiceID(appID, microServiceName, version, "")
	assert.Nil(t, err)
	return sid
}
