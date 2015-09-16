package appraisal

import (
	"encoding/xml"
	"github.com/homdna/homdna-models"
)

type FilePayload struct {
	file_payload string
}

type DocumentPayload struct {
	Uuid         *string
	DocumentType string
	DocumentName string
	FileUuids    []string
}

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
