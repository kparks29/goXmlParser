package appraisal

import (
	// "bytes"
	"code.google.com/p/go-uuid/uuid"
	"crypto/md5"
	"encoding/base64"
	// "encoding/json"
	"encoding/xml"
	// "fmt"
	"github.com/homdna/homdna-models"
	"io/ioutil"
	// "net/http"
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

func CreateLotFeatureModels(lotFeatures []LotFeature) []*models.FeatureModel {
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

func CreateApplianceModels(kitchenAppliances []Appliance) []*models.ApplianceModel {
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

func GetFileMd5(file *string) [16]byte {
	decoded, _ := base64.StdEncoding.DecodeString(*file)
	md5Hash := md5.Sum(decoded)
	return md5Hash
}

// func PostDocument(homdnaId string, document *Document, serviceApiKey *string) string {
// 	id := uuid.New()
// 	url := "https://dev.homdna.com/homdnas/" + homdnaId + "/documents"
// 	docType := "Appraisal"
// 	payload, _ := json.Marshal(DocumentPayload{
// 		Uuid:         &id,
// 		DocumentName: document.Name,
// 		DocumentType: docType,
// 	})

// 	client := &http.Client{}
// 	request, err := http.NewRequest("POST", url, bytes.NewReader(payload))
// 	request.Header["X-Service-Api-Key"] = []string{*serviceApiKey}
// 	response, err := client.Do(request)

// 	if err != nil {
// 		os.Exit(4)
// 	}
// 	if response.StatusCode != http.StatusOK {
// 		fmt.Sprintf("Status code: %v Status: %v", response.StatusCode, response.Status)
// 		os.Exit(5)
// 	}
// 	body, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		os.Exit(6)
// 	}

// 	fmt.Printf("\n\nbody\n%v\n\n", body)
// 	return "Success"
// }

// func PostFile(homdna *HomdnaResponse, md5Hash [16]byte, document *Document, serviceApiKey *string) string {
// 	documentId := PostDocument(homdna.id, document, serviceApiKey)
// 	client := &http.Client{}
// 	url := "https://dev.homdna.com/homdnas/" + homdna.id + "/documents/" + documentId + "/files"
// 	payload, _ := json.Marshal(FilePayload{
// 		file_payload: document.File,
// 	})

// 	request, err := http.NewRequest("POST", url, bytes.NewReader(payload))
// 	request.Header["Content-Type"] = []string{document.MIMEType}
// 	request.Header["Content-MD5"] = []string{string(md5Hash[:])}
// 	request.Header["X-Service-Api-Key"] = []string{*serviceApiKey}
// 	response, err := client.Do(request)

// 	if err != nil {
// 		os.Exit(4)
// 	}
// 	if response.StatusCode != http.StatusOK {
// 		fmt.Sprintf("Status code: %v Status: %v", response.StatusCode, response.Status)
// 		os.Exit(5)
// 	}
// 	body, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		os.Exit(6)
// 	}

// 	fmt.Printf("\n\nbody\n%v\n\n", body)
// 	return "Success"
// }
