package poker

import (
	"fmt"
	"strconv"
	"strings"
)

// Card card type
type Card byte

//Pokers
const (
	AceSpades   = Card(0x14)
	TwoSpades   = Card(0x24)
	ThreeSpades = Card(0x34)
	FourSpades  = Card(0x44)
	FiveSpades  = Card(0x54)
	SixSpades   = Card(0x64)
	SevenSpades = Card(0x74)
	EightSpades = Card(0x84)
	NineSpades  = Card(0x94)
	TenSpades   = Card(0xA4)
	JackSpades  = Card(0xB4)
	QueenSpades = Card(0xC4)
	KingSpades  = Card(0xD4)

	AceHearts   = Card(0x13)
	TwoHearts   = Card(0x23)
	ThreeHearts = Card(0x33)
	FourHearts  = Card(0x43)
	FiveHearts  = Card(0x53)
	SixHearts   = Card(0x63)
	SevenHearts = Card(0x73)
	EightHearts = Card(0x83)
	NineHearts  = Card(0x93)
	TenHearts   = Card(0xA3)
	JackHearts  = Card(0xB3)
	QueenHearts = Card(0xC3)
	KingHearts  = Card(0xD3)

	AceClubs   = Card(0x12)
	TwoClubs   = Card(0x22)
	ThreeClubs = Card(0x32)
	FourClubs  = Card(0x42)
	FiveClubs  = Card(0x52)
	SixClubs   = Card(0x62)
	SevenClubs = Card(0x72)
	EightClubs = Card(0x82)
	NineClubs  = Card(0x92)
	TenClubs   = Card(0xA2)
	JackClubs  = Card(0xB2)
	QueenClubs = Card(0xC2)
	KingClubs  = Card(0xD2)

	AceDiamonds   = Card(0x11)
	TwoDiamonds   = Card(0x21)
	ThreeDiamonds = Card(0x31)
	FourDiamonds  = Card(0x41)
	FiveDiamonds  = Card(0x51)
	SixDiamonds   = Card(0x61)
	SevenDiamonds = Card(0x71)
	EightDiamonds = Card(0x81)
	NineDiamonds  = Card(0x91)
	TenDiamonds   = Card(0xA1)
	JackDiamonds  = Card(0xB1)
	QueenDiamonds = Card(0xC1)
	KingDiamonds  = Card(0xD1)

	BlackJoker = Card(0xE0)
	RedJoker   = Card(0xF0)
)

var valStr = []string{
	"_", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "♚", "♔",
}

var suitStr = []string{
	"", "♦", "♣", "♥", "♠",
}

func (card Card) String() string {
	num := card.Value()
	suit := card.Suit()

	if int(num) < len(valStr) && int(suit) < len(suitStr) {
		return fmt.Sprintf("%s%s", valStr[num], suitStr[suit])
	}

	return fmt.Sprintf("%d", uint32(card))
}

func (card Card) Value() uint32 {
	return uint32(card) >> 4
}

func MakeCard(value uint32, suit uint32) Card {
	return Card(value<<4 | suit)
}

func (card *Card) SetValue(v uint32) {
	*card = Card(card.Suit() | v<<4)
}

func (card Card) Suit() uint32 {
	return uint32(card) & 0xF
}

func (card Card) Number() uint32 {
	v := card.Value()
	if v >= 10 {
		v = 0
	}
	return v
}

func (card Card) Front() string {
	return fmt.Sprintf("%d%d", card.Value(), card.Suit())
}

func (card *Card) SetFromFront(f string) error {
	v, err := strconv.Atoi(f)
	if err != nil {
		return err
	}
	suit := uint32(v % 10)
	val := uint32(v / 10)
	c := MakeCard(val, suit)
	*card = c
	return nil
}

var Deck = []Card{
	AceSpades, TwoSpades, ThreeSpades, FourSpades, FiveSpades, SixSpades, SevenSpades, EightSpades, NineSpades, TenSpades, JackSpades, QueenSpades, KingSpades,
	AceHearts, TwoHearts, ThreeHearts, FourHearts, FiveHearts, SixHearts, SevenHearts, EightHearts, NineHearts, TenHearts, JackHearts, QueenHearts, KingHearts,
	AceClubs, TwoClubs, ThreeClubs, FourClubs, FiveClubs, SixClubs, SevenClubs, EightClubs, NineClubs, TenClubs, JackClubs, QueenClubs, KingClubs,
	AceDiamonds, TwoDiamonds, ThreeDiamonds, FourDiamonds, FiveDiamonds, SixDiamonds, SevenDiamonds, EightDiamonds, NineDiamonds, TenDiamonds, JackDiamonds, QueenDiamonds, KingDiamonds,
}

var Deck54 = []Card{}

func init() {
	// Deck include joker
	Deck54 = append(Deck, BlackJoker, RedJoker)
}

func String2pokers(s string) []Card {
	strs := strings.Split(s, ",")
	cards := make([]Card, len(strs))

	for i := 0; i < len(strs); i++ {
		cards[i].SetFromFront(strings.Trim(strs[i], " "))
	}
	return cards
}

func CardsToUint32s(c []Card) []uint32 {
	ret := make([]uint32, len(c))
	for i := 0; i < len(c); i++ {
		ret[i] = uint32(c[i])
	}
	return ret
}
