package model

import (
	"github.com/google/uuid"
)

type Player struct {
	ID    uuid.UUID
	hold  bool
	money int
	name  string
}

func (p Player) On() {
	p.hold = true
}

func (p Player) Off() {
	p.hold = false
}
