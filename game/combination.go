package game

const(
	HighCard Combination = iota
	Pair
	TwoPairs
	Three
	Straight
	Flush
	FullHouse
	Four
	StraightFlush
	RoyalFlush
)

type Combination int

var CombinationLevel = make([]Combination, 10, 10)

func init() {
	CombinationLevel = []Combination{
		HighCard,
		Pair,
		TwoPairs,
		Three,
		Straight,
		Flush,
		FullHouse,
		Four,
		StraightFlush,
		RoyalFlush,
	}
}

