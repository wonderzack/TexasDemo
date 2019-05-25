package texas_holdem

import (
	"fmt"
	"strings"

	"github.com/zack-wong/TexasDemo/poker"

	"github.com/inconshreveable/log15"
)

func CardsStr2Cards(cardstr string) (cards []poker.Card) {
	h := strings.Split(cardstr, "-")
	var c poker.Card
	for _, str := range h {
		e := c.SetFromFront(str)
		if e == nil {
			cards = append(cards, poker.Card(c))
		}
	}
	return
}

func WinloseAnalyzeByPokerStr(poker string) (winOrTieRate float32, err error) {
	cards := strings.Split(poker, ":")
	log15.Warn(fmt.Sprintf("===> win lost poker %s cards %#v ", poker, cards))
	if len(cards) != 3 {
		err = fmt.Errorf("err poker str")
		return
	}

	player := CardsStr2Cards(cards[0])
	banker := CardsStr2Cards(cards[1])
	log15.Warn(fmt.Sprintf("===> win lost player %#v banker %#v ", player, banker))
	return WinloseAnalyze(player, banker), nil
}

const (
	counter_win = iota
	counter_lose
	counter_tie
)

func WinloseAnalyze(player, public []poker.Card) (winOrTieRate float32) {
	counter := make([]int32, 3)
	cardShowed := make(map[poker.Card]bool)
	for _, c := range player {
		cardShowed[c] = true
	}
	for _, c := range public {
		cardShowed[c] = true
	}

	handPlayer := NewHand()
	handPlayer.SetNeedCalIndex(false)
	handBanker := NewHand()
	handBanker.SetNeedCalIndex(false)

	banker := make([]poker.Card, 0, 2)

	if len(public) >= 5 {
		_innerAnalyze(player, public, banker, 0, cardShowed, handPlayer, handBanker, counter)
	} else {
		for i := 0; i < len(poker.Deck); i++ {
			if cardShowed[poker.Deck[i]] {
				continue
			}
			public2 := append(public, poker.Deck[i])
			_innerAnalyze(player, public2, banker, i, cardShowed, handPlayer, handBanker, counter)
		}
	}

	winOrTieCount := counter[counter_win] + counter[counter_tie]
	totalCount := winOrTieCount + counter[counter_lose]

	return float32(winOrTieCount) / float32(totalCount)

}

func _innerAnalyze(player, public, banker []poker.Card,
	i int, cardShowed map[poker.Card]bool, handPlayer, handBanker *Hand, counter []int32) {
	handPlayer.SetCard(append(player, public...))

	for j := i + 1; j < len(poker.Deck); j++ {
		if cardShowed[poker.Deck[j]] {
			continue
		}
		banker1 := append(banker, poker.Deck[j])
		for k := j + 1; k < len(poker.Deck); k++ {
			if cardShowed[poker.Deck[k]] {
				continue
			}
			banker2 := append(banker1, poker.Deck[k])
			handBanker.SetCard(append(banker2, public...))
			if handPlayer.Win(handBanker) {
				counter[counter_win]++
			} else if handPlayer.Tie(handBanker) {
				counter[counter_tie]++
			} else {
				counter[counter_lose]++
			}
		}
	}
}
