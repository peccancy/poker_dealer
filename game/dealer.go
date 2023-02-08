package game

import (
	"github.com/google/uuid"
	"github.com/peccancy/poker_dealer/model"
)

func New(p model.Player, lock bool, pass string, bet int) uuid.UUID {
	table := model.New(p, lock, pass, bet)

	return table.ID
}

func Start() {

}
