package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	cardsFile := "cards/cards" + time.Now().Format("2006-01-02") + ".json"
	if err := GetScryfallCards(cardsFile); err != nil {
		log.Fatalf("failed to get Scryfall cards: %v", err)
	}
}

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

func (s ScryfallBulkData) String() string {
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

func GetScryfallCards(filename string) error {
	log.Printf("Getting Scryfall cards...")
	var bulkData ScryfallBulkData
	err := GetJSON(SCRYFALL_URL, &bulkData)
	if err != nil {
		return err
	}

	return GetFile(bulkData.DownloadURI, filename)
}

const SCRYFALL_URL = "https://api.scryfall.com/bulk-data/default-cards"

func GetJSON(url string, target any) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func GetFile(url string, targetFile string) error {
	log.Println("Writing file to", targetFile)
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	file, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	return err
}
