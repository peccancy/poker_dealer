package game

import (
	"math/rand"
	deck "github.com/peccancy/poker_dealer/data_card"
)

type Table struct {
	Deck []deck.Card
}

func Shuffle(a []deck.Card, source rand.Source) {
	random := rand.New(source)
	for i := len(a) - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func Start()  {

}
