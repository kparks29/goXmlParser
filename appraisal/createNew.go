package appraisal

// import (
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"github.com/homdna/homdna-models"
// 	// "github.com/homdna/homdna-service/requests"
// 	"io/ioutil"
// 	"net/http"
// )

// func PostNewHomdna(homdnaPayload *models.HomdnaModel, serviceApiKey *string) (*HomdnaResponse, error) {
// 	client := &http.Client{}
// 	url := "https://dev.homdna.com/homdnas"
// 	payload, _ := json.Marshal(*homdnaPayload)

// 	request, err := http.NewRequest("POST", url, bytes.NewReader(payload))
// 	fmt.Printf("request: %#v", request)
// 	request.Header["Content-Type"] = []string{"application/json"}
// 	request.Header["X-Service-Api-Key"] = []string{*serviceApiKey}
// 	response, err := client.Do(request)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if response.StatusCode != http.StatusOK {
// 		fmt.Printf("Status code: %v Status: %v RESPONSE: %v", response.StatusCode, response.Status)
// 		return nil, errors.New("Bad Status Code")
// 	}

// 	bodyBytes, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	homdna := &HomdnaResponse{}
// 	err = json.Unmarshal(bodyBytes, homdna)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return homdna, nil
// }
