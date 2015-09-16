package appraisal

import (
	"github.com/homdna/homdna-models"
)

func UpdateHomdnaModel(file *[]byte) (models.HomdnaModel, *[]byte, string, string, error) {

	result, err := ParseXml(file)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	address := appraisal.CreateAddressModel(result)
	appliances := appraisal.CreateApplianceModels(result.Property.Structure.KitchenAppliances)
	lotFeatures := appraisal.CreateLotFeatureModels(result.Property.Structure.LotFeatures)
	structureFeatures := appraisal.CreateStructureFeatureModels(result.Property.Structure.StructureFeatures)
	lotSize := appraisal.GetSize(result.Property.LotInfo.LotSize)
	lot := appraisal.CreateLot(&lotSize, lotFeatures)
	structures := appraisal.CreateStructures(result.Property.Structure, structureFeatures, appliances)
	homdna := appraisal.CreateHomdnaModel(address, lot, structures)

	if document, err := base64.StdEncoding.DecodeString(&result.Report.Document.File); err != nil {
		return nil, nil, nil, nil, err
	}

	return homdna, document, &result.Report.Document.MIMEType, &result.Report.Document.Name, nil
}
