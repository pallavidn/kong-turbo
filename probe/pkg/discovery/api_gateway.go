package discovery

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type APIGatewayInterface interface {
	GetServices() ([]*APIService, error)
	GetRoutes() ([]*FunctionRoute, error)
}

type APIService struct {
	ServiceName  string
	ServiceId    string
	HostName     string
	FunctionName string
}

type FunctionRoute struct {
	Path         string
	ServiceId    string
	FunctionHost string
	RouteId      string
}

type Function struct {
	Method   string
	Resource string
}

type Service struct {
	Name      string
	Functions []*Function
}

func DiscoverAPIGateway(gatewayAddr string) ([]*proto.EntityDTO, error) {

	kongClient, err := NewKongHttpClient(gatewayAddr)

	if err != nil {
		return nil, fmt.Errorf("Error creating client for Kong API Gateway object: %v", err)

	}

	serviceDtoBuilder := APIServiceDTOyBuilder{}

	// Fetch services
	apiServices, _ := kongClient.GetServices()
	apiServicesMap := make(map[string]*APIService)
	for _, apiService := range apiServices {
		apiServicesMap[apiService.ServiceId] = apiService
	}

	// Fetch api endpoint routes
	functionRoutes, _ := kongClient.GetRoutes()

	var discoveryResult []*proto.EntityDTO
	for _, functionRoute := range functionRoutes {
		apiService := apiServicesMap[functionRoute.ServiceId]

		dtoBuilder, err := serviceDtoBuilder.buildDto(apiService, functionRoute)
		if err != nil {
			glog.Errorf("%s", err)
			fmt.Printf("Error while building entity : %v\n", err)
		}
		if dtoBuilder == nil {
			fmt.Printf("%v\n", err)
		}
		dto, _ := dtoBuilder.Create()
		fmt.Printf("DTO %++v\n", dto)
		discoveryResult = append(discoveryResult, dto)
	}

	return discoveryResult, nil
}
