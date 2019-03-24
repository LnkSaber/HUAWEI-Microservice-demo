package servicecenter_test

import (
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	client "github.com/go-chassis/go-chassis/pkg/scclient"
	_ "github.com/go-chassis/go-chassis/security/plugins/plain"
	"github.com/stretchr/testify/assert"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegistryClient_Health(t *testing.T) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server"))
	config.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	registryClient := &client.RegistryClient{}
	err := registryClient.Initialize(
		client.Options{
			Addrs: []string{"127.0.0.1:30100"},
		})
	assert.NoError(t, err)
	instances, err := registryClient.Health()
	t.Log("testing health of SC, health response : ", instances)
	assert.NoError(t, err)
	assert.NotZero(t, len(instances))

	services, err := registryClient.GetAllResources("instances")
	assert.NoError(t, err)
	for _, service := range services {
		for _, inst := range service.Instances {
			for _, uri := range inst.Endpoints {
				u, err := url.Parse(uri)
				if err != nil {
					lager.Logger.Error("Wrong URI: " + err.Error())
					continue
				}
				u.Host = strings.Split(u.Host, ":")[0]
				t.Log(u.Host, service.MicroService)
				//no need to analyze each endpoint
				break
			}
		}
	}
}
