package mismo

import (
	"encoding/json"
	"errors"
	"os"
)

type AppraisalConfig struct {
	ServiceApiKey string `json:"service_api_key"`
}

func (this *AppraisalConfig) Validate() error {
	if this.ServiceApiKey == "" {
		return errors.New("A valid api key must be set. The hash key is used by other microservices to make requests to special APIs not intended for public consumption.")
	}
	return nil
}

func LoadAppraisalConfig(filePath string) (*AppraisalConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configData := &AppraisalConfig{}
	if err := decoder.Decode(configData); err != nil {
		return nil, err
	}
	return configData, nil
}
