package main

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/homdna/homdna-models"
	"github.com/homdna/homdna-service/domain"
	"github.com/homdna/homdna-service/requests"
	"github.com/kparks29/Document_Parser/documentParser"
	"github.com/kparks29/Document_Parser/mismo"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type LatestHomdna struct {
	Homdna  models.HomdnaModel
	version int
}

func main() {
	// 1) from command line get file_location, mimetype, standard, homdna_uuid
	fmt.Println("Parsing arguments from terminal")
	if len(os.Args) < 5 {
		log.Fatalln("\n Missing Arguments! Need the following:  File Location, MIMEtype, Standard, and Homdna Uuid.")
	}
	filePath, mimeType, standard, homdna_uuid := os.Args[1], os.Args[2], os.Args[3], os.Args[4]
	parser := mismo.MismoDocumentParser{}

	if !parser.SupportStandard(standard) {
		log.Fatalln("Only Supports mismo standard")
	}

	serviceApiKey, err := documentParser.GetApiKey()
	if err != nil {
		log.Fatalln(err)
	}

	// 2) get homdna version
	fmt.Println("Get latest Homdna version from server")
	latestHomdna, err := GetLatestHomdna(homdna_uuid, *serviceApiKey)
	if err != nil {
		log.Fatalln(err)
	}

	// 3) parse & merge homdna
	fmt.Println("Parsing the file")
	file, err := mismo.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	parsedResult, err := parser.Parse(*file, mimeType, &latestHomdna.Homdna)
	if err != nil {
		log.Fatalln(err)
	}

	// 5) post document & file
	fmt.Println("Posting Document & adding id's to Homdna.Documents")
	for _, doc := range parsedResult.Documents {
		documentId, err := getDocumentId(&doc, homdna_uuid, serviceApiKey)
		if err != nil {
			log.Fatalln(err)
		}
		// 6) add document id to homdna model
		parsedResult.Homdna.Documents = append(parsedResult.Homdna.Documents, *documentId)
	}

	// 7) post new version
	fmt.Println("Posting new homdna version")
	_, err = PostNewHomdna(homdna_uuid, parsedResult.Homdna, latestHomdna.version, serviceApiKey)
	if err != nil {
		fmt.Println("\n\n failed at post version \n\n")
		log.Fatalln(err)
	}
}

func GetLatestHomdna(homdnaId string, serviceApiKey string) (*LatestHomdna, error) {
	url := "https://dev.homdna.com/homdnas/" + homdnaId + "/versions/latest"
	payload := []byte{}

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

	version, err := strconv.Atoi(response.Header["X-Homdna-Version"][0])
	if err != nil {
		return nil, err
	}

	latestHomdna := LatestHomdna{
		Homdna:  homdna,
		version: version,
	}

	return &latestHomdna, nil
}

func getDocumentId(doc *documentParser.ParsedFile, homdna_uuid string, serviceApiKey *string) (*string, error) {
	// Prepare Docoument
	documentPayload, err := createDocumentRequest(doc)
	if err != nil {
		return nil, err
	}
	// POST DOCUMENT
	documentResponse, err := PostDocument(&homdna_uuid, documentPayload, serviceApiKey)
	if err != nil {
		return nil, err
	}

	document := domain.HomdnaDocument{}
	err = json.Unmarshal(*documentResponse, &document)
	if err != nil {
		return nil, err
	}

	// PREPARE FILE
	body := base64.StdEncoding.EncodeToString(doc.Body)
	md5Hash, err := GetFileMd5(&doc.Body)
	if err != nil {
		return nil, err
	}
	// POST FILE
	_, err = PostFile(homdna_uuid, document.Uuid, md5Hash, body, doc, serviceApiKey)
	if err != nil {
		return nil, err
	}

	return &document.Uuid, nil
}

func createDocumentRequest(document *documentParser.ParsedFile) (*[]byte, error) {
	uuid := uuid.New()
	docType := "Appraisal"
	var fileIds []string
	payload, err := json.Marshal(requests.HomdnaDocumentCreationRequest{
		Uuid:         &uuid,
		DocumentType: docType,
		DocumentName: document.Name,
		FileUuids:    fileIds,
	})
	if err != nil {
		return nil, err
	}
	return &payload, err
}

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

func GetFileMd5(file *[]byte) (*string, error) {
	h := md5.New()
	_, err := h.Write(*file)
	if err != nil {
		return nil, err
	}

	// must convert to string becuase the EncodeToString expects a byte array represented by ascii characters and the md5.sum returns the numerical value
	md5String := fmt.Sprintf("%x", h.Sum(nil))
	generatedChecksum := base64.StdEncoding.EncodeToString([]byte(md5String))

	return &generatedChecksum, nil
}

func PostDocument(homdnaId *string, documentPayload *[]byte, serviceApiKey *string) (*[]byte, error) {
	url := "https://dev.homdna.com/homdnas/" + *homdnaId + "/documents"

	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewReader(*documentPayload))
	if err != nil {
		return nil, err
	}
	request.Header["X-Service-Api-Key"] = []string{*serviceApiKey}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		fmt.Printf("Status code: %v Status: %v", response.StatusCode, response.Status)
		return nil, errors.New("Bad Request")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

func PostFile(homdnaId string, documentId string, md5Hash *string, encodedDoc string, document *documentParser.ParsedFile, serviceApiKey *string) (*[]byte, error) {
	client := &http.Client{}
	url := "https://dev.homdna.com/homdnas/" + homdnaId + "/documents/" + documentId + "/files"

	request, err := http.NewRequest("POST", url, bytes.NewReader([]byte(encodedDoc)))
	request.Header["Content-Type"] = []string{document.MimeType}
	request.Header["Content-MD5"] = []string{*md5Hash}
	request.Header["X-Service-Api-Key"] = []string{*serviceApiKey}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		fmt.Printf("Status code: %v Status: %v", response.StatusCode, response.Status)
		return nil, errors.New("Bad Request")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

func PostNewHomdna(homdnaId string, homdna *models.HomdnaModel, version int, serviceApiKey *string) (*[]byte, error) {
	client := &http.Client{}
	url := "https://dev.homdna.com/homdnas/" + homdnaId + "/versions"
	jsonPayload, err := json.Marshal(*homdna)
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(jsonPayload))
	// payload := base64.StdEncoding.EncodeToString(jsonPayload)
	md5Hash, err := GetFileMd5(&jsonPayload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader([]byte(jsonPayload)))
	if err != nil {
		return nil, err
	}
	request.Header["Content-MD5"] = []string{*md5Hash}
	request.Header["X-Service-Api-Key"] = []string{*serviceApiKey}
	request.Header["X-Homdna-Modified-Version"] = []string{strconv.Itoa(version)}
	request.Header["X-Homdna-Version"] = []string{strconv.Itoa(version + 1)}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("\n ERROR MESSAGE: %v\n\n", string(body))
		fmt.Printf("Status code: %v Status: %v", response.StatusCode, response.Status)
		return nil, errors.New("Bad Request")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}
