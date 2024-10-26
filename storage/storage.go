package storage

import "github.com/diagmatrix/fblthp-archive/model"

// Type alias
type Card = model.Card

type Storage interface {
	// Card
	CreateCard(card Card) (int, error)
	GetCard(id int) (Card, error)
	ListCard() ([]Card, error)
	UpadateCard(id int, card Card) error
	DeleteCard(id int) error
}
