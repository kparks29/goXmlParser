package main

import (
	"../../../code.google.com/p/go-uuid/uuid"
	"../../../github.com/homdna/homdna-service/models"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Document struct {
	Name         string `xml:"_Name,attr"`
	EncodingType string `xml:"_EncodingType,attr"`
	MIMEType     string `xml:"_MIMEType,attr"`
	Type         string `xml:"_Type,attr"`
	File         string `xml:"DOCUMENT"`
}

type MarketData struct {
	Type           string `xml:"_Type,attr"`
	TrendType      string `xml:"_TrendType,attr"`
	MonthRangeType string `xml:"_MonthRangeType,attr"`
	Amount         string `xml:"_Amount,attr"`
	Count          string `xml:"_Count,attr"`
	Rate           string `xml:"_Rate,attr"`
}

type Report struct {
	AppraisalFormType              string       `xml:"AppraisalFormType,attr"`
	TitleDescription               string       `xml:"_TitleDescription,attr"`
	AppraisalFormVersionIdentifier string       `xml:"AppraisalFormVersionIdentifier,attr"`
	AppraisalPurposeType           string       `xml:"AppraisalPurposeType,attr"`
	MarketData                     []MarketData `xml:"FORM>MARKET>MARKET_INVENTORY"`
	Document                       Document     `xml:"EMBEDDED_FILE"`
}

type Valuation struct {
	Amount string `xml:"PropertyAppraisedValueAmount,attr"`
	Date   string `xml:"AppraisalEffectiveDate,attr"`
}

type CarStorageLocation struct {
	Type            string `xml:"_Type,attr"`
	ExistsIndicator string `xml:"__ExistsIndicator,attr"`
}

type CarStorage struct {
	ParkingSpacesCount         string               `xml:"ParkingSpacesCount,attr"`
	ParkingSpaceIdentifier     string               `xml:"ParkingSpaceIdentifier,attr"`
	ParkingSpaceAssignmentType string               `xml:"ParkingSpaceAssignmentType,attr"`
	CarStorageLocations        []CarStorageLocation `xml:"CAR_STORAGE_LOCATION"`
}

type SiteUtility struct {
	Type   string `xml:"_Type,attr"`
	Public string `xml:"_PublicIndicator,attr"`
}

type SiteFeature struct {
	Type    string `xml:"_Type,attr"`
	Comment string `xml:"_Comment,attr"`
}

type LotInfo struct {
	ZoningClassification string        `xml:"_ZoningClassificationDescription,attr"`
	LotSize              string        `xml:"SquareFeetCount,attr"`
	SiteFeatures         []SiteFeature `xml:"SITE_FEATURE"`
	SiteUtilities        []SiteUtility `xml:"SITE_UTILITY"`
}

type Owner struct {
	Name string `xml:"_Name,attr"`
}

type OtherStructure struct {
	SequenceId string `xml:"PropertyFeatureSequenceIdentifier,attr"`
	Name       string `xml:"PropertyFeatureName,attr"`
}

type Appliance struct {
	Type            string `xml:"_Type,attr"`
	Count           string `xml:"_Count,attr"`
	ExistsIndicator string `xml:"_ExistsIndicator,attr"`
}

type Heating struct {
	Type        string `xml:"_UnitDescription,attr"`
	Description string `xml:"_FuelDescription,attr"`
}

type Cooling struct {
	Centralized string `xml:"_CentralizedIndicator,attr"`
}

type LotFeature struct {
	Name        string `xml:"_Type,attr"`
	Description string `xml:"_Description,attr"`
}

type StructureFeature struct {
	Name        string `xml:"_Type,attr"`
	Description string `xml:"_ConditionDescription,attr"`
}

type Amenity struct {
	Type                 string `xml:"_Type,attr"`
	Count                string `xml:"_Count,attr"`
	ExistsIndicator      string `xml:"_ExistsIndicator,attr"`
	DetailedDescription  string `xml:"_DetailedDescription,attr"`
	TypeOtherDescription string `xml:"_TypeOtherDescription,attr"`
}

type Structure struct {
	Levels            string             `xml:"StoriesCount,attr"`
	YearBuilt         string             `xml:"PropertyStructureBuiltYear,attr"`
	Bedrooms          string             `xml:"TotalBedroomCount,attr"`
	Bathrooms         string             `xml:"TotalBathroomCount,attr"`
	Size              string             `xml:"GrossLivingAreaSquareFeetCount,attr"`
	StructureFeatures []StructureFeature `xml:"INTERIOR_FEATURE"`
	LotFeatures       []LotFeature       `xml:"EXTERIOR_FEATURE"`
	Heating           Heating            `xml:"HEATING"`
	Cooling           Cooling            `xml:"COOLING"`
	KitchenAppliances []Appliance        `xml:"KITCHEN_EQUIPMENT"`
	Amenities         []Amenity          `xml:"AMENITY"`
}

type Property struct {
	StreetAddress string    `xml:"_StreetAddress,attr"`
	City          string    `xml:"_City,attr"`
	State         string    `xml:"_State,attr"`
	PostalCode    string    `xml:"_PostalCode,attr"`
	County        string    `xml:"_County,attr"`
	Structure     Structure `xml:"STRUCTURE"`
	Owner         Owner     `xml:"_OWNER"`
	LotInfo       LotInfo   `xml:"SITE"`
}

type Result struct {
	XMLName         xml.Name         `xml:"VALUATION_RESPONSE"`
	Property        Property         `xml:"PROPERTY"`
	OtherStructures []OtherStructure `xml:"VALUATION_METHODS>SALES_COMPARISON>OTHER_FEATURE"`
	Valuation       Valuation        `xml:"VALUATION"`
	MISMOVersionId  string           `xml:"MISMOVersionID,attr"`
	Report          Report           `xml:"REPORT"`
}

func getSize(sizeString string) float64 {
	sizeString = strings.Split(sizeString, " ")[0]
	sizeString = strings.Join(strings.Split(sizeString, ","), "")
	size, _ := strconv.Atoi(sizeString)
	return float64(size)
}

func createRooms(bedrooms string, bathrooms string, kitchenAppliances []*models.ApplianceModel) []*models.RoomModel {
	var rooms []*models.RoomModel

	bedroomCount, _ := strconv.Atoi(bedrooms)
	bathroomCount, _ := strconv.Atoi(bathrooms)

	for i := 0; i < bedroomCount; i++ {
		name := "Bedroom"
		id := uuid.New()
		room := models.RoomModel{
			Name:     &name,
			RoomType: "bedroom",
			Level:    1,
			SpaceModel: models.SpaceModel{
				Id: &id,
			},
		}
		rooms = append(rooms, &room)
	}
	for i := 0; i < bathroomCount; i++ {
		name := "Bathroom"
		id := uuid.New()
		room := models.RoomModel{
			Name:     &name,
			RoomType: "bathroom",
			Level:    1,
			SpaceModel: models.SpaceModel{
				Id: &id,
			},
		}
		rooms = append(rooms, &room)
	}

	kitchenName := "Kitchen"
	kitchenId := uuid.New()
	kitchen := models.RoomModel{
		Name:     &kitchenName,
		RoomType: "kitchen",
		Level:    1,
		SpaceModel: models.SpaceModel{
			Id:         &kitchenId,
			Appliances: kitchenAppliances,
		},
	}

	rooms = append(rooms, &kitchen)

	return rooms
}

func createStructures(structureInfo Structure, structureFeatures []*models.FeatureModel, kitchenAppliances []*models.ApplianceModel) []*models.StructureModel {
	var structures []*models.StructureModel

	name := "House"
	size := getSize(structureInfo.Size)
	levels, _ := strconv.Atoi(structureInfo.Levels)
	id := uuid.New()
	space := models.SpaceModel{
		Id:       &id,
		Features: structureFeatures,
	}
	rooms := createRooms(structureInfo.Bedrooms, structureInfo.Bathrooms, kitchenAppliances)

	structure := &models.StructureModel{
		SpaceModel:    space,
		Name:          &name,
		StructureType: "main_house",
		Size:          &size,
		Levels:        &levels,
		Rooms:         rooms,
	}

	return append(structures, structure)
}

func createLot(size *float64, features []*models.FeatureModel) *models.LotModel {
	lotSize := size
	id := uuid.New()
	space := models.SpaceModel{
		Features: features,
		Id:       &id,
	}
	return &models.LotModel{
		LotSize:    lotSize,
		SpaceModel: space,
	}
}

func createStructureFeatureModels(structureFeatures []StructureFeature) []*models.FeatureModel {
	var features []*models.FeatureModel

	for i := 0; i < len(structureFeatures); i++ {
		feature := new(models.FeatureModel)
		feature.ComputeId()
		featureDescription := strings.Split(structureFeatures[i].Description, "/")

		feature.Name = structureFeatures[i].Name
		feature.Notes = &featureDescription[0]

		var condition int
		if featureDescription[1] == "Gd" {
			condition = 2
		}

		feature.Condition = &condition

		features = append(features, feature)
	}
	return features
}

func createLotFeatureModels(lotFeatures []LotFeature) []*models.FeatureModel {
	var features []*models.FeatureModel
	for i := 0; i < len(lotFeatures); i++ {
		feature := new(models.FeatureModel)
		feature.ComputeId()
		feature.Name = lotFeatures[i].Name
		feature.Notes = &lotFeatures[i].Description

		features = append(features, feature)
	}
	return features
}

func createApplianceModels(kitchenAppliances []Appliance) []*models.ApplianceModel {
	var appliances []*models.ApplianceModel
	// loop through all the appliances and if they exist or the count is greater than 0 add the model\
	// TODO Match them up with an ApplianceId from the Appliance DB
	for i := 0; i < len(kitchenAppliances); i++ {
		if count, _ := strconv.Atoi(kitchenAppliances[i].Count); count > 0 || kitchenAppliances[i].ExistsIndicator == "Y" {
			appliance := new(models.ApplianceModel)
			appliance.ComputeId()
			appliance.Name = &kitchenAppliances[i].Type

			appliances = append(appliances, appliance)
		}
	}
	return appliances
}

func createAddressModel(data *Result) *models.AddressModel {
	return &models.AddressModel{
		StreetAddress: data.Property.StreetAddress,
		City:          data.Property.City,
		State:         data.Property.State,
		PostalCode:    data.Property.PostalCode,
	}
}

func createHomdnaModel(address *models.AddressModel, lot *models.LotModel, structures []*models.StructureModel) *models.HomdnaModel {
	model := &models.HomdnaModel{
		Address:    *address,
		Lot:        *lot,
		Structures: structures,
	}
	model.ComputeIds()
	return model
}

func getFile() *Result {

	var result = new(Result)

	xmlFile, loadErr := os.Open(os.Args[1])
	if loadErr != nil {
		fmt.Println("Error opening file:", loadErr)
		os.Exit(1)
	}
	defer xmlFile.Close()

	readFile, _ := ioutil.ReadAll(xmlFile)

	err := xml.Unmarshal(readFile, &result)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(2)
	}

	return result
}

func main() {

	xmlFile := getFile()

	address := createAddressModel(xmlFile)
	appliances := createApplianceModels(xmlFile.Property.Structure.KitchenAppliances)
	lotFeatures := createLotFeatureModels(xmlFile.Property.Structure.LotFeatures)
	structureFeatures := createStructureFeatureModels(xmlFile.Property.Structure.StructureFeatures)
	lotSize := getSize(xmlFile.Property.LotInfo.LotSize)
	lot := createLot(&lotSize, lotFeatures)
	structures := createStructures(xmlFile.Property.Structure, structureFeatures, appliances)
	homdna := createHomdnaModel(address, lot, structures)

}
