package documentParser

import (
	"github.com/homdna/homdna-models"
)

type DocumentParser interface {
	Parse(body []byte, mimeType string, homdna *models.HomdnaModel) (ParsedDocument, error)
	SupportStandard(standard string) bool
}
