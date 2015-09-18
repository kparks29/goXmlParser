package main

import (
	"fmt"
	// "github.com/kparks29/Document_Parser/documentParser"
	// "github.com/kparks29/Document_Parser/mismo"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/homdna/homdna-models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	// 1) from command line get file_location, mimetype, standard, homdna_uuid
	if len(os.Args) < 5 {
		log.Fatalln("\n Missing Arguments! Need the following:  File Location, MIMEtype, Standard, and Homdna Uuid.")
	}
	filePath, mimeType, standard, homdna_uuid := os.Args[1], os.Args[2], os.Args[3], os.Args[4]

	serviceApiKey, err := GetApiKey()
	if err != nil {
		log.Fatalln(err)
	}
	// 2) get homdna version
	homdna, err := GetLatestHomdna(homdna_uuid, *serviceApiKey)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("\n homdna: %#v \n\n", *homdna)

	// 3) parse
	// 4) iterate through parsed files
	// 5) pos document & file
	// 6) add file to homdna model
	// 7) post new version
}

type Config struct {
	ServiceApiKey string `json:"service_api_key"`
}

func (this *Config) Validate() error {
	if this.ServiceApiKey == "" {
		return errors.New("A valid api key must be set. The hash key is used by other microservices to make requests to special APIs not intended for public consumption.")
	}
	return nil
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configData := &Config{}
	if err := decoder.Decode(configData); err != nil {
		return nil, err
	}
	return configData, nil
}

func GetApiKey() (*string, error) {
	config, err := LoadConfig("./document_parser.conf")
	if err != nil {
		return nil, errors.New("Error loading configuration file.")
	}
	if err = config.Validate(); err != nil {
		return nil, err
	}

	serviceApiKey := &config.ServiceApiKey
	return serviceApiKey, nil
}

func GetLatestHomdna(homdnaId string, serviceApiKey string) (*models.HomdnaModel, error) {
	url := "https://dev.homdna.com/homdnas/" + homdnaId + "/versions/latest"
	payload := []byte{}
	fmt.Printf("\n %v \n %v \n %v \n %v \n\n", "GET", url, payload, serviceApiKey)
	response, err := makeRequest("GET", url, &payload, &serviceApiKey)
	if err != nil {
		return nil, err
	}
	homdnaResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	homdna := models.HomdnaModel{}
	err = json.Unmarshal(homdnaResponse, &homdna)
	if err != nil {
		return nil, err
	}
	return &homdna, nil
}

// ALREADY WRITTEN MOST OF THIS HERE

// import (
// 	"bytes"
// 	"code.google.com/p/go-uuid/uuid"
// 	"encoding/base64"
// 	"encoding/json"
// 	"errors"
// 	"github.com/homdna/homdna-models"
// 	"github.com/homdna/homdna-service/domain"
// 	"github.com/homdna/homdna-service/requests"
// 	"io/ioutil"
// 	"net/http"
// 	"strings"
// )

// func main() {

// xmlFile, err := os.Open(filePath)
// if err != nil {
// 	log.Fatalln(err)
// }

// defer xmlFile.Close()

// readFile, err := ioutil.ReadAll(xmlFile)
// if err != nil {
// 	log.Fatalln(err)
// }

// return &readFile, nil

// 	// PASS IN A FILE TO RECEIVE A HOMDNA BACK (WORKFLOW 1)
// 	appraisalResponse, err := appraisal.UpdateHomdnaModel(file)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	// START OF POSTING HOMDNA OBJECT REMOTELY (WORKFLOW 2)
// 	appraisalConfig, err := appraisal.LoadAppraisalConfig("./appraisal.conf")
// 	if err != nil {
// 		log.Fatalln("Error loading HOMDNA appraisal configuration file.", err)
// 	}
// 	if err = appraisalConfig.Validate(); err != nil {
// 		log.Fatalln(err)
// 	}

// 	serviceApiKey := &appraisalConfig.ServiceApiKey

// 	// CREATE USER ACCOUNT REQUEST OBJECT
// 	// fullName := appraisalResponse.ParsedXML.Property.Owner.Name
// 	// accountCreationRequest := CreateAccount(os.Args[2], fullName, serviceApiKey)
// 	// address := &appraisalResponse.Homdna.Address

// 	// CREATE HOMDNA REQUEST OBJECT
// 	// homdnaRequest := &requests.HomdnaRequest{
// 	// 	StreetAddress:    address.StreetAddress,
// 	// 	City:             address.City,
// 	// 	State:            address.State,
// 	// 	PostalCode:       address.PostalCode,
// 	// 	PrimaryHomeOwner: *accountCreationRequest,
// 	// }

// 	CREATE HOMDNA
// 	if err = PostNewHomdna(appraisalResponse.Homdna, homdnaRequest, serviceApiKey); err != nil {
// 		log.Fatalln(err)
// 	}

// 	// PREPARE DOCUMENT
// 	homdnaId := "2f6a2416-90a9-47bb-a52e-b77248da5f3d" //temp id
// 	documentPayload, err := createDocumentRequest(&appraisalResponse.ParsedXML.Report.Document)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	// POST DOCUMENT
// 	documentResponse, err := appraisal.PostDocument(&homdnaId, documentPayload, serviceApiKey)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	document := domain.HomdnaDocument{}
// 	err = json.Unmarshal(*documentResponse, &document)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	// PREPARE FILE
// 	body, err := base64.StdEncoding.DecodeString(appraisalResponse.ParsedXML.Report.Document.File)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	md5Hash, err := appraisal.GetFileMd5(&body)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	// POST FILE
// 	_, err = appraisal.PostFile(homdnaId, md5Hash, document.Uuid, &appraisalResponse.ParsedXML.Report.Document, serviceApiKey)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	// ADD DOCUMENT ID TO HOMDNA
// 	appraisalResponse.Homdna.Documents = append(appraisalResponse.Homdna.Documents, document.Uuid)

// 	// POST VERSION
// 	_, err = appraisal.PostFirstVersion(homdnaId, appraisalResponse.Homdna, serviceApiKey)
// 	if err != nil {
// 		fmt.Println("\n\n failed at post version \n\n")
// 		log.Fatalln(err)
// 	}
// }

// func createDocumentRequest(document *appraisal.Document) (*[]byte, error) {
// 	uuid := uuid.New()
// 	docType := "Appraisal"
// 	var fileIds []string
// 	payload, err := json.Marshal(requests.HomdnaDocumentCreationRequest{
// 		Uuid:         &uuid,
// 		DocumentType: docType,
// 		DocumentName: document.Name,
// 		FileUuids:    fileIds,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &payload, err
// }

// func parseName(name string) (first string, last string) {
// 	first = ""
// 	last = strings.TrimSpace(name)

// 	idx := strings.LastIndex(last, " ")
// 	if idx > 0 {
// 		first = last[:idx]
// 		last = last[idx+1:]
// 	}
// 	return
// }

func makeRequest(method string, url string, payload *[]byte, serviceApiKey *string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, bytes.NewReader(*payload))
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
	return response, nil
}

// func CreateAccount(email string, fullName string, serviceApiKey *string) *requests.AccountCreationRequest {
// 	firstName, lastName := parseName(fullName)
// 	accountCreationRequest := &requests.AccountCreationRequest{
// 		Role:         "Home Owner",
// 		FullName:     &fullName,
// 		FirstName:    &firstName,
// 		LastName:     lastName,
// 		EmailAddress: email,
// 	}
// 	return accountCreationRequest
// }

// func PostNewHomdna(homdnaPayload *models.HomdnaModel, homdnaRequest *requests.HomdnaRequest, serviceApiKey *string) error {

// 	payload, err := json.Marshal(*homdnaRequest)
// 	if err != nil {
// 		return err
// 	}
// 	response, err := makeRequest("POST", "https://dev.homdna.com/homdnas", &payload, serviceApiKey)
// 	if err != nil {
// 		return err
// 	}
// 	bodyBytes, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Printf("\n response %v\n\n", string(bodyBytes))

// 	return nil
// }
