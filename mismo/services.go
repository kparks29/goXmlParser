package mismo

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/homdna/homdna-models"
	"github.com/kparks29/Document_Parser/documentParser"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetSize(sizeString string) float64 {
	sizeString = strings.Split(sizeString, " ")[0]
	sizeString = strings.Join(strings.Split(sizeString, ","), "")
	size, _ := strconv.Atoi(sizeString)
	return float64(size)
}

func CreateRooms(bedrooms string, bathrooms string, kitchenAppliances []*models.ApplianceModel) []*models.RoomModel {
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

func CreateStructures(structureInfo Structure, structureFeatures []*models.FeatureModel, kitchenAppliances []*models.ApplianceModel) []*models.StructureModel {
	var structures []*models.StructureModel

	name := "House"
	size := GetSize(structureInfo.Size)
	levels, _ := strconv.Atoi(structureInfo.Levels)
	id := uuid.New()
	space := models.SpaceModel{
		Id:       &id,
		Features: structureFeatures,
	}
	rooms := CreateRooms(structureInfo.Bedrooms, structureInfo.Bathrooms, kitchenAppliances)

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

func CreateLot(size *float64, features []*models.FeatureModel) *models.LotModel {
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

func CreateStructureFeatureModels(structureFeatures []StructureFeature) []*models.FeatureModel {
	var features []*models.FeatureModel

	for i := 0; i < len(structureFeatures); i++ {
		feature := new(models.FeatureModel)
		feature.ComputeId()
		featureDescription := strings.Split(structureFeatures[i].Description, "/")

		feature.Name = structureFeatures[i].Name
		feature.Notes = &featureDescription[0]

		var condition *int
		if featureDescription[1] == "Gd" {
			rating := 2
			condition = &rating
		}

		feature.Condition = condition

		features = append(features, feature)
	}
	return features
}

func CreateLotFeatureModels(lotFeatures []LotFeature) ([]*models.FeatureModel, error) {
	var features []*models.FeatureModel
	elementFeatures, err := GetFeatures("lot_features")
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n\n elementFeatures %#v \n\n", *elementFeatures)
	for _, lotFeature := range lotFeatures {
		feature := new(models.FeatureModel)
		feature.ComputeId()
		feature.Name = lotFeature.Name
		feature.Notes = &lotFeature.Description

		features = append(features, feature)
	}
	return features, err
}

func CreateApplianceModels(kitchenAppliances []Appliance) ([]*models.ApplianceModel, error) {
	var appliances []*models.ApplianceModel

	applianceElements, err := GetAppliances()
	if err != nil {
		return nil, err
	}
	// loop through all the appliances and if they exist or the count is greater than 0 add the model\
	for i := 0; i < len(kitchenAppliances); i++ {
		if count, _ := strconv.Atoi(kitchenAppliances[i].Count); count > 0 || kitchenAppliances[i].ExistsIndicator == "Y" {
			appliance := new(models.ApplianceModel)
			appliance.ComputeId()
			appliance.Name = &kitchenAppliances[i].Type
			appliance.ApplianceId = GetApplianceId(*appliance, *applianceElements)

			appliances = append(appliances, appliance)
		}
	}
	return appliances, nil
}

func CreateAddressModel(data *Result) *models.AddressModel {
	return &models.AddressModel{
		StreetAddress: data.Property.StreetAddress,
		City:          data.Property.City,
		State:         data.Property.State,
		PostalCode:    data.Property.PostalCode,
	}
}

func CreateHomdnaModel(address *models.AddressModel, lot *models.LotModel, structures []*models.StructureModel) *models.HomdnaModel {
	model := &models.HomdnaModel{
		Address:    *address,
		Lot:        *lot,
		Structures: structures,
	}
	model.ComputeIds()
	return model
}

func ReadFile(filePath string) (*[]byte, error) {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer xmlFile.Close()

	readFile, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}

	return &readFile, nil
}

func ParseXml(file *[]byte) (*Result, error) {

	result := &Result{}

	err := xml.Unmarshal(*file, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetAppliances() (*map[string]string, error) {
	serviceApiKey, err := documentParser.GetApiKey()
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://appliance.homdna.com/appliances", bytes.NewReader([]byte{}))
	request.Header["Content-Type"] = []string{"application/json"}
	request.Header["X-Service-Api-Key"] = []string{*serviceApiKey}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("Status code: %v Status: %v RESPONSE: %v", response.StatusCode, response.Status, string(body))
		return nil, errors.New("Bad Status Code")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	parsedAppliances := []ApplianceElement{}
	err = json.Unmarshal(body, &parsedAppliances)
	if err != nil {
		return nil, err
	}

	appliances := make(map[string]string)
	for _, appliance := range parsedAppliances {
		if appliance.Brand.Name == "other" {
			appliances[appliance.Type.Name] = appliance.ApplianceId
		}
	}

	return &appliances, nil
}

func GetApplianceId(appliance models.ApplianceModel, appliances map[string]string) string {
	applianceType := ""
	switch {
	case *appliance.Name == "RangeOven":
		applianceType = "stove_top_oven"
		break
	case *appliance.Name == "Disposal":
		applianceType = "garbage_disposal"
		break
	case *appliance.Name == "Dishwasher":
		applianceType = "dishwasher"
		break
	}
	return appliances[applianceType]
}

func GetFeatures(filter string) (*[]FeatureElement, error) {
	serviceApiKey, err := documentParser.GetApiKey()
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://appliance.homdna.com/elements", bytes.NewReader([]byte{}))
	request.Header["Content-Type"] = []string{"application/json"}
	request.Header["X-Service-Api-Key"] = []string{*serviceApiKey}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("Status code: %v Status: %v RESPONSE: %v", response.StatusCode, response.Status, string(body))
		return nil, errors.New("Bad Status Code")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var elementFeatures = Elements{}

	err = json.Unmarshal(body, &elementFeatures)
	if err != nil {
		return nil, err
	}

	features := make(map[string]interface{})
	for _, feature := range elementFeatures.RoomFeatures {
		features[feature.feature_name] = feature
	}

	fmt.Printf("\n\n %#v \n\n", features)

	switch {
	case filter == "lot_features":
		return &elementFeatures.LotFeatures, nil
	case filter == "exterior_features":
		return &elementFeatures.StructureFeatures, nil
	case filter == "room_features":
		return &elementFeatures.RoomFeatures, nil
	}

	return nil, errors.New("Feature type does not exist")

}

func GetFeatureType(feature models.FeatureModel, features []FeatureElement) FeatureType {
	featureType := ""
	switch {
	case feature.Name == "RangeOven":
		featureType = "stove_top_oven"
		break
	case feature.Name == "Disposal":
		featureType = "garbage_disposal"
		break
	case feature.Name == "Dishwasher":
		featureType = "dishwasher"
		break
	}
	return features[featureType]
}
