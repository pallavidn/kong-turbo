package discovery

import "github.com/turbonomic/turbo-go-sdk/pkg/proto"
import (
	"fmt"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
)

const (
	StitchingAttr            string = "GatewayVappId"	//"VappIds"	//"vAppUuid"
	FunctionIdAttr            string = "GatewayVappId"
	FunctionEndpointAttr            string = "FunctionEndpoint"

	LocalNameAttr    string = "LocalName"
	AltNameAttr      string = "altName"
	ExternalNameAttr string = "externalnames"

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
	var vappId, localId, altId string
	if len(apiSvc.Functions) > 0 {
		function := apiSvc.Functions[0]
		localId = function.FunctionName
	} else {
		localId = route.FunctionHost
	}

	if localId == "" {
		return nil, fmt.Errorf("Cannot create function vapp without ID %++v", route)
	}

	altId = route.Path	//functionEndpoint
	fmt.Printf("**** local name: %s\n", localId)
	fmt.Printf("**** alt name: %s\n", altId)

	// uuid and display name
	vappId = fmt.Sprintf("%s/%s/%s", "kong", route.Path, localId)
	fmt.Printf("**** vappUuid : %s\n", vappId)

	commodities := []*proto.CommodityDTO{}
	commodity, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_TRANSACTION).Create()
	commodities = append(commodities, commodity)

	vAppType := proto.EntityDTO_VIRTUAL_APPLICATION
	entityDTOBuilder := builder.NewEntityDTOBuilder(proto.EntityDTO_VIRTUAL_APPLICATION, vappId).
		DisplayName(vappId).
		//WithProperty(getEntityProperty(FunctionIdAttr, vappId)).
		//WithProperty(getEntityProperty(FunctionEndpointAttr, functionEndpoint)).
		WithProperty(getEntityProperty(LocalNameAttr, localId)).
		WithProperty(getEntityProperty(AltNameAttr, altId)).
		ReplacedBy(getReplacementMetaData(vAppType)).
		SellsCommodities(commodities)

	return entityDTOBuilder, nil
}

func getReplacementMetaData(entityType proto.EntityDTO_EntityType,
) *proto.EntityDTO_ReplacementEntityMetaData {
	extAttr := ExternalNameAttr	//StitchingAttr
	intAttr := LocalNameAttr
	useTopoExt := true

	b := builder.NewReplacementEntityMetaDataBuilder().
		Matching(intAttr).
		MatchingExternal(&proto.ServerEntityPropDef{
			Entity:    &entityType,
			Attribute: &extAttr,
			UseTopoExt: &useTopoExt,
		})
		//PatchBuyingWithProperty(proto.CommodityDTO_TRANSACTION, []string{PropertyUsed}).
		//PatchBuyingWithProperty(proto.CommodityDTO_RESPONSE_TIME, []string{PropertyUsed}).
		//PatchSellingWithProperty(proto.CommodityDTO_TRANSACTION, []string{PropertyUsed, PropertyCapacity}).
		//PatchSellingWithProperty(proto.CommodityDTO_RESPONSE_TIME, []string{PropertyUsed, PropertyCapacity})

	return b.Build()
}

func getEntityProperty(attr, value string) *proto.EntityDTO_EntityProperty {
	ns := DefaultPropertyNamespace

	return &proto.EntityDTO_EntityProperty{
		Namespace: &ns,
		Name:      &attr,
		Value:     &value,
	}
}

type FunctionDTOBuilder struct {
}
