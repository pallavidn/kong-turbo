package discovery

import "github.com/turbonomic/turbo-go-sdk/pkg/proto"
import (
	"fmt"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"strings"
)

const (
	StitchingAttr            string = "GatewayVappId"	//"VappIds"	//"vAppUuid"
	GatewayStitchingAttr            string = "GatewayVappId"
	DefaultPropertyNamespace string = "DEFAULT"
	PropertyUsed        = "used"
	PropertyCapacity    = "capacity"
)

type APIServiceDTOyBuilder struct {
}

func (dtoBuilder *APIServiceDTOyBuilder) buildDto(apiSvc *APIService, route *FunctionRoute) (*builder.EntityDTOBuilder, error) {
	if apiSvc == nil {
		return nil, fmt.Errorf("Null service for %++v", route)
	}

	if route == nil {
		return nil, fmt.Errorf("Null route")
	}

	// id.
	var vappId string
	if apiSvc.FunctionName != "" {
		vappId = apiSvc.FunctionName
		idx := strings.LastIndex(vappId, ":")
		if idx > -1 {
			vappId_new := vappId[0:idx]
			vappId = vappId_new
			fmt.Printf("### Changed vapp id : %s\n", vappId)
		}
	} else {
		vappId = route.FunctionHost
		//// function name and namespace
		//substr := strings.Split(vappId, ".")
		//if len(substr) > 2 {
		//	funcName := substr[0]
		//	namespace := substr[1]
		//	vappId = fmt.Sprintf("%s/%s-%s", namespace, funcName, "00001-service")
		//}
	}

	if vappId == "" {
		return nil, fmt.Errorf("Cannot create function vapp without ID %++v", route)
	}

	fmt.Printf("**** vapp id : %s\n", vappId)

	// uuis
	vappUuid := strings.Join( []string{"kong-vapp",vappId}, "-")

	// display name.
	displayName := strings.Join( []string{"kong-vapp",vappId}, "-")

	fmt.Printf("**** vapp vappUuid : %s\n", vappUuid)
	fmt.Printf("**** vapp displayName : %s\n", displayName)


	commodities := []*proto.CommodityDTO{}
	commodity, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_TRANSACTION).Create()
	commodities = append(commodities, commodity)

	vAppType := proto.EntityDTO_VIRTUAL_APPLICATION
	entityDTOBuilder := builder.NewEntityDTOBuilder(proto.EntityDTO_VIRTUAL_APPLICATION, vappUuid).
		DisplayName(displayName).
		WithProperty(getEntityProperty(vappId)).
		ReplacedBy(getReplacementMetaData(vAppType)).
		SellsCommodities(commodities)

	return entityDTOBuilder, nil
}

func getReplacementMetaData(entityType proto.EntityDTO_EntityType,
) *proto.EntityDTO_ReplacementEntityMetaData {
	extAttr := StitchingAttr
	intAttr := GatewayStitchingAttr
	useTopoExt := true

	b := builder.NewReplacementEntityMetaDataBuilder().
		Matching(intAttr).
		MatchingExternal(&proto.ServerEntityPropDef{
			Entity:    &entityType,
			Attribute: &extAttr,
			UseTopoExt: &useTopoExt,
		}).
		PatchBuyingWithProperty(proto.CommodityDTO_TRANSACTION, []string{PropertyUsed}).
		PatchBuyingWithProperty(proto.CommodityDTO_RESPONSE_TIME, []string{PropertyUsed}).
		PatchSellingWithProperty(proto.CommodityDTO_TRANSACTION, []string{PropertyUsed, PropertyCapacity}).
		PatchSellingWithProperty(proto.CommodityDTO_RESPONSE_TIME, []string{PropertyUsed, PropertyCapacity})

	return b.Build()
}

func getEntityProperty(value string) *proto.EntityDTO_EntityProperty {
	attr := GatewayStitchingAttr	//StitchingAttr
	ns := DefaultPropertyNamespace

	return &proto.EntityDTO_EntityProperty{
		Namespace: &ns,
		Name:      &attr,
		Value:     &value,
	}
}

type FunctionDTOBuilder struct {
}
