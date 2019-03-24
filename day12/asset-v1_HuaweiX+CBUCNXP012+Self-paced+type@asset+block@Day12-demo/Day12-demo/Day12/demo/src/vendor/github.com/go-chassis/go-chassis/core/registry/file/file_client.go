package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"

	"github.com/go-chassis/go-chassis/pkg/scclient/proto"
)

const (
	// ServiceJSON service json
	ServiceJSON = "service.json"
)

type localFileData struct {
	ServiceName  string `json:"serviceName,omitempty"`
	InstanceData *proto.MicroServiceInstance
}

// Options struct having addresses
type Options struct {
	Addrs []string
}

type fileClient struct {
	Addresses []string
}

type serviceData struct {
	Service []*service `json:"service,omitempty"`
}
type service struct {
	Name     string   `json:"name,omitempty"`
	Instance []string `json:"instance,omitempty"`
}

func (f *fileClient) Initialize(opt Options) {
	f.Addresses = opt.Addrs
}

func (f *fileClient) FindMicroServiceInstances(microServiceName string) ([]*proto.MicroServiceInstance, error) {
	var instanceData []*proto.MicroServiceInstance

	data := f.getInstanceDataFromFile()
	if data == nil {
		return instanceData, fmt.Errorf("failed to get instance information")
	}

	localData := &localFileData{}
	for _, value := range data.Service {
		if value.Name == microServiceName {
			insData := &proto.MicroServiceInstance{
				Endpoints: value.Instance,
			}
			localData.ServiceName = value.Name
			localData.InstanceData = insData

			instanceData = append(instanceData, localData.InstanceData)
			return instanceData, nil
		}
	}

	return instanceData, nil
}

func (f *fileClient) getInstanceDataFromFile() *serviceData {
	var data *serviceData
	path := strings.Join(f.Addresses, "")

	if path == "" {
		cwd, _ := fileutil.GetWorkDir()
		path = filepath.Join(cwd, "disco", ServiceJSON)
	}

	file, err := os.Open(path)
	if err != nil {
		lager.Logger.Warnf("failed to open a file", err)
	}
	defer file.Close()

	plan, err := ioutil.ReadFile(path)
	if err != nil {
		lager.Logger.Warnf("failed to do readfile operation", err)
	}

	err = json.Unmarshal(plan, &data)
	if err != nil {
		lager.Logger.Warnf("failed to do unmarshall", err)
	}

	return data
}
