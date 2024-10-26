package utils

import (
	"fmt"
	"time"
)

func GetScryfallCards(filename string) error {
	bulkData := ScryfallBulkData{}
	err := GetJSON(SCRYFALL_BULK_URL, &bulkData)
	if err != nil {
		return err
	}

	return GetFile(bulkData.DownloadURI, filename)
}

const SCRYFALL_BULK_URL = "https://api.scryfall.com/bulk-data/default-cards"

type ScryfallBulkData struct {
	ID              string    `json:"id"`
	URI             string    `json:"uri"`
	Type            string    `json:"type"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	DownloadURI     string    `json:"download_uri"`
	UpdatedAt       time.Time `json:"updated_at"`
	Size            int       `json:"size"`
	ContentEncoding string    `json:"content_encoding"`
}

func (s ScryfallBulkData) String() string { // For debugging purposes
	return fmt.Sprintf(
		"<ScryfallBulkData: ID=%s, URI=%s, Type=%s, Name=%s, Description=%s, DownloadURI=%s, UpdatedAt=%s, Size=%d, ContentEncoding=%s>",
		s.ID,
		s.URI,
		s.Type,
		s.Name,
		s.Description,
		s.DownloadURI,
		s.UpdatedAt,
		s.Size,
		s.ContentEncoding,
	)
}
