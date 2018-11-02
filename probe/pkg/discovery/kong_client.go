package discovery

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"strconv"
	"strings"
)

const (
	DEFAULT_HTTP_PROTOCOL   = "http"
	DEFAULT_KONG_ADMIN_PORT = 8001
)

type KongHttpClient struct {
	httpClient *HttpClient
	hostname   string
}

func NewKongHttpClient(hostname string) (*KongHttpClient, error) {
	hostUrl := fmt.Sprintf(DEFAULT_HTTP_PROTOCOL + "://" + hostname + ":" + strconv.Itoa(DEFAULT_KONG_ADMIN_PORT))
	fmt.Printf("Http host url: %s\n", hostUrl)
	httpClient, err := NewHttpClient(hostUrl)
	if err != nil {
		return nil, err
	}
	return &KongHttpClient{
		hostname:   hostname,
		httpClient: httpClient,
	}, nil
}

func (kongClient *KongHttpClient) GetServices() ([]*APIService, error) {
	hclient := kongClient.httpClient

	var resp []byte
	resp, err := hclient.DoGet("services")
	if err != nil {
		glog.Errorf("Error: %v", err)
		return nil, err
	}

	var kongServiceList []*APIService

	var serviceData interface{}
	err = json.Unmarshal(resp, &serviceData)
	if err != nil {
		glog.Errorf("JSON error %s", err)
		return nil, fmt.Errorf("Error in json unmarshal for stats response : %s ", err)
	}

	serviceDataMap, isMap := serviceData.(map[string]interface{})
	if isMap {
		var serviceList []interface{}
		serviceList, exists := serviceDataMap["data"].([]interface{})
		if exists {
			for _, tgt := range serviceList {
				serviceInstance, ok := tgt.(map[string]interface{})
				if !ok {
					continue
				}

				var serviceName, serviceId, hostName string
				serviceName, _ = serviceInstance["name"].(string)
				serviceId, _ = serviceInstance["id"].(string)
				hostName, _ = serviceInstance["host"].(string)

				kongService := &APIService{
					ServiceName: serviceName,
					HostName:    hostName,
					ServiceId:   serviceId,
				}
				_ = kongClient.getServicePlugins(kongService)
				//fmt.Printf("**** Found service - %s::%s::%s\n", serviceId, serviceName, hostName)
				kongServiceList = append(kongServiceList, kongService)
				fmt.Printf("**** Service - %++v\n", kongService)
			}
		}
	}

	return kongServiceList, nil
}

// Parse service plugin for Function data
func (kongClient *KongHttpClient) getServicePlugins(kongSvc *APIService) error {
	hclient := kongClient.httpClient

	var resp []byte
	routePath := fmt.Sprintf("services/%s/plugins", kongSvc.ServiceId)
	resp, err := hclient.DoGet(routePath)
	if err != nil {
		glog.Errorf("Error: %v", err)
		return err
	}

	var pluginData interface{}
	err = json.Unmarshal(resp, &pluginData)
	if err != nil {
		glog.Errorf("JSON error %s", err)
		return fmt.Errorf("Error in json unmarshal for stats response : %s ", err)
	}

	var functionName string
	pluginDataMap, isMap := pluginData.(map[string]interface{})
	if isMap {
		var pluginList []interface{}
		pluginList, exists := pluginDataMap["data"].([]interface{})
		if exists {
			for _, tgt := range pluginList {
				pluginInstance, ok := tgt.(map[string]interface{})
				if !ok {
					continue
				}
				name, _ := pluginInstance["name"]
				if name != "aws-lambda" {
					continue
				}
				//id, _ := pluginInstance["service_id"]
				//fmt.Printf("*** plugin for service %s:%s\n", id, name)
				// Look for host headers
				configMap, isConfigMap := pluginInstance["config"].(map[string]interface{})
				if isConfigMap {
					functionName, _ = configMap["function_name"].(string)
					if strings.HasPrefix(functionName, "arn:aws") {
						if strings.HasSuffix(functionName, ":$LATEST") {
							idx := strings.LastIndex(functionName, ":")
							if idx > -1 {
								temp := functionName[0:idx]
								functionName = temp
								fmt.Printf("### Changed function : %s\n", functionName)
							}
						}
					}
					break
				}
			}
		}
	}

	if functionName == "" {
		return fmt.Errorf("Function name not found for service %s", kongSvc.ServiceName)
	}
	function := &Function{
		FunctionName: functionName,
	}
	kongSvc.Functions = append(kongSvc.Functions, function)
	return nil
}

func (kongClient *KongHttpClient) GetRoutes() ([]*FunctionRoute, error) {
	hclient := kongClient.httpClient

	var resp []byte
	resp, err := hclient.DoGet("routes")
	if err != nil {
		glog.Errorf("Error: %v", err)
		return nil, err
	}

	var kongRouteList []*FunctionRoute

	var routeData interface{}
	err = json.Unmarshal(resp, &routeData)
	if err != nil {
		glog.Errorf("JSON error %s", err)
		return nil, fmt.Errorf("Error in json unmarshal for stats response : %s ", err)
	}

	routeDataMap, isMap := routeData.(map[string]interface{})
	if isMap {
		var routeList []interface{}
		routeList, exists := routeDataMap["data"].([]interface{})
		if exists {
			for _, tgt := range routeList {
				routeInstance, ok := tgt.(map[string]interface{})
				if !ok {
					continue
				}

				service, _ := routeInstance["service"].(map[string]interface{})
				serviceId, _ := service["id"].(string)
				paths, _ := routeInstance["paths"].([]interface{})
				path, _ := paths[0].(string)
				id, _ := routeInstance["id"].(string)

				kongRoute := &FunctionRoute{
					RouteId:   id,
					ServiceId: serviceId,
					Path:      path[1:],
				}

				//for _, path := range paths {
				//	fmt.Printf("**** path - %v\n", path)
				//}
				_ = kongClient.getRoutePlugins(kongRoute)
				kongRouteList = append(kongRouteList, kongRoute)

				fmt.Printf("**** Route - %++v\n", kongRoute)
			}
		}
	}

	return kongRouteList, nil
}

// Parse route plugin for Function data
func (kongClient *KongHttpClient) getRoutePlugins(kongRoute *FunctionRoute) error {
	hclient := kongClient.httpClient

	var resp []byte
	routePath := fmt.Sprintf("routes/%s/plugins", kongRoute.RouteId)
	resp, err := hclient.DoGet(routePath)
	if err != nil {
		glog.Errorf("Error: %v", err)
		return err
	}

	var pluginData interface{}
	err = json.Unmarshal(resp, &pluginData)
	if err != nil {
		glog.Errorf("JSON error %s", err)
		return fmt.Errorf("Error in json unmarshal for stats response : %s ", err)
	}

	var functionRoute string
	pluginDataMap, isMap := pluginData.(map[string]interface{})
	if isMap {
		var pluginList []interface{}
		pluginList, exists := pluginDataMap["data"].([]interface{})
		if exists {
			for _, tgt := range pluginList {

				pluginInstance, ok := tgt.(map[string]interface{})
				if !ok {
					continue
				}
				name, _ := pluginInstance["name"]
				if name != "request-transformer" {
					continue
				}
				//id, _ := pluginInstance["route_id"]
				//fmt.Printf("*** plugin for route %s:%s\n", id, name)
				configMap, isConfigMap := pluginInstance["config"].(map[string]interface{})
				if isConfigMap {
					appendHeadersMap, isHeaderMap := configMap["append"].(map[string]interface{})
					if isHeaderMap {
						headers := appendHeadersMap["headers"]
						headersList, isHeadersList := headers.([]interface{})
						if isHeadersList {
							for _, header := range headersList {
								headerStr := header.(string)
								fmt.Printf("**** header - %v\n", headerStr)
								if strings.HasPrefix(headerStr, "Host:") {
									functionRoute = strings.TrimPrefix(headerStr, "Host:")
									break
								}
							}
						}
					}
				}
			}
		}
	}

	if functionRoute == "" {
		return fmt.Errorf("Function name not found for route %s:%s",
			kongRoute.ServiceId, kongRoute.RouteId)
	}

	kongRoute.FunctionHost = functionRoute
	return nil
}
