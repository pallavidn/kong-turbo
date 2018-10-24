package registration

import (
"github.com/turbonomic/turbo-go-sdk/pkg/proto"
"github.com/turbonomic/turbo-go-sdk/pkg/supplychain"
)

var (
	transactionType = proto.CommodityDTO_TRANSACTION
	key             = "key-placeholder"

	transactionTemplateComm *proto.TemplateCommodity = &proto.TemplateCommodity{
		CommodityType: &transactionType,
		//Key:           &key,
	}
)

type SupplyChainFactory struct{}

func (f *SupplyChainFactory) CreateSupplyChain() ([]*proto.TemplateDTO, error) {
	vAppNode, err := f.buildVAppSupplyBuilder()
	if err != nil {
		return nil, err
	}

	return supplychain.NewSupplyChainBuilder().Top(vAppNode).
		Create()
}

func (f *SupplyChainFactory) buildVAppSupplyBuilder() (*proto.TemplateDTO, error) {
	builder := supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_VIRTUAL_APPLICATION).
		Sells(transactionTemplateComm)
	builder.SetPriority(-1000)
	//builder.SetTemplateType(proto.TemplateDTO_BASE)
	//builder.SetTemplateType(proto.TemplateDTO_EXTENSION)

	return builder.Create()
}
