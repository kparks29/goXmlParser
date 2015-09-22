package documentParser

import (
	"encoding/json"
	"errors"
	"github.com/homdna/homdna-models"
	"os"
)

type DocumentParser interface {
	Parse(body []byte, mimeType string, homdna *models.HomdnaModel) (*ParsedDocument, error)
	SupportStandard(standard string) bool
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
