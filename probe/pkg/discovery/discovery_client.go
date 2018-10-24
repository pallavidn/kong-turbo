package discovery

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/pallavidn/kong-turbo/probe/pkg/registration"
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

// Implements the TurboDiscoveryClient interface
type KongDiscoveryClient struct {
	targetAddr string
}

func NewDiscoveryClient(targetAddr string) *KongDiscoveryClient {
	glog.V(2).Infof("New Discovery client with target address %s", targetAddr)
	return &KongDiscoveryClient{
		targetAddr: targetAddr,
	}
}

// Get the Account Values to create VMTTarget in the turbo server corresponding to this client
func (d *KongDiscoveryClient) GetAccountValues() *probe.TurboTargetInfo {
	targetIdFieldName := registration.TargetIdentifierField
	targetIdVal := &proto.AccountValue{
		Key:         &targetIdFieldName,
		StringValue: &d.targetAddr,
	}

	accountValues := []*proto.AccountValue{
		targetIdVal,
	}

	targetInfo := probe.NewTurboTargetInfoBuilder(registration.ProbeCategory, registration.TargetType,
		registration.TargetIdentifierField, accountValues).Create()

	return targetInfo
}

// Validate the Target
func (d *KongDiscoveryClient) Validate(accountValues []*proto.AccountValue) (*proto.ValidationResponse, error) {
	glog.V(2).Infof("Validating Kong target %s", accountValues)
	// TODO: Add logic for validation
	validationResponse := &proto.ValidationResponse{}

	// Validation fails if no exporter responses
	return validationResponse, nil
}

// Discover the Target Topology
func (d *KongDiscoveryClient) Discover(accountValues []*proto.AccountValue) (*proto.DiscoveryResponse, error) {
	glog.V(2).Infof("Discovering Kong target %s", accountValues)
	var entities []*proto.EntityDTO

	var discoveryResponse *proto.DiscoveryResponse
	entities, err := DiscoverAPIGateway(d.targetAddr)
	if err != nil {
		return d.failDiscovery(), nil
	}

	discoveryResponse = &proto.DiscoveryResponse {
		EntityDTO: entities,
	}

	return discoveryResponse, nil
}

func (d *KongDiscoveryClient) failDiscovery() *proto.DiscoveryResponse {
	description := fmt.Sprintf("Kongturbo probe discovery failed")
	glog.Errorf(description)
	severity := proto.ErrorDTO_CRITICAL
	errorDTO := &proto.ErrorDTO{
		Severity:    &severity,
		Description: &description,
	}
	discoveryResponse := &proto.DiscoveryResponse{
		ErrorDTO: []*proto.ErrorDTO{errorDTO},
	}
	return discoveryResponse
}

func (d *KongDiscoveryClient) failValidation() *proto.ValidationResponse {
	description := fmt.Sprintf("Kongturbo probe validation failed")
	glog.Errorf(description)
	severity := proto.ErrorDTO_CRITICAL
	errorDto := &proto.ErrorDTO{
		Severity:    &severity,
		Description: &description,
	}

	validationResponse := &proto.ValidationResponse{
		ErrorDTO: []*proto.ErrorDTO{errorDto},
	}
	return validationResponse
}
