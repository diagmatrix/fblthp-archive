package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Card struct {
	layout          string
	SubCards        []Card
	Name            string
	CMC             float64
	Colors          []string
	ColorIdentity   []string
	Types           []string
	Subtypes        []string
	SetID           string
	OracleText      string
	keywords        []string
	CollectorNumber string
	ArtistID        string
}

func (c Card) String() string {
	return fmt.Sprintf(
		"<Card: Name=%s, CMC=%f, Colors=%v, ColorIdentity=%v, Types=%v, Subtypes=%v, SetID=%s, OracleText=%s, keywords=%v, CollectorNumber=%s, ArtistID=%s>",
		c.Name,
		c.CMC,
		c.Colors,
		c.ColorIdentity,
		c.Types,
		c.Subtypes,
		c.SetID,
		c.OracleText,
		c.keywords,
		c.CollectorNumber,
		c.ArtistID,
	)
}

func NewCardsFromJSON(filename string) ([]Card, error) {
	rawCards, err := NewRawCardsFromJSON(filename)
	if err != nil {
		return nil, err
	}
	cards := []Card{}
	for _, rawCard := range rawCards {
		card, err := rawCard.ToCard()
		if err != nil {
			return nil, err
		}
		cards = append(cards, *card)
	}
	return cards, nil
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
	cards := []RawCard{}
	if err := readJSON(filename, &cards); err != nil {
		return nil, err
	}

	return cards, nil
}

func (c RawCard) ToCard() (*Card, error) {
	card := &Card{}
	if c.Name == "" {
		return nil, errors.New("Card name is required")
	}
	if c.SetID == "" {
		return nil, errors.New("Set ID is required")
	}
	if c.CollectorNumber == "" {
		return nil, errors.New("Collector number is required")
	}
	card.Name = c.Name
	card.CMC = c.CMC
	card.Colors = c.Colors
	card.ColorIdentity = c.ColorIdentity
	types, subtypes, err := parseTypeLine(c.TypeLine)
	if err != nil {
		return nil, err
	}
	card.Types = types
	card.Subtypes = subtypes
	card.SetID = c.Set
	card.OracleText = c.OracleText
	card.keywords = c.Keywords
	card.CollectorNumber = c.CollectorNumber
	card.ArtistID = c.ArtistIDs[0] // TODO: WTF is this

	for _, face := range c.CardFaces {
		subcard, err := face.ToCard()
		if err != nil {
			return nil, err
		}
		card.SubCards = append(card.SubCards, *subcard)
	}

	return card, nil
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
		return nil, errors.New("Card name is required")
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
	card.ArtistID = c.ArtistID
	return card, nil
}

type RelatedCard struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Component string `json:"component"`
	Name      string `json:"name"`
	TypeLine  string `json:"type_line"`
	URI       string `json:"uri"`
}

func readJSON(filename string, target *[]RawCard) error { // TODO: Change target to any
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}

func parseTypeLine(typeLine string) ([]string, []string, error) {
	typeAndSubtype := strings.Split(typeLine, " — ")
	if len(typeAndSubtype) != 2 {
		return nil, nil, errors.New("Invalid type line")
	}

	return strings.Split(typeAndSubtype[0], " "), strings.Split(typeAndSubtype[1], " "), nil
}
