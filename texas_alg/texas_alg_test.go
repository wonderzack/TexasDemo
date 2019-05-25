package texas_holdem

import (
	"fmt"
	"testing"
	"time"

	"github.com/zack-wong/TexasDemo/poker"
)

type test_case struct {
	cards    string
	handtype HandType
	match    string
}

func reverseString(s string) string {
	runes := []rune(s)

	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}

	return string(runes)
}

func CasesTest(c []test_case, t *testing.T) {
	h := NewHand()

	for _, c := range c {
		h.SetCard(poker.String2pokers(c.cards))
		h.debug()
		if h.Level != c.handtype {
			t.Error("err level")
		}
		if fmt.Sprintf("%07b", h.MatchFlag) != reverseString(c.match) {
			t.Error("err match")
		}

	}
}

func TestRoyalFlush(t *testing.T) {
	var testcases = []test_case{
		//=======================
		{"11, 131, 121, 111, 101, 21, 31", RoyalFlush, "1111100"},
		{"12, 132, 122, 112, 102, 92, 82", RoyalFlush, "1111100"},
		{"13, 133, 123, 113, 103, 91, 81", RoyalFlush, "1111100"},
		{"14, 134, 124, 114, 104, 84, 74", RoyalFlush, "1111100"},
	}
	CasesTest(testcases, t)
}

func TestStraightFlush(t *testing.T) {
	var testcases = []test_case{
		//=======================
		{"91, 131, 121, 111, 101, 21, 31", StraightFlush, "1111100"},
		{"22, 132, 122, 112, 102, 92, 82", StraightFlush, "0111110"},
		{"11, 133, 123, 113, 103, 91, 93", StraightFlush, "0111101"},
		{"11, 131, 124, 114, 104, 94, 84", StraightFlush, "0011111"},
	}
	CasesTest(testcases, t)
}

func TestFourOfAKind(t *testing.T) {
	var testcases = []test_case{
		//=======================
		{"101, 102, 103, 104, 91, 21, 31", FourOfAKind, "1111100"},
		{"22, 21, 91, 92, 93, 94, 81", FourOfAKind, "0011111"},
	}
	CasesTest(testcases, t)

	h1 := NewHand()
	h1.SetCard(poker.String2pokers("21, 91, 92, 93, 94, 81, 22"))

	h2 := NewHand()
	h2.SetCard(poker.String2pokers("21, 91, 92, 93, 94, 101, 23"))

	if h1.Win(h2) {
		t.Error("not win")
	}
}

func TestFullHouse(t *testing.T) {
	var testcases = []test_case{
		//=======================
		{"101, 102, 103, 91, 92, 93, 31", FullHouse, "1111100"},
		{"22, 21, 91, 92, 93, 84, 81", FullHouse, "0011111"},
	}
	CasesTest(testcases, t)
}

func TestFlush(t *testing.T) {
	var testcases = []test_case{
		//=======================
		{"101, 91, 81, 71, 51, 52, 53", Flush, "1111100"},
		{"102, 92, 82, 71, 51, 52, 42", Flush, "1110011"},
	}
	CasesTest(testcases, t)

	h1 := NewHand()
	h1.SetCard(poker.String2pokers("101, 91, 81, 71, 51, 52, 121"))
	h1.debug()

	h2 := NewHand()
	h2.SetCard(poker.String2pokers("101, 91, 81, 71, 51, 72, 131"))
	h2.debug()

	if h1.Win(h2) {
		t.Error("not win")
	}
}

func TestStraight(t *testing.T) {
	var testcases = []test_case{
		//=======================
		{"101, 91, 81, 72, 62, 52, 53", Straight, "1111100"},
		{"102, 92, 82, 51, 41, 72, 61", Straight, "1110011"},
	}
	CasesTest(testcases, t)

}

func TestThreeOfAKind(t *testing.T) {
	var testcases = []test_case{
		//=======================
		{"91, 92, 93, 81, 72, 63, 43", ThreeOfAKind, "1111100"},
		{"91, 92, 93, 51, 41, 132, 121", ThreeOfAKind, "1110011"},
	}
	CasesTest(testcases, t)
}

func TestTwoPairs(t *testing.T) {
	var testcases = []test_case{
		//=======================
		{"91, 92, 63, 61, 82, 83, 103", TwoPairs, "1100111"},
	}
	CasesTest(testcases, t)
}

func TestOnePair(t *testing.T) {
	var testcases = []test_case{
		//=======================
		{"91, 92, 23, 51, 81, 114, 123", OnePair, "1100111"},
	}
	CasesTest(testcases, t)
}

func TestHighCard(t *testing.T) {
	var testcases = []test_case{
		//=======================
		{"11, 32, 53, 74, 91, 112, 133", HighCard, "1001111"},
	}
	CasesTest(testcases, t)
}

func Test5card(t *testing.T) {
	h1 := NewHand()
	h1.SetCard(poker.String2pokers("101, 91, 81, 71, 51"))
	h1.debug()
}

func TestMatchBit(t *testing.T) {
	var testcases = []test_case{
		//=======================
		{"21, 12, 113, 43, 52, 82, 34", Straight, "1101101"},
	}
	CasesTest(testcases, t)
}

func Test_win(t *testing.T) {

	deck := poker.NewDealer(1, poker.Deck)
	deck.Shuffle()

	player, _ := deck.SimpleDeal(2)
	public, _ := deck.SimpleDeal(3)

	fmt.Println("hand", player)
	fmt.Println("public", public)
	tick := time.Now()
	fmt.Println("winOrTieRate", WinloseAnalyze(player, public))
	delta := time.Since(tick)
	fmt.Println("use time", delta)
}

func Test_win2(t *testing.T) {

	tick := time.Now()
	w, _ := WinloseAnalyzeByPokerStr("141-142:143-144-134-133:")
	fmt.Println("winOrTieRate", w)
	delta := time.Since(tick)
	fmt.Println("use time", delta)
}

func TestSplitPond(t *testing.T) {

	betStatuss := make([]IBetStatus, 10)

	for i := 0; i < 10; i++ {
		b := &BetStatus{
			seatID:    int32(i + 20),
			winVal:    uint32(i / 2),
			betAmount: int64(i * 100),
		}
		betStatuss[i] = b
		fmt.Printf("%s", b)
	}
	tick := time.Now()

	seatID2WinAmount := DistributePond(betStatuss)
	delta := time.Since(tick)
	fmt.Println("use time", delta)
	fmt.Println("奖池分配")
	for i := 0; i < 10; i++ {
		seat := int32(i + 20000)
		fmt.Printf(" %d:%d\n", seat, seatID2WinAmount[seat])
	}
}
