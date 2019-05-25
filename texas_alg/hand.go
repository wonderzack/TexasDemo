package texas_holdem

// 修改自： https://github.com/SongLiangChen/TexasHoldemAI

/* 你应该先看一下这里的介绍，对你阅读会事半功倍
---------------------------------------------------------------
这里详细介绍一下各种数据是如何储存的
关于牌，我们用一个14位的二进制区间来存储，例如一副牌2-A，将表示为
11111111111110，即第二位表示牌2，第三位表示牌3，以此类推。那么第一位是干嘛的呢？请继续看下面

这样的储存方式除了省空间外还有什么优势呢？我们顺子的判断为例：
例如顺子10~~A，在二进制区间将表示为11111000000000，叫它S
现在我们有手牌2 3 10 J Q K A，那么它的二进制表示是11111000000110，叫它T
那么T&S==S的话，就可以说明T包含一个顺子，并且顺子是10~~A
S转化为10进制的话是15872
类似的我们将所有可能的顺子预先保存好，可以看到“StraightValue”这个数组，就是干这个用的
由于德州扑克里面A 2 3 4 5是最小的顺子，现在你可以明白二进制区间里第0位的作用了，和最高位一样也是保存A

我们将所有牌型做一个分级：
皇家同花顺：10
同花顺    ：9
四条      ：8
葫芦      ：7
同花      ：6
顺子      ：5
三条      ：4
两对      ：3
一对      ：2
高牌      ：1

比较的时候先比较两手牌的等级，等级相同的情况下，我们进一步分析匹配出来的5张牌的value值。
我的value的算法如下：
对你的手牌进行排序，排序规则是出现次数多的优先，次数相同的则值大的优先，比如：
7 8 4 2 2 A K，排序后为:22AK874, 匹配牌为22AK8可以理解为16进制：0x22ED8
需要注意的是，顺子和同花顺不适应此算法，顺子的value就是该顺子的最高牌；

各种牌型的判断以及比较：
1、皇家同花顺 royal flush
这个最简单了，直接用四个花色的牌集(详见straight数组)，去和15872相与即可(原理见上)
一场牌局只可能出现唯一皇家同花顺，所以只需要记录等级即可，因为只可能win or tie(五张公牌)

2、同花顺straight flush
和皇家同花顺类似，从大到小遍历所有可能的顺子，和它们做与操作。value值是该顺子中最大的高牌

3、四条 four of a kind
维护一个数组count []int用于对每一种牌值进行计数即可。
还有一种方法是将四种花色的牌集相与，最后二进制区间内还是1的那些就是我们要的。

4、葫芦 full house
根据排序好的card 如 222244K 选取前面4个 ， 再取剩下三张牌中最大的

5、同花 flush
看各个花色的个数 suitCount

6、顺子 straight
取所有花色牌的集合，去和所有可能的顺子做与操作

7、三条 three of a kind；两对 two pairs； 一对 one pair
运用count数组计数器
*/

import (
	"errors"
	"fmt"
	"sort"

	"github.com/zack-wong/TexasDemo/poker"
)

const (
	SUIT_SIZE uint32 = 4 //四种花色
	ACE_VALUE uint32 = 14
)

type CardInfo struct {
	p        poker.Card
	Showtime int // 这张牌的 牌值 出现的次数， 用來排序
}

//实现sort包中的排序接口
type Cards []*CardInfo

func (c Cards) Len() int {
	return len(c)
}

func (c Cards) Less(i, j int) bool {
	if c[i].Showtime > c[j].Showtime {
		return true
	} else if c[i].Showtime < c[j].Showtime {
		return false
	} else {
		return c[i].p.Value() > c[j].p.Value()
	}
}

func (c Cards) Swap(i, j int) {
	tmp := c[i]
	c[i] = c[j]
	c[j] = tmp
}

type HandType = uint32

const (
	HandTypeUnknown = HandType(0)  // 未知
	HighCard        = HandType(1)  // 高牌
	OnePair         = HandType(2)  // 一对
	TwoPairs        = HandType(3)  // 两对
	ThreeOfAKind    = HandType(4)  // 三条
	Straight        = HandType(5)  // 顺子
	Flush           = HandType(6)  // 同花
	FullHouse       = HandType(7)  // 葫芦
	FourOfAKind     = HandType(8)  // 四条
	StraightFlush   = HandType(9)  // 同花顺
	RoyalFlush      = HandType(10) // 皇家同花顺
)

var handTypeName = []string{
	"未知",
	"高牌",
	"一对",
	"两对",
	"三条",
	"顺子",
	"同花",
	"葫芦",
	"四条",
	"同花顺",
	"皇家同花顺",
}

func HandTypeName(t HandType) string {
	if t > 0 && t < uint32(len(handTypeName)) {
		return handTypeName[t]
	}
	return handTypeName[HandTypeUnknown]
}

// StraightValue 各个顺子代表的值
// 62    => 0b000000000111110   A2345
// 31744 => 0b111110000000000   10JQKA
var StraightValue = []uint32{31744, 15872, 7936, 3968, 1984, 992, 496, 248, 124, 62}

type Hand struct {
	cards              Cards // 储存发下来的手牌
	needCalIndex       bool  // 是否需要计算哪些牌的index 被 选中
	card2index         map[poker.Card]uint32
	suitCount          [SUIT_SIZE]uint32  //用于判断是否有同花
	valCount           [ACE_VALUE + 1]int //记录每种牌值出现的次数
	straightFlushFlags [SUIT_SIZE]uint32  //记录每种花色的牌集，用于判断同花顺
	straightFlag       uint32             //手头上四种花色牌的并集，用于判断是否有顺子

	/*
		一副手牌的最终值，Level相同的情况下，我们用SubLevel来比较大小
		例如一副手牌：3 3 3 7 7 A K，
		它的值是：33377
		数据排序规则是，出现次数多者优先，次数相同则大小优先
	*/
	SubLevel   uint32
	Level      HandType
	MatchFlag  uint32
	MatchCards []poker.Card
}

func (h *Hand) FinalLevel() uint32 {
	return h.Level<<20 | h.SubLevel // 20 mean FFFFF, SubLevel 最多占用20个位
}

func (h *Hand) HighCard() uint32 {
	if len(h.MatchCards) != 0 {
		return h.MatchCards[0].Value()
	}
	return 0
}

func (h *Hand) debug() {
	fmt.Printf("card              : ")
	for _, v := range h.cards {
		fmt.Printf(" %v", v.p)
	}
	fmt.Printf("\n")
	fmt.Printf("matchCard         : %v\n", h.MatchCards)
	fmt.Printf("valCount          : %v\n", h.valCount)
	fmt.Printf("suitCount         : %v\n", h.suitCount)
	for _, v := range h.straightFlushFlags {
		fmt.Printf("straightFlushFlags: %015b\n", v)
	}
	fmt.Printf("straightFlag      : %014b\n", h.straightFlag)
	fmt.Printf("HandTypeName      : %s\n", HandTypeName(h.Level))
	fmt.Printf("SubLevel        : %d\n", h.SubLevel)
	fmt.Printf("======================\n")
}

func NewHand() *Hand {
	h := new(Hand)
	h.needCalIndex = true
	return h
}

func (h *Hand) SetNeedCalIndex(need bool) {
	h.needCalIndex = need
}

func (h *Hand) Win(otherH *Hand) bool {
	return h.FinalLevel() > otherH.FinalLevel()
}

func (h *Hand) Tie(otherH *Hand) bool {
	return h.FinalLevel() == otherH.FinalLevel()
}

func (h *Hand) Reset() {
	i := uint32(0) // temp val

	for i = 0; i < SUIT_SIZE; i++ {
		h.straightFlushFlags[i] = 0
		h.suitCount[i] = 0
	}
	for i = 1; i < ACE_VALUE+1; i++ {
		h.valCount[i] = 0
	}

	for i = 0; i < uint32(len(h.cards)); i++ {
		if h.cards[i] != nil {
			h.cards[i].Showtime = 0
			continue
		}
		h.cards[i] = new(CardInfo)
	}
	if h.needCalIndex {
		h.card2index = map[poker.Card]uint32{}
	}
	h.straightFlag = 0
	h.Level = HandTypeUnknown
	h.SubLevel = 0
	h.MatchFlag = 0
	h.MatchCards = h.MatchCards[:0]

}

func (h *Hand) SetCard(c []poker.Card) error {
	if len(c) < 5 && len(c) > 7 {
		return errors.New("卡牌个数不支持")
	}

	if len(c) != len(h.cards) {
		h.cards = make(Cards, len(c))
	}
	h.Reset()

	for i, p := range c {
		cardInfo := h.cards[i]
		cardInfo.p = p
		if cardInfo.p.Value() == 1 {
			if h.needCalIndex {
				h.card2index[cardInfo.p] = uint32(i)
			}
			cardInfo.p.SetValue(ACE_VALUE)
		}
		if h.needCalIndex {
			h.card2index[cardInfo.p] = uint32(i)
		}
	}

	h.analyseHand()
	return nil
}

func _findBiggestCard(cards Cards) poker.Card {
	biggestCard := poker.Card(0x21)
	for _, v := range cards {
		if v.p.Value() >= biggestCard.Value() {
			biggestCard = v.p
		}
	}
	return biggestCard
}

func (h *Hand) _appendFirstNToMatch(N int) {
	if N > len(h.cards) {
		N = len(h.cards)
	}
	for index := 0; index < N; index++ {
		h.MatchCards = append(h.MatchCards, h.cards[index].p)
	}
}

func (h *Hand) _matchCard2Flag() {
	if h.needCalIndex {
		for _, c := range h.MatchCards {
			cardIndex, ok := h.card2index[c]
			if ok {
				h.MatchFlag |= (1 << cardIndex)
			}
		}
	}
}

func (h *Hand) _analyseIsRoyalFlush() bool {
	//由大到小来判断手牌等级
	//判断是否有皇家同花顺
	for suit := uint32(0); suit < SUIT_SIZE; suit++ {
		if h.straightFlushFlags[suit]&StraightValue[0] == StraightValue[0] {
			h.Level = RoyalFlush
			for v := uint32(10); v <= ACE_VALUE; v++ { // 10 J Q K A
				h.MatchCards = append(h.MatchCards, poker.MakeCard(v, suit+1))
			}
			return true
		}
	}
	return false
}

func (h *Hand) _analyseIsStraightFlush() bool {
	//判断是否有同花顺，由于只有可能出现一个花色的同花顺，所以记录高牌的值即可比较两个同花顺大小
	for suit := uint32(0); suit < SUIT_SIZE; suit++ {
		if h.straightFlushFlags[suit] < 31 /*0b11111*/ {
			continue
		}

		for j := 1; j < len(StraightValue); j++ {
			if h.straightFlushFlags[suit]&StraightValue[j] == StraightValue[j] {
				h.Level = StraightFlush
				highCardValue := uint32(len(StraightValue) - j + 4)
				h.SubLevel = highCardValue
				for v := highCardValue - 4; v <= highCardValue; v++ {
					h.MatchCards = append(h.MatchCards, poker.MakeCard(v, suit+1))
				}
				return true
			}
		}
	}
	return false
}

func (h *Hand) _analyseIsFourOfAKind() bool {

	if h.cards[0].Showtime == 4 {
		biggestCard := _findBiggestCard(h.cards[4:])
		h.Level = FourOfAKind
		h._appendFirstNToMatch(4)
		h.MatchCards = append(h.MatchCards, biggestCard)
		h.SubLevel = turnToValue(h.MatchCards)
		return true
	}
	return false
}

func (h *Hand) _analyseIsFullHouse() bool {
	//判断葫芦，和四条同理
	if h.cards[0].Showtime == 3 && h.cards[3].Showtime >= 2 {
		h.Level = FullHouse
		h._appendFirstNToMatch(5)
		h.SubLevel = turnToValue(h.MatchCards)

		return true
	}
	return false
}

func (h *Hand) _analyseIsFlush() bool {
	/*判断同花，
	还是同理，场上有且只有可能出现一个花色的同花
	都是同花的情况下，就比较谁的同花大
	*/
	for suit := uint32(0); suit < SUIT_SIZE; suit++ {
		suitCount := h.suitCount[suit]
		if suitCount >= 5 {
			h.Level = Flush
			straightFlushFlags := h.straightFlushFlags[suit]
			for value := ACE_VALUE; value >= uint32(2); value-- {
				if straightFlushFlags&(1<<value) != 0 {
					h.MatchCards = append(h.MatchCards, poker.MakeCard(value, suit+1))
				}
				if len(h.MatchCards) == 5 {
					break
				}
			}
			h.SubLevel = turnToValue(h.MatchCards)
			return true
		}
	}
	return false
}

func (h *Hand) _analyseIsStraight() bool {
	//判断顺子，handvalue保存的是所有花色rank的并集，和同花顺同理
	for i := 0; i < len(StraightValue); i++ {
		if h.straightFlag&StraightValue[i] == StraightValue[i] {
			highCardValue := uint32(len(StraightValue) - i + 4)
			h.Level = Straight
			h.SubLevel = highCardValue

			usebit := uint32(0)
			for _, card := range h.cards {
				if (card.p.Value() <= highCardValue && card.p.Value() >= highCardValue-4) ||
					(highCardValue == 5 && card.p.Value() == ACE_VALUE) {
					if usebit&(1<<card.p.Value()) == 0 {
						h.MatchCards = append(h.MatchCards, card.p)
						usebit |= (1 << card.p.Value())
					}
				}
			}
			return true
		}
	}
	return false
}

func (h *Hand) _analyseIsThreeOfAKind() bool {
	//判断三条
	if h.cards[0].Showtime == 3 &&
		h.cards[3].Showtime == 1 &&
		h.cards[4].Showtime == 1 {
		h.Level = ThreeOfAKind
		h._appendFirstNToMatch(5)
		h.SubLevel = turnToValue(h.MatchCards)

		return true
	}
	return false
}

func (h *Hand) _analyseIsTwoPairs() bool {
	/*
		判断两对，首先我们要确定有没有两对，都有两对的情况下，也有可能出现平局的情况
		所以判断依据是将手牌排序，出现次数多的牌优先，次数相同的情况下，牌值大的优先
		最后将排序转化成16进制int，直接比较即可
	*/
	if h.cards[0].Showtime == 2 &&
		h.cards[2].Showtime == 2 {
		biggestCard := _findBiggestCard(h.cards[4:])
		h.Level = TwoPairs
		h._appendFirstNToMatch(4)
		h.MatchCards = append(h.MatchCards, biggestCard)
		h.SubLevel = turnToValue(h.MatchCards)
		return true
	}
	return false
}

func (h *Hand) _analyseIsOnePair() bool {
	//判断一对
	if h.cards[0].Showtime == 2 {
		h.Level = OnePair
		h._appendFirstNToMatch(5)
		h.SubLevel = turnToValue(h.MatchCards)
		return true

	}
	return false
}

func (h *Hand) analyseHand() {
	h._analyCards()
	sort.Sort(h.cards)

	// 利用 运算符 || 的特殊性质，只要有真 就不往下判断
	if h._analyseIsRoyalFlush() ||
		h._analyseIsStraightFlush() ||
		h._analyseIsFourOfAKind() ||
		h._analyseIsFullHouse() ||
		h._analyseIsFlush() ||
		h._analyseIsStraight() ||
		h._analyseIsThreeOfAKind() ||
		h._analyseIsTwoPairs() ||
		h._analyseIsOnePair() {
		h._matchCard2Flag()
		return
	}
	//判断高牌
	h.Level = HighCard
	h._appendFirstNToMatch(5)
	h.SubLevel = turnToValue(h.MatchCards)
	h._matchCard2Flag()
	return
}

//将手牌转化成整数形式
func turnToValue(cards []poker.Card) uint32 {
	/*
		一副手牌的最终值，Level相同的情况下，我们用SubLevel来比较大小
		例如一副手牌：3 3 3 7 7 A K，
		它的值是：33377AK
		数据排序规则是，出现次数多者优先，次数相同则大小优先
	*/
	res := uint32(0)
	for _, c := range cards {
		res = (res << 4) | c.Value()
	}
	return res
}

func (h *Hand) _analyCards() {
	var suitIndex, pokerVal uint32

	for _, c := range h.cards {
		suitIndex = c.p.Suit() - 1
		pokerVal = c.p.Value()

		h.suitCount[suitIndex]++
		h.straightFlushFlags[suitIndex] |= 1 << uint(pokerVal)
		h.straightFlag |= 1 << uint(pokerVal)

		if pokerVal == ACE_VALUE { //A也保存到第2位中去，用于判断A2345这样的顺子
			h.straightFlushFlags[suitIndex] |= 2 // 2 mean 1 << 1
			h.straightFlag |= 2                  // 2 mean 1 << 1
		}

		h.valCount[pokerVal]++
	}

	for _, c := range h.cards {
		c.Showtime = h.valCount[c.p.Value()]
	}
}
