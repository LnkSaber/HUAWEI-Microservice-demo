/*
 * Copyright 2017 Huawei Technologies Co., Ltd
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

//Package configcentersource created on 2017/6/22.
package configcentersource

import (
	"crypto/tls"
	"errors"
	"sync"
	"time"

	"github.com/go-chassis/go-archaius/core"
	"github.com/go-chassis/go-cc-client/configcenter-client"

	"fmt"
	"github.com/go-chassis/go-cc-client"
	"github.com/go-chassis/go-cc-client/serializers"
	"github.com/go-mesh/openlogging"
	"github.com/gorilla/websocket"
	"net/url"
	"os"
	"reflect"
	"strings"
)

const (
	defaultTimeout = 10 * time.Second
	//ConfigCenterSourceConst variable of type string
	ConfigCenterSourceConst    = "ConfigCenterSource"
	configCenterSourcePriority = 0
	dimensionsInfo             = "dimensionsInfo"
	dynamicConfigAPI           = `/configuration/refresh/items`
	getConfigAPI               = `/configuration/items`
	defaultContentType         = "application/json"
	//ConfigServerMemRefreshError is error message
	ConfigServerMemRefreshError = "error in poulating config server member"
)

var (
	//ConfigPath is a variable of type string
	ConfigPath = ""
	//ConfigRefreshPath is a variable of type string
	ConfigRefreshPath = ""
)

//Handler handles configs from config center
type Handler struct {
	MemberDiscovery              configcenterclient.MemberDiscovery
	dynamicConfigHandler         *DynamicConfigHandler
	dimensionsInfo               string
	dimensionInfoMap             map[string]string
	Configurations               map[string]interface{}
	dimensionsInfoConfiguration  map[string]map[string]interface{}
	dimensionsInfoConfigurations []map[string]map[string]interface{}
	initSuccess                  bool
	connsLock                    sync.Mutex
	sync.RWMutex
	TLSClientConfig *tls.Config
	TenantName      string
	RefreshMode     int
	RefreshInterval time.Duration
	Version         string
	RefreshPort     string
	Environment     string
}

//ConfigCenterConfig is pointer of config center source
var ConfigCenterConfig *Handler

//NewConfigCenterSource initializes all components of configuration center
func NewConfigCenterSource(memberDiscovery configcenterclient.MemberDiscovery, dimInfo string, tlsConfig *tls.Config, tenantName string,
	refreshMode, refreshInterval int, enableSSL bool, version, refreshPort, env string) core.ConfigSource {

	if ConfigCenterConfig == nil {
		ConfigCenterConfig = new(Handler)
		ConfigCenterConfig.MemberDiscovery = memberDiscovery
		ConfigCenterConfig.dimensionsInfo = dimInfo
		ConfigCenterConfig.initSuccess = true
		ConfigCenterConfig.TLSClientConfig = tlsConfig
		ConfigCenterConfig.TenantName = tenantName
		ConfigCenterConfig.RefreshMode = refreshMode
		ConfigCenterConfig.RefreshInterval = time.Second * time.Duration(refreshInterval)
		ConfigCenterConfig.Version = version
		ConfigCenterConfig.RefreshPort = refreshPort
		ConfigCenterConfig.Environment = env

		//Read the version for yaml file
		//Set Default api version to V3
		var apiVersion string
		switch version {
		case "v2":
			apiVersion = "v2"
		case "V2":
			apiVersion = "v2"
		case "v3":
			apiVersion = "v3"
		case "V3":
			apiVersion = "v3"
		default:
			apiVersion = "v3"
		}
		//Update the API Base Path based on the Version
		updateAPIPath(apiVersion)

	}
	return ConfigCenterConfig
}

//Update the Base PATH and HEADERS Based on the version of Configcenter used.
func updateAPIPath(apiVersion string) {

	//Check for the env Name in Container to get Domain Name
	//Default value is  "default"
	projectID, isExsist := os.LookupEnv("cse.config.client.tenantName")
	if !isExsist {
		projectID = "default"
	}
	switch apiVersion {
	case "v3":
		ConfigPath = "/v3/" + projectID + getConfigAPI
		ConfigRefreshPath = "/v3/" + projectID + dynamicConfigAPI
	case "v2":
		ConfigPath = "/configuration/v2/items"
		ConfigRefreshPath = "/configuration/v2/refresh/items"
	default:
		ConfigPath = "/v3/" + projectID + getConfigAPI
		ConfigRefreshPath = "/v3/" + projectID + dynamicConfigAPI
	}
}

//GetConfigAPI is map
type GetConfigAPI map[string]map[string]interface{}

//CreateConfigAPI creates a configuration API
type CreateConfigAPI struct {
	DimensionInfo string                 `json:"dimensionsInfo"`
	Items         map[string]interface{} `json:"items"`
}

// ensure to implement config source
var _ core.ConfigSource = &Handler{}

//GetConfigurations gets a particular configuration
func (cfgSrcHandler *Handler) GetConfigurations() (map[string]interface{}, error) {
	configMap := make(map[string]interface{})

	err := cfgSrcHandler.refreshConfigurations("")
	if err != nil {
		return nil, err
	}
	if cfgSrcHandler.RefreshMode == 1 {
		go cfgSrcHandler.refreshConfigurationsPeriodically("")
	}

	cfgSrcHandler.Lock()
	for key, value := range cfgSrcHandler.Configurations {
		configMap[key] = value
	}
	cfgSrcHandler.Unlock()
	return configMap, nil
}

//GetConfigurationsByDI gets required configurations for particular dimension info
func (cfgSrcHandler *Handler) GetConfigurationsByDI(dimensionInfo string) (map[string]interface{}, error) {
	configMap := make(map[string]interface{})

	err := cfgSrcHandler.refreshConfigurations(dimensionInfo)
	if err != nil {
		return nil, err
	}

	if cfgSrcHandler.RefreshMode == 1 {
		go cfgSrcHandler.refreshConfigurationsPeriodically(dimensionInfo)
	}

	cfgSrcHandler.Lock()
	for key, value := range cfgSrcHandler.dimensionsInfoConfiguration {
		configMap[key] = value
	}
	cfgSrcHandler.Unlock()
	return configMap, nil
}

func (cfgSrcHandler *Handler) refreshConfigurationsPeriodically(dimensionInfo string) {
	ticker := time.Tick(cfgSrcHandler.RefreshInterval)
	isConnectionFailed := false
	for range ticker {
		err := cfgSrcHandler.refreshConfigurations(dimensionInfo)
		if err == nil {
			if isConnectionFailed {
				openlogging.GetLogger().Infof("Recover configurations from config center server")
			}
			isConnectionFailed = false
		} else {
			isConnectionFailed = true
		}
	}
}

func (cfgSrcHandler *Handler) refreshConfigurations(dimensionInfo string) error {
	var (
		config     map[string]interface{}
		configByDI map[string]map[string]interface{}
		err        error
		events     []*core.Event
	)

	if dimensionInfo == "" {
		config, err = client.DefaultClient.PullConfigs(cfgSrcHandler.dimensionsInfo, "", "", "")
		if err != nil {
			openlogging.GetLogger().Warnf("Failed to pull configurations from config center server", err) //Warn
			return err
		}
		//Populate the events based on the changed value between current config and newly received Config
		events, err = cfgSrcHandler.populateEvents(config)
	} else {
		var diInfo string
		for _, value := range cfgSrcHandler.dimensionInfoMap {
			if value == dimensionInfo {
				diInfo = dimensionInfo
			}
		}
		configByDI, err = client.DefaultClient.PullConfigsByDI(dimensionInfo, diInfo)
		if err != nil {
			openlogging.GetLogger().Warnf("Failed to pull configurations from config center server", err) //Warn
			return err
		}
		//Populate the events based on the changed value between current config and newly received Config based dimension info
		events, err = cfgSrcHandler.setKeyValueByDI(configByDI, dimensionInfo)
	}

	if err != nil {
		openlogging.GetLogger().Warnf("error in generating event", err)
		return err
	}

	//Generate OnEvent Callback based on the events created
	if cfgSrcHandler.dynamicConfigHandler != nil {
		openlogging.GetLogger().Debugf("event On Receive %+v", events)
		for _, event := range events {
			cfgSrcHandler.dynamicConfigHandler.EventHandler.Callback.OnEvent(event)
		}
	}

	cfgSrcHandler.Lock()
	cfgSrcHandler.updatedimensionsInfoConfigurations(dimensionInfo, configByDI, config)
	cfgSrcHandler.Unlock()

	return nil
}

func (cfgSrcHandler *Handler) updatedimensionsInfoConfigurations(dimensionInfo string,
	configByDI map[string]map[string]interface{}, config map[string]interface{}) {

	if dimensionInfo == "" {
		cfgSrcHandler.Configurations = config

	} else {
		if len(cfgSrcHandler.dimensionsInfoConfigurations) != 0 {
			for _, j := range cfgSrcHandler.dimensionsInfoConfigurations {
				// This condition is used to add the information of dimension info if there are 2 dimension
				if len(j) == 0 {
					cfgSrcHandler.dimensionsInfoConfigurations = append(cfgSrcHandler.dimensionsInfoConfigurations, configByDI)
				}
				for p := range j {
					if (p != dimensionInfo && len(cfgSrcHandler.dimensionInfoMap) > len(cfgSrcHandler.dimensionsInfoConfigurations)) || (len(j) == 0) {
						cfgSrcHandler.dimensionsInfoConfigurations = append(cfgSrcHandler.dimensionsInfoConfigurations, configByDI)
					}
					_, ok := j[dimensionInfo]
					if ok {
						delete(j, dimensionInfo)
						cfgSrcHandler.dimensionsInfoConfigurations = append(cfgSrcHandler.dimensionsInfoConfigurations, configByDI)
					}
				}
			}
			// This for loop to remove the emty map "map[]" from cfgSrcHandler.dimensionsInfoConfigurations
			for i, v := range cfgSrcHandler.dimensionsInfoConfigurations {
				if len(v) == 0 && len(cfgSrcHandler.dimensionsInfoConfigurations) > 1 {
					cfgSrcHandler.dimensionsInfoConfigurations = append(cfgSrcHandler.dimensionsInfoConfigurations[:i], cfgSrcHandler.dimensionsInfoConfigurations[i+1:]...)
				}
			}
		} else {
			cfgSrcHandler.dimensionsInfoConfigurations = append(cfgSrcHandler.dimensionsInfoConfigurations, configByDI)
		}

	}
}

//GetConfigurationByKey gets required configuration for a particular key
func (cfgSrcHandler *Handler) GetConfigurationByKey(key string) (interface{}, error) {
	cfgSrcHandler.Lock()
	configSrcVal, ok := cfgSrcHandler.Configurations[key]
	cfgSrcHandler.Unlock()
	if ok {
		return configSrcVal, nil
	}

	return nil, errors.New("key not exist")
}

//GetConfigurationByKeyAndDimensionInfo gets required configuration for a particular key and dimension pair
func (cfgSrcHandler *Handler) GetConfigurationByKeyAndDimensionInfo(key, dimensionInfo string) (interface{}, error) {
	var (
		configSrcVal interface{}
		actualValue  interface{}
		exist        bool
	)

	cfgSrcHandler.Lock()
	for _, v := range cfgSrcHandler.dimensionsInfoConfigurations {
		value, ok := v[dimensionInfo]
		if ok {
			actualValue, exist = value[key]
		}
	}
	cfgSrcHandler.Unlock()

	if exist {
		configSrcVal = actualValue
		return configSrcVal, nil
	}

	return nil, errors.New("key not exist")
}

//AddDimensionInfo adds dimension info for a configuration
func (cfgSrcHandler *Handler) AddDimensionInfo(dimensionInfo string) (map[string]string, error) {
	if len(cfgSrcHandler.dimensionInfoMap) == 0 {
		cfgSrcHandler.dimensionInfoMap = make(map[string]string)
	}

	for i := range cfgSrcHandler.dimensionInfoMap {
		if i == dimensionInfo {
			openlogging.GetLogger().Errorf("dimension info already exist")
			return cfgSrcHandler.dimensionInfoMap, errors.New("dimension info allready exist")
		}
	}

	cfgSrcHandler.dimensionInfoMap[dimensionInfo] = dimensionInfo

	return cfgSrcHandler.dimensionInfoMap, nil
}

//GetSourceName returns name of the configuration
func (*Handler) GetSourceName() string {
	return ConfigCenterSourceConst
}

//GetPriority returns priority of a configuration
func (*Handler) GetPriority() int {
	return configCenterSourcePriority
}

//DynamicConfigHandler dynamically handles a configuration
func (cfgSrcHandler *Handler) DynamicConfigHandler(callback core.DynamicConfigCallback) error {
	if cfgSrcHandler.initSuccess != true {
		return errors.New("config center source initialization failed")
	}

	dynCfgHandler, err := newDynConfigHandlerSource(cfgSrcHandler, callback)
	if err != nil {
		openlogging.GetLogger().Error("failed to initialize dynamic config center Handler:" + err.Error())
		return errors.New("failed to initialize dynamic config center Handler")
	}
	cfgSrcHandler.dynamicConfigHandler = dynCfgHandler

	if cfgSrcHandler.RefreshMode == 0 {
		// Pull All the configuration for the first time.
		cfgSrcHandler.refreshConfigurations("")
		//Start a web socket connection to recieve change events.
		dynCfgHandler.startDynamicConfigHandler(cfgSrcHandler.RefreshPort)
	}

	return nil
}

//Cleanup cleans the particular configuration up
func (cfgSrcHandler *Handler) Cleanup() error {
	cfgSrcHandler.connsLock.Lock()
	defer cfgSrcHandler.connsLock.Unlock()

	if cfgSrcHandler.dynamicConfigHandler != nil {
		cfgSrcHandler.dynamicConfigHandler.Cleanup()
	}

	cfgSrcHandler.dynamicConfigHandler = nil
	cfgSrcHandler.Configurations = nil

	return nil
}

func (cfgSrcHandler *Handler) populateEvents(updatedConfig map[string]interface{}) ([]*core.Event, error) {
	events := make([]*core.Event, 0)
	newConfig := make(map[string]interface{})
	cfgSrcHandler.Lock()
	defer cfgSrcHandler.Unlock()

	currentConfig := cfgSrcHandler.Configurations

	// generate create and update event
	for key, value := range updatedConfig {
		newConfig[key] = value
		currentValue, ok := currentConfig[key]
		if !ok { // if new configuration introduced
			events = append(events, cfgSrcHandler.constructEvent(core.Create, key, value))
		} else if !reflect.DeepEqual(currentValue, value) {
			events = append(events, cfgSrcHandler.constructEvent(core.Update, key, value))
		}
	}

	// generate delete event
	for key, value := range currentConfig {
		_, ok := newConfig[key]
		if !ok { // when old config not present in new config
			events = append(events, cfgSrcHandler.constructEvent(core.Delete, key, value))
		}
	}

	// update with latest config
	cfgSrcHandler.Configurations = newConfig

	return events, nil
}

func (cfgSrcHandler *Handler) setKeyValueByDI(updatedConfig map[string]map[string]interface{}, dimensionInfo string) ([]*core.Event, error) {
	events := make([]*core.Event, 0)
	newConfigForDI := make(map[string]map[string]interface{})
	cfgSrcHandler.Lock()
	defer cfgSrcHandler.Unlock()

	currentConfig := cfgSrcHandler.dimensionsInfoConfiguration

	// generate create and update event
	for key, value := range updatedConfig {
		if key == dimensionInfo {
			newConfigForDI[key] = value
			for k, v := range value {
				if len(currentConfig) == 0 {
					events = append(events, cfgSrcHandler.constructEvent(core.Create, k, v))
				}
				for diKey, val := range currentConfig {
					if diKey == dimensionInfo {
						currentValue, ok := val[k]
						if !ok { // if new configuration introduced
							events = append(events, cfgSrcHandler.constructEvent(core.Create, k, v))
						} else if currentValue != v {
							events = append(events, cfgSrcHandler.constructEvent(core.Update, k, v))
						}
					}
				}
			}
		}
	}

	// generate delete event
	for key, value := range currentConfig {
		if key == dimensionInfo {
			for k, v := range value {
				for _, val := range newConfigForDI {
					_, ok := val[k]
					if !ok {
						events = append(events, cfgSrcHandler.constructEvent(core.Delete, k, v))
					}
				}
			}
		}
	}

	// update with latest config
	cfgSrcHandler.dimensionsInfoConfiguration = newConfigForDI

	return events, nil
}

func (cfgSrcHandler *Handler) constructEvent(eventType string, key string, value interface{}) *core.Event {
	newEvent := new(core.Event)
	newEvent.EventSource = ConfigCenterSourceConst
	newEvent.EventType = eventType
	newEvent.Key = key
	newEvent.Value = value

	return newEvent
}

//DynamicConfigHandler is a struct
type DynamicConfigHandler struct {
	dimensionsInfo  string
	EventHandler    *ConfigCenterEventHandler
	dynamicLock     sync.Mutex
	wsDialer        *websocket.Dialer
	wsConnection    *websocket.Conn
	memberDiscovery configcenterclient.MemberDiscovery
}

func newDynConfigHandlerSource(cfgSrc *Handler, callback core.DynamicConfigCallback) (*DynamicConfigHandler, error) {
	eventHandler := newConfigCenterEventHandler(cfgSrc, callback)
	dynCfgHandler := new(DynamicConfigHandler)
	dynCfgHandler.EventHandler = eventHandler
	dynCfgHandler.dimensionsInfo = cfgSrc.dimensionsInfo
	dynCfgHandler.wsDialer = &websocket.Dialer{
		TLSClientConfig:  cfgSrc.TLSClientConfig,
		HandshakeTimeout: defaultTimeout,
	}
	dynCfgHandler.memberDiscovery = cfgSrc.MemberDiscovery
	return dynCfgHandler, nil
}

func (dynHandler *DynamicConfigHandler) getWebSocketURL(refreshPort string) (*url.URL, error) {

	var defaultTLS bool
	var parsedEndPoint []string
	var host string

	configCenterEntryPointList, err := dynHandler.memberDiscovery.GetConfigServer()
	if err != nil {
		openlogging.GetLogger().Error("error in member discovery:" + err.Error())
		return nil, err
	}
	activeEndPointList, err := dynHandler.memberDiscovery.GetWorkingConfigCenterIP(configCenterEntryPointList)
	if err != nil {
		openlogging.GetLogger().Error("failed to get ip list:" + err.Error())
	}
	for _, server := range activeEndPointList {
		parsedEndPoint = strings.Split(server, `://`)
		hostArr := strings.Split(parsedEndPoint[1], `:`)
		port := refreshPort
		if port == "" {
			port = "30104"
		}
		host = hostArr[0] + ":" + port
		if host == "" {
			host = "localhost"
		}
	}

	if dynHandler.wsDialer.TLSClientConfig != nil {
		defaultTLS = true
	}
	if host == "" {
		err := errors.New("host must be a URL or a host:port pair")
		openlogging.GetLogger().Error("empty host for watch action:" + err.Error())
		return nil, err
	}
	hostURL, err := url.Parse(host)
	if err != nil || hostURL.Scheme == "" || hostURL.Host == "" {
		scheme := "ws://"
		if defaultTLS {
			scheme = "wss://"
		}
		hostURL, err = url.Parse(scheme + host)
		if err != nil {
			return nil, err
		}
		if hostURL.Path != "" && hostURL.Path != "/" {
			return nil, fmt.Errorf("host must be a URL or a host:port pair: %q", host)
		}
	}
	return hostURL, nil
}

func (dynHandler *DynamicConfigHandler) startDynamicConfigHandler(refreshPort string) error {
	parsedDimensionInfo := strings.Replace(dynHandler.dimensionsInfo, "#", "%23", -1)
	refreshConfigPath := ConfigRefreshPath + `?` + dimensionsInfo + `=` + parsedDimensionInfo
	if dynHandler != nil && dynHandler.wsDialer != nil {
		/*-----------------
		1. Decide on the URL
		2. Create WebSocket Connection
		3. Call KeepAlive in seperate thread
		3. Generate events on Recieve Data
		*/
		baseURL, err := dynHandler.getWebSocketURL(refreshPort)
		if err != nil {
			error := errors.New("error in getting default server info")
			return error
		}
		url := baseURL.String() + refreshConfigPath
		dynHandler.dynamicLock.Lock()
		dynHandler.wsConnection, _, err = dynHandler.wsDialer.Dial(url, nil)
		if err != nil {
			dynHandler.dynamicLock.Unlock()
			return fmt.Errorf("watching config-center dial catch an exception error:%s", err.Error())
		}
		dynHandler.dynamicLock.Unlock()
		keepAlive(dynHandler.wsConnection, 15*time.Second)
		go func() error {
			for {
				messageType, message, err := dynHandler.wsConnection.ReadMessage()
				if err != nil {
					break
				}
				if messageType == websocket.TextMessage {
					dynHandler.EventHandler.OnReceive(message)
				}
			}
			err = dynHandler.wsConnection.Close()
			if err != nil {
				return fmt.Errorf("CC watch Conn close failed error:%s", err.Error())
			}
			return nil
		}()
	}
	return nil
}

func keepAlive(c *websocket.Conn, timeout time.Duration) {
	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})
	go func() {
		for {
			err := c.WriteMessage(websocket.PingMessage, []byte("keepalive"))
			if err != nil {
				return
			}
			time.Sleep(timeout / 2)
			if time.Now().Sub(lastResponse) > timeout {
				c.Close()
				return
			}
		}
	}()
}

//Cleanup cleans particular dynamic configuration Handler up
func (dynHandler *DynamicConfigHandler) Cleanup() error {
	dynHandler.dynamicLock.Lock()
	defer dynHandler.dynamicLock.Unlock()
	if dynHandler.wsConnection != nil {
		dynHandler.wsConnection.Close()
	}
	dynHandler.wsConnection = nil
	return nil
}

//ConfigCenterEventHandler handles a event of a configuration center
type ConfigCenterEventHandler struct {
	ConfigSource *Handler
	Callback     core.DynamicConfigCallback
}

//ConfigCenterEvent stores info about an configuration center event
type ConfigCenterEvent struct {
	Action string `json:"action"`
	Value  string `json:"value"`
}

func newConfigCenterEventHandler(cfgSrc *Handler, callback core.DynamicConfigCallback) *ConfigCenterEventHandler {
	eventHandler := new(ConfigCenterEventHandler)
	eventHandler.ConfigSource = cfgSrc
	eventHandler.Callback = callback
	return eventHandler
}

//OnConnect is a method
func (*ConfigCenterEventHandler) OnConnect() {
	return
}

//OnConnectionClose is a method
func (*ConfigCenterEventHandler) OnConnectionClose() {
	return
}

//OnReceive initializes all necessary components for a configuration center
func (eventHandler *ConfigCenterEventHandler) OnReceive(actionData []byte) {
	configCenterEvent := new(ConfigCenterEvent)
	err := serializers.Decode(serializers.JsonEncoder, actionData, &configCenterEvent)
	if err != nil {
		openlogging.GetLogger().Errorf(fmt.Sprintf("error in unmarshalling data on event receive with error %s", err.Error()))
		return
	}

	sourceConfig := make(map[string]interface{})
	err = serializers.Decode(serializers.JsonEncoder, []byte(configCenterEvent.Value), &sourceConfig)
	if err != nil {
		openlogging.GetLogger().Errorf(fmt.Sprintf("error in unmarshalling config values %s", err.Error()))
		return
	}

	events, err := eventHandler.ConfigSource.populateEvents(sourceConfig)
	if err != nil {
		openlogging.GetLogger().Error("error in generating event:" + err.Error())
		return
	}

	openlogging.GetLogger().Debugf("event On Receive", events)
	for _, event := range events {
		eventHandler.Callback.OnEvent(event)
	}

	return
}

//InitConfigCenter is a function which initializes the memberDiscovery of go-cc-client
func InitConfigCenter(ccEndpoint, dimensionInfo, tenantName string, enableSSL bool, tlsConfig *tls.Config, refreshMode int,
	refreshInterval int, autoDiscovery bool, clientType, apiVersion, refreshPort, environment string) (core.ConfigSource, error) {
	memDiscovery := configcenterclient.NewConfiCenterInit(tlsConfig, tenantName, enableSSL, apiVersion, autoDiscovery, environment)

	configCenters := strings.Split(ccEndpoint, ",")
	cCenters := make([]string, 0)
	for _, value := range configCenters {
		value = strings.Replace(value, " ", "", -1)
		cCenters = append(cCenters, value)
	}

	memDiscovery.ConfigurationInit(cCenters)

	if enabledAutoDiscovery(autoDiscovery) {
		refreshError := memDiscovery.RefreshMembers()
		if refreshError != nil {
			openlogging.GetLogger().Error(ConfigServerMemRefreshError + refreshError.Error())
			return nil, errors.New(ConfigServerMemRefreshError)
		}
	}

	configCenterSource := NewConfigCenterSource(
		memDiscovery, dimensionInfo, tlsConfig, tenantName, refreshMode,
		refreshInterval, enableSSL, apiVersion, refreshPort, environment)

	configcenterclient.MemberDiscoveryService = memDiscovery
	if err := installPlugin(clientType); err != nil {
		return nil, err
	}
	return configCenterSource, nil
}

func installPlugin(clientType string) error {
	return client.Enable(clientType)
}

func enabledAutoDiscovery(autoDiscovery bool) bool {
	if autoDiscovery {
		return true
	}
	return false
}