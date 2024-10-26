package model

import (
	"fmt"
	"log"
	"strings"

	errs "github.com/diagmatrix/fblthp-archive/errors"

	"github.com/diagmatrix/fblthp-archive/utils"
)

type Card struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Layout          string   `json:"layout"`
	Finish          string   `json:"finish"`
	SubCards        []Card   `json:"subcards"`
	CMC             float64  `json:"cmc"`
	Colors          []string `json:"colors"`
	ColorIdentity   []string `json:"color_identity"`
	Types           []string `json:"types"`
	Subtypes        []string `json:"subtypes"`
	SetID           string   `json:"set_id"`
	OracleText      string   `json:"oracle_text"`
	Keywords        []string `json:"keywords"`
	CollectorNumber string   `json:"collector_number"`
	ArtistIDs       []string `json:"artist_id"`
}

func (c Card) String() string { /// For debugging purposes
	return fmt.Sprintf(
		"<Card: ID=%d, Name=%s, SetID=%s, CollectorNumber=%s, Finish=%s>",
		c.ID,
		c.Name,
		c.SetID,
		c.CollectorNumber,
		c.Finish,
	)
}

func (c Card) ToJSON(filename string) error {
	return utils.WriteJSON(c, filename)
}

func NewCardsFromJSON(filename string) ([]Card, error) {
	rawCards, err := NewRawCardsFromJSON(filename)
	if err != nil {
		return nil, err
	}

	log.Println("Populating cards from JSON...")
	cardList := []Card{}
	for _, rawCard := range rawCards {
		cards, err := rawCard.ToCards()
		if err != nil {
			return nil, err
		}
		cardList = append(cardList, *cards...)
	}
	return cardList, nil
}

type RawCard struct {
	// Core card fields
	ArenaID           int    `json:"arena_id"`
	ID                string `json:"id"`
	MtgoID            int    `json:"mtgo_id"`
	MtgoFoilID        int    `json:"mtgo_foil_id"`
	MultiverseIDs     []int  `json:"multiverse_id"`
	TgplayerID        int    `json:"tcgplayer_id"`
	TgcplayerEtchedID int    `json:"tcgplayer_etched_id"`
	CardmarketID      int    `json:"cardmarket_id"`
	Object            string `json:"object"`
	Layout            string `json:"layout"`
	OracleID          string `json:"oracle_id"`
	PrintsSearchURI   string `json:"prints_search_uri"`
	RulingsURI        string `json:"rulings_uri"`
	ScryfallURI       string `json:"scryfall_uri"`
	URI               string `json:"uri"`
	// Gameplay fields
	AllParts       []RelatedCard     `json:"all_parts"`
	CardFaces      []CardFace        `json:"card_faces"`
	CMC            float64           `json:"cmc"`
	ColorIdentity  []string          `json:"color_identity"`
	ColorIndicator []string          `json:"color_indicator"`
	Colors         []string          `json:"colors"`
	EdhrecRank     int               `json:"edhrec_rank"`
	HandModifier   string            `json:"hand_modifier"`
	Keywords       []string          `json:"keywords"`
	Legalities     map[string]string `json:"legalities"`
	LifeModifier   string            `json:"life_modifier"`
	Loyalty        string            `json:"loyalty"`
	ManaCost       string            `json:"mana_cost"`
	Name           string            `json:"name"`
	OracleText     string            `json:"oracle_text"`
	PennyRank      int               `json:"penny_rank"`
	Power          string            `json:"power"`
	ProducedMana   []string          `json:"produced_mana"`
	Reserved       bool              `json:"reserved"`
	Toughness      string            `json:"toughness"`
	TypeLine       string            `json:"type_line"`
	// Print fields
	Artist           string            `json:"artist"`
	ArtistIDs        []string          `json:"artist_ids"`
	AttractionLights []any             `json:"attribution_lights"` // TODO: What is this?
	Booster          bool              `json:"booster"`
	BorderColor      string            `json:"border_color"`
	CardBackID       string            `json:"card_back_id"`
	CollectorNumber  string            `json:"collector_number"`
	ContentWarning   bool              `json:"content_warning"`
	Digital          bool              `json:"digital"`
	Finishes         []string          `json:"finishes"`
	FlavorName       string            `json:"flavor_name"`
	FlavorText       string            `json:"flavor_text"`
	FrameEffects     []string          `json:"frame_effects"`
	Frame            string            `json:"frame"`
	FullArt          bool              `json:"full_art"`
	Games            []string          `json:"games"`
	HighresImage     bool              `json:"highres_image"`
	IllustrationID   string            `json:"illustration_id"`
	ImageStatus      string            `json:"image_status"`
	ImageURIs        map[string]string `json:"image_uris"`
	Oversized        bool              `json:"oversized"`
	Prices           map[string]string `json:"prices"`
	PrintedName      string            `json:"printed_name"`
	PrintedText      string            `json:"printed_text"`
	PrintedTypeLine  string            `json:"printed_type_line"`
	Promo            bool              `json:"promo"`
	PromoTypes       []string          `json:"promo_types"`
	PurchaseURIs     map[string]string `json:"purchase_uris"`
	Rarity           string            `json:"rarity"`
	RelatedURIs      map[string]string `json:"related_uris"`
	ReleasedAt       string            `json:"released_at"`
	Reprint          bool              `json:"reprint"`
	ScryfallSetURI   string            `json:"scryfall_set_uri"`
	SetName          string            `json:"set_name"`
	SetSearchURI     string            `json:"set_search_uri"`
	SetType          string            `json:"set_type"`
	SetURI           string            `json:"set_uri"`
	Set              string            `json:"set"`
	SetID            string            `json:"set_id"`
	StorySpotlight   bool              `json:"story_spotlight"`
	Textless         bool              `json:"textless"`
	Variation        bool              `json:"variation"`
	VariationOf      string            `json:"variation_of"`
	Watermark        string            `json:"watermark"`
}

func NewRawCardsFromJSON(filename string) ([]RawCard, error) {
	log.Println("Populating raw cards from JSON...")
	cards := []RawCard{}
	if err := utils.ReadJSON(filename, &cards); err != nil {
		return nil, err
	}

	return cards, nil
}

func (c RawCard) ToCards() (*[]Card, error) {
	cards := []Card{}
	for _, finish := range c.Finishes {
		card := &Card{}
		if c.Name == "" {
			return nil, errs.NewNoCardNameError()
		}
		if c.SetID == "" {
			return nil, errs.NewNoSetIDError()
		}
		if c.CollectorNumber == "" {
			return nil, errs.NewNoCollectorNumberError()
		}

		types, subtypes, err := parseTypeLine(c.TypeLine)
		if err != nil {
			return nil, err
		}

		card.Finish = finish
		card.Name = c.Name
		card.CMC = c.CMC
		card.Colors = c.Colors
		card.ColorIdentity = c.ColorIdentity
		card.Types = types
		card.Subtypes = subtypes
		card.SetID = c.Set
		card.OracleText = c.OracleText
		card.Keywords = c.Keywords
		card.CollectorNumber = c.CollectorNumber
		card.ArtistIDs = c.ArtistIDs

		for _, face := range c.CardFaces {
			subcard, err := face.ToCard()
			if err != nil {
				return nil, err
			}
			subcard.Finish = finish
			card.SubCards = append(card.SubCards, *subcard)
		}

		cards = append(cards, *card)
	}

	return &cards, nil
}

type RelatedCard struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Component string `json:"component"`
	Name      string `json:"name"`
	TypeLine  string `json:"type_line"`
	URI       string `json:"uri"`
}

type CardFace struct {
	Artist          string            `json:"artist"`
	ArtistID        string            `json:"artist_id"`
	CMC             float64           `json:"cmc"`
	ColorIndicator  []string          `json:"color_indicator"`
	Colors          []string          `json:"colors"`
	Defense         string            `json:"defense"`
	FlavorText      string            `json:"flavor_text"`
	IllustrationID  string            `json:"illustration_id"`
	ImageURIs       map[string]string `json:"image_uris"`
	Layout          string            `json:"layout"`
	Loyalty         string            `json:"loyalty"`
	ManaCost        string            `json:"mana_cost"`
	Name            string            `json:"name"`
	Object          string            `json:"object"`
	OracleID        string            `json:"oracle_id"`
	OracleText      string            `json:"oracle_text"`
	Power           string            `json:"power"`
	PrintedName     string            `json:"printed_name"`
	PrintedText     string            `json:"printed_text"`
	PrintedTypeLine string            `json:"printed_type_line"`
	Toughness       string            `json:"toughness"`
	TypeLine        string            `json:"type_line"`
	Watermark       string            `json:"watermark"`
}

func (c CardFace) ToCard() (*Card, error) {
	card := &Card{}
	if c.Name == "" {
		return nil, errs.NewNoCardNameError()
	}
	card.Name = c.Name
	card.CMC = c.CMC
	card.Colors = c.Colors
	card.ColorIdentity = c.ColorIndicator
	types, subtypes, err := parseTypeLine(c.TypeLine)
	if err != nil {
		return nil, err
	}
	card.Types = types
	card.Subtypes = subtypes
	card.OracleText = c.OracleText
	card.CollectorNumber = c.PrintedName
	card.ArtistIDs = []string{c.ArtistID}
	return card, nil
}

func parseTypeLine(typeLine string) ([]string, []string, error) {
	typeLines := strings.Split(typeLine, " // ")
	if len(typeLines) > 1 {
		firstTypes, firstSubtypes, err := parseTypeLine(typeLines[0])
		if err != nil {
			return nil, nil, errs.NewCannotParseTypelineError(typeLine, err.Error())
		}
		secondTypes, secondSubtypes, err := parseTypeLine(typeLines[1])
		if err != nil {
			return nil, nil, errs.NewCannotParseTypelineError(typeLine, err.Error())
		}
		types := append(firstTypes, secondTypes...)
		subtypes := append(firstSubtypes, secondSubtypes...)
		return types, subtypes, nil
	}

	typeAndSubtype := strings.Split(typeLine, " — ")
	if len(typeAndSubtype) > 2 {
		return nil, nil, errs.NewCannotParseTypelineError(typeLine, "Too many dashes")
	}

	types := strings.Split(typeAndSubtype[0], " ")

	var subtypes []string
	if len(typeAndSubtype) == 2 {
		subtypes = strings.Split(typeAndSubtype[1], " ")
	}
	return types, subtypes, nil
}
