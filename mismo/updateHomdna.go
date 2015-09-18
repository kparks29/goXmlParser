package mismo

import (
	"encoding/base64"
)

func UpdateHomdnaModel(file *[]byte) (*AppraisalResponse, error) {

	result, err := ParseXml(file)
	if err != nil {
		return nil, err
	}

	address := CreateAddressModel(result)
	appliances := CreateApplianceModels(result.Property.Structure.KitchenAppliances)
	lotFeatures := CreateLotFeatureModels(result.Property.Structure.LotFeatures)
	structureFeatures := CreateStructureFeatureModels(result.Property.Structure.StructureFeatures)
	lotSize := GetSize(result.Property.LotInfo.LotSize)
	lot := CreateLot(&lotSize, lotFeatures)
	structures := CreateStructures(result.Property.Structure, structureFeatures, appliances)
	homdna := CreateHomdnaModel(address, lot, structures)

	document, err := base64.StdEncoding.DecodeString(result.Report.Document.File)
	if err != nil {
		return nil, err
	}

	appraisalResponse := &AppraisalResponse{
		Homdna:            homdna,
		AppraisalDocument: &document,
		MIMEType:          &result.Report.Document.MIMEType,
		DocumentName:      &result.Report.Document.Name,
		ParsedXML:         result,
	}

	return appraisalResponse, nil
}
