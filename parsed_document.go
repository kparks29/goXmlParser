package documentParser

import (
	"github.com/homdna/homdna-models"
)

type ParsedDocument struct {
	Documents []ParsedFile
	Homdna    *models.HomdnaModel
}

type ParsedFile struct {
	Body     []byte
	MimeType string
	Name     string
}
