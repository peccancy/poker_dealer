package model

import "github.com/google/uuid"

type Table struct {
	ID      uuid.UUID
	Deck    CardDeck
	Players []Player
	Dealer  Player
	Owner   Player
	Locked  bool
	Pass    string
	Bet     int
}

func New(p Player, lock bool, pass string, bet int) Table {
	return Table{
		ID:      uuid.New(),
		Deck:    InitCardDeck(),
		Players: []Player{p},
		Dealer:  p,
		Owner:   p,
		Locked:  lock,
		Pass:    pass,
		Bet:     bet,
	}
}

func (t *Table) AddPlayer(p Player) {
	t.Players = append(t.Players, p)
}
