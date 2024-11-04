package storage

import (
	errs "github.com/diagmatrix/fblthp-archive/exceptions"

	"github.com/diagmatrix/fblthp-archive/utils"
)

type MemoryStorage struct {
	cards       []Card
	cardCounter utils.Counter
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		cards:       []Card{},
		cardCounter: *utils.NewCounter(),
	}
}

func (s *MemoryStorage) GetCard(id int) (Card, error) {
	if id < 1 || id > s.cardCounter.Current() {
		return Card{}, errs.NewCardNotFoundError(id)
	}
	for _, card := range s.cards {
		if card.ID == id {
			return card, nil
		}
	}
	return Card{}, errs.NewCardNotFoundError(id)
}
