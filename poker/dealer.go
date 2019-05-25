package poker

import (
	"bytes"
	"strconv"
	"sync"

	"errors"
	"math/rand"
	"time"
)

const (
	CardNPos int = -1
)

// 发牌器
type Dealer struct {
	DealerMutex sync.RWMutex
	rand        *rand.Rand
	cards       []Card
	left        int
	next        int
}

func NewDealer(deckNum int, rawDeck []Card) *Dealer {
	Dealer := Dealer{
		rand:  rand.New(rand.NewSource(time.Now().UnixNano())),
		cards: make([]Card, deckNum*len(rawDeck)),
	}
	for i := 0; i < deckNum; i++ {
		copy(Dealer.cards[i*len(rawDeck):], rawDeck)
	}

	return &Dealer
}

func (d *Dealer) Shuffle() {
	d.DealerMutex.Lock()
	d.rand.Seed(time.Now().UnixNano())

	le := len(d.cards)

	for i := 0; i < le; i++ {
		j := d.rand.Intn(le-i) + i
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}

	d.next = 0
	d.left = len(d.cards)
	d.DealerMutex.Unlock()
}

func (d *Dealer) ReserveFindIf(pred func(card Card) bool) (index int) {
	d.DealerMutex.RLock()
	defer d.DealerMutex.RUnlock()

	index = d._totalPoker() - 1
	for index >= d.next {
		if pred(d.cards[index]) {
			return index
		}
		index--
	}
	return CardNPos
}

func (d *Dealer) FindIf(pred func(card Card) bool) (index int) {
	d.DealerMutex.RLock()
	defer d.DealerMutex.RUnlock()

	index = d.next
	for index < d._totalPoker() {
		if pred(d.cards[index]) {
			return index
		}
		index++
	}
	return CardNPos
}

func (d *Dealer) SwapPoker(i, j int) {
	d.DealerMutex.Lock()
	defer d.DealerMutex.Unlock()

	d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
}

func (d *Dealer) TotalPoker() int {
	d.DealerMutex.RLock()
	defer d.DealerMutex.RUnlock()

	return d._totalPoker()
}

func (d *Dealer) _totalPoker() int {
	return len(d.cards)
}

func (d *Dealer) PokerAt(index int) (Card, error) {
	d.DealerMutex.RLock()
	defer d.DealerMutex.RUnlock()

	if index < 0 || index >= d._totalPoker() {
		return 0, errors.New("index out of range")
	}
	return d.cards[index], nil
}

func (d *Dealer) LeftPoker() int {
	return d.left
}

func (d *Dealer) DealOne() (card Card, index int, err error) {
	d.DealerMutex.Lock()
	card, index, err = d._dealOne()
	d.DealerMutex.Unlock()
	return
}

func (d *Dealer) _dealOne() (Card, int, error) {
	if d.next >= d._totalPoker() {
		return Card(0), 0, errors.New("poker use out")
	}
	res := d.cards[d.next]
	d.left--
	d.next++
	return res, d.next - 1, nil
}

func (d *Dealer) FirstCardFront() string {
	d.DealerMutex.RLock()
	defer d.DealerMutex.RUnlock()

	if d.next == 0 || d._totalPoker() == 0 {
		return ""
	}
	return d.cards[0].Front()

}

func (d *Dealer) FirstCard() Card {
	d.DealerMutex.RLock()
	defer d.DealerMutex.RUnlock()

	if d.next == 0 || d._totalPoker() == 0 {
		return 0
	}
	return d.cards[0]
}

func (d *Dealer) SimpleDeal(n int) ([]Card, error) {
	d.DealerMutex.Lock()
	defer d.DealerMutex.Unlock()

	var res []Card
	for n != 0 {
		card, _, err := d._dealOne()
		if err == nil {
			res = append(res, card)
			n--
		} else {
			return res, err
		}
	}
	return res, nil
}

// 回退
func (d *Dealer) Back(n int) {
	d.DealerMutex.Lock()
	defer d.DealerMutex.Unlock()

	for n > 0 {
		d.left++
		d.next--
		n--
	}
}

func (d *Dealer) Trace() string {
	d.DealerMutex.RLock()
	defer d.DealerMutex.RUnlock()

	return d._trace()
}

func (d *Dealer) _trace() string {
	var buffer bytes.Buffer
	for _, card := range d.cards {
		buffer.WriteString(strconv.Itoa(int(card)))
		buffer.WriteString(";")
	}
	return buffer.String()
}
