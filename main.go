package main

import (
	"log"
	"time"

	"github.com/diagmatrix/fblthp-archive/model"
	"github.com/diagmatrix/fblthp-archive/utils"
)

func main() {
	cardsFile := "raw/cards" + time.Now().Format("2006-01-02") + ".json"
	if err := utils.GetScryfallCards(cardsFile); err != nil {
		log.Fatalf("failed to get Scryfall cards: %v", err)
	}
	cards, err := model.NewCardsFromJSON(cardsFile)
	if err != nil {
		log.Fatalf("failed to get cards from JSON: %v", err)
	}

	for _, card := range cards {
		err = card.ToJSON("card/" + card.SetID + "_" + card.CollectorNumber + "_" + card.Finish + ".json")
		if err != nil {
			log.Fatalf("failed to write card to JSON: %v", err)
		}
	}
}
