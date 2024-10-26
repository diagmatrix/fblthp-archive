package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

// -----------------------------------------------------------------------------
// Generic files functions

func GetFile(url string, targetFile string) error {
	log.Println("Getting file from", url, "...")
	if fileExists(targetFile) {
		log.Println("File already exists:", targetFile)
		return nil
	}

	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	log.Println("Writing file to", targetFile, "...")
	file, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	return err
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

// -----------------------------------------------------------------------------
// JSON functions

func ReadJSON(filename string, target any) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(target)
}

func WriteJSON(item any, filename string) error {
	log.Println("Writing JSON to", filename, "...")
	if fileExists(filename) {
		log.Println("File already exists:", filename)
		return nil
	}

	cardJSON, err := json.Marshal(item)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(cardJSON)
	if err != nil {
		return err
	}
	return nil
}

func GetJSON(url string, target any) error {
	log.Println("Getting JSON from", url, "...")
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
