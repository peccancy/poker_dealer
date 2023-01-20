package model

var (
	Club    = Suit("club")
	Diamond = Suit("diamond")
	Heart   = Suit("hearts")
	Spade   = Suit("spade")

	FV2  = FaceValue("2")
	FV3  = FaceValue("3")
	FV4  = FaceValue("4")
	FV5  = FaceValue("5")
	FV6  = FaceValue("6")
	FV7  = FaceValue("7")
	FV8  = FaceValue("8")
	FV9  = FaceValue("9")
	FV10 = FaceValue("10")
	FVJ  = FaceValue("J")
	FVQ  = FaceValue("Q")
	FVK  = FaceValue("K")
	FVA  = FaceValue("A")
)

type CardDeck []Card

func InitCardDeck() CardDeck {
	return CardDeck{
		{Club, FV2},
		{Club, FV3},
		{Club, FV4},
		{Club, FV5},
		{Club, FV6},
		{Club, FV7},
		{Club, FV8},
		{Club, FV9},
		{Club, FV10},
		{Club, FVJ},
		{Club, FVQ},
		{Club, FVK},
		{Club, FVA},

		{Diamond, FV2},
		{Diamond, FV3},
		{Diamond, FV4},
		{Diamond, FV5},
		{Diamond, FV6},
		{Diamond, FV7},
		{Diamond, FV8},
		{Diamond, FV9},
		{Diamond, FV10},
		{Diamond, FVJ},
		{Diamond, FVQ},
		{Diamond, FVK},
		{Diamond, FVA},

		{Heart, FV2},
		{Heart, FV3},
		{Heart, FV4},
		{Heart, FV5},
		{Heart, FV6},
		{Heart, FV7},
		{Heart, FV8},
		{Heart, FV9},
		{Heart, FV10},
		{Heart, FVJ},
		{Heart, FVQ},
		{Heart, FVK},
		{Heart, FVA},

		{Spade, FV2},
		{Spade, FV3},
		{Spade, FV4},
		{Spade, FV5},
		{Spade, FV6},
		{Spade, FV7},
		{Spade, FV8},
		{Spade, FV9},
		{Spade, FV10},
		{Spade, FVJ},
		{Spade, FVQ},
		{Spade, FVK},
		{Spade, FVA},
	}
}
