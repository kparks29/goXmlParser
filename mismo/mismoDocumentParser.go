package mismo

import (
	"encoding/base64"
	"errors"
	"github.com/homdna/homdna-models"
	"github.com/kparks29/Document_Parser/documentParser"
)

type MismoDocumentParser struct {
}

func (this *MismoDocumentParser) Parse(body []byte, mimeType string, homdna *models.HomdnaModel) (*documentParser.ParsedDocument, error) {
	if mimeType != "application/xml" && mimeType != "text/xml" {
		return nil, errors.New("Only supports xml")
	}
	result, err := ParseXml(&body)
	if err != nil {
		return nil, err
	}

	address := CreateAddressModel(result)
	appliances, err := CreateApplianceModels(result.Property.Structure.KitchenAppliances)
	if err != nil {
		return nil, err
	}
	lotFeatures, err := CreateLotFeatureModels(result.Property.Structure.LotFeatures)
	if err != nil {
		return nil, err
	}
	structureFeatures := CreateStructureFeatureModels(result.Property.Structure.StructureFeatures)
	lotSize := GetSize(result.Property.LotInfo.LotSize)
	lot := CreateLot(&lotSize, lotFeatures)
	structures := CreateStructures(result.Property.Structure, structureFeatures, appliances)
	parsedHomdna := CreateHomdnaModel(address, lot, structures)
	mergedHomdna := MergeHomdnas(*homdna, *parsedHomdna)

	document, err := base64.StdEncoding.DecodeString(result.Report.Document.File)
	if err != nil {
		return nil, err
	}

	parsedFile := []documentParser.ParsedFile{}

	parsedFile = append(parsedFile, documentParser.ParsedFile{
		Body:     document,
		MimeType: result.Report.Document.MIMEType,
		Name:     result.Report.Document.Name,
	})

	parsedDocument := &documentParser.ParsedDocument{
		Documents: parsedFile,
		Homdna:    &mergedHomdna,
	}

	return parsedDocument, nil
}

func (this *MismoDocumentParser) SupportStandard(standard string) bool {
	if standard == "mismo" {
		return true
	}
	return false
}
