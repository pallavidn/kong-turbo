package conf

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/service"
	"io/ioutil"
)

type KongTurboServiceSpec struct {
	*service.TurboCommunicationConfig `json:"communicationConfig,omitempty"`
	*KongTurboTargetConf              `json:"kongturboTargetConfig,omitempty"`
}

type KongTurboTargetConf struct {
	ProbeCategory string `json:"probeCategory,omitempty"`
	//TargetType       string `json:"targetType,omitempty"`
	TargetAddress string `json:"kongAdmin,omitempty"`
}

func NewKongTurboServiceSpec(configFilePath string) (*KongTurboServiceSpec, error) {

	glog.Infof("Read configuration from %s", configFilePath)
	tapSpec, err := readConfig(configFilePath)

	if err != nil {
		return nil, err
	}

	if tapSpec.TurboCommunicationConfig == nil {
		return nil, fmt.Errorf("Unable to read the turbo communication config from %s", configFilePath)
	}

	if tapSpec.KongTurboTargetConf == nil {
		return nil, fmt.Errorf("Unable to read the target config from %s", configFilePath)
	}

	return tapSpec, nil
}

func readConfig(path string) (*KongTurboServiceSpec, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Errorf("File error: %v\n", err)
		return nil, err
	}
	glog.Infoln(string(file))

	var spec KongTurboServiceSpec
	err = json.Unmarshal(file, &spec)

	if err != nil {
		glog.Errorf("Unmarshall error :%v\n", err)
		return nil, err
	}
	glog.Infof("Results: %+v\n", spec)

	return &spec, nil
}
